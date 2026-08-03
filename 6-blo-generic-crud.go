// ------------------------------------------------------------------------------------------------
// Here we implement the generic business logic involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import "github.com/aldesgroup/goald/features/utils"

// CreateBusinessObjects creates a new business object in the database, and all the other business objects
// that are linked exclusively to it (children) if any, in a single transaction/
// WARNING: all the business objects must be of the same type (same model), this does not handle polymorphism
func CreateBusinessObjects[BOTYPE IBusinessObject](bloCtx BloContext, bObjs ...BOTYPE) (createErr error) {
	if len(bObjs) == 0 {
		return nil
	}

	// we know the BOs here are all of the same model
	model := bObjs[0].getModel(bObjs[0])

	// let's start a transaction if none is already started
	beginTransactionHere, errBegin := bloCtx.BeginTransaction(model)
	if errBegin != nil {
		return ErrorC(errBegin, "Could not create object since a transaction could not be started")
	}

	// let's make sure the transaction is always ended, even if a panic occurs below,
	// so that we never leave a dangling transaction / corrupt the BloContext's tx state
	if beginTransactionHere {
		bloCtx.Trace("New transaction started here!")
		defer func() {
			if r := recover(); r != nil {
				bloCtx.Trace("Recovered from panic while creating object(s): %v", r)
				errEnd := bloCtx.EndTransaction(Error("Recovered from panic while creating object(s): %v", r))
				panic(errEnd)
			}

			if errEnd := bloCtx.EndTransaction(createErr); errEnd != nil {
				createErr = ErrorC(errEnd, "Could not create object since the current transaction could not be terminated")
			} else {
				bloCtx.Trace("Transaction ended here!")
			}
		}()
	}

	// let's try to create the given entity
	createErr = doCreateBusinessObjects(bloCtx, model, bObjs...)

	return
}

// doCreateBO does the actual creation of a new business object in the database
func doCreateBusinessObjects[BOTYPE IBusinessObject](bloCtx BloContext, model IBusinessObjectModel, bObjs ...BOTYPE) error {
	if len(bObjs) == 0 {
		return nil
	}

	// we need this type conversion later on, let's do it now since we're iterating over the given business objects anyway
	iBObjs := make([]IBusinessObject, len(bObjs))

	// performing some checks on the given business objects
	for i, bObj := range bObjs {
		// the business object should not have an ID already
		if bObj.GetID() != 0 {
			return Error("Could not create object since it already has an ID (%d)", bObj.GetID())
		}

		// do we have stuff to perform on the __BOBJ__ before inserting it ?
		if err := bObj.ChangeBeforeInsert(bloCtx); err != nil {
			return ErrorC(err, "Could not create object since the pre-insert got an error")
		}

		// check of the validity regarding the model constraints
		if err := bObj.IsModelValid(); err != nil {
			return ErrorC(err, "Could not create object since it is not valid regarding its model constraints")
		}

		// check of "functional / business" validity
		if err := bObj.IsValid(bloCtx); err != nil {
			return ErrorC(err, "Could not create object since it is not valid")
		}

		// let's also keep track of the BO number in the given slice, so that we can use it later on for DB ID consolidation
		bObj.setPreID(i + 1)

		// adding the business object to the interface slice
		iBObjs[i] = bObj
	}

	// // setting some tracking info
	// // bObj.SetCreatedByID(biContext.GetCurrentUser().GetID())
	// // bObj.SetCreatedBy(biContext.GetCurrentUser().GetLabel())
	// // bObj.SetCreation(core.Now())
	// // bObj.Set__BOBJ__Status(__BOBJ__StatusCREATED)

	// pushing to the DB ! We're going to add a new line within the __BOBJ__'s table
	if err := dbInsert(bloCtx.daoFor(model), iBObjs...); err != nil {
		return err
	}

	// TODO inserting all the links that waited for the current entity to be inserted in DB before getting inserted themselves
	if err := doCreateDependentBusinessObjects(bloCtx, model, iBObjs...); err != nil {
		return err
	}

	// we have stuff to do after the insertion ? yeah ? really ? let's do it now !
	for _, bObj := range bObjs {
		if err := bObj.ChangeAfterInsert(bloCtx); err != nil {
			return ErrorC(err, "Could not post-insert object since it got an error")
		}
	}

	// 'guess everything is alrite here
	return nil
}

// doCreateDependentBusinessObjects creates all the business objects that are linked exclusively to the given business objects (children) if any, in a single transaction
func doCreateDependentBusinessObjects(bloCtx BloContext, model IBusinessObjectModel, iBObjs ...IBusinessObject) error {
	// handling all the children relationships this model of business objects may have, if any
	for _, relationship := range model.getRelationshipsWithRequiredBackref() {
		// gathering all the children, by type - because we may have polymorphic children, and we want to create them in batches of the same type
		childrenByType := make(map[utils.ModelName][]IBusinessObject)
		for _, bObj := range iBObjs {
			// getting the children of this business object for this relationship
			children, errGet := bObj.GetMultipleRelationshipValue(relationship.name)
			if errGet != nil {
				return errGet
			}

			// grouping the children by type, so that we can create them in batches of the same type
			for _, child := range children {
				childrenByType[child.GetModelName()] = append(childrenByType[child.GetModelName()], child)
			}
		}

		// creating the children, by type
		for childType, children := range childrenByType {
			childModel := modelFor(childType, true)
			if err := doCreateBusinessObjects(bloCtx, childModel, children...); err != nil {
				return ErrorC(err, "Could not create child objects for relationship '%s.%s'", model.getName(), relationship.name)
			}
		}
	}

	return nil
}

func SearchBusinessObjects(bloCtx BloContext, query IQuery, values IQueryParamsObject) ([]IBusinessObject, error) {
	// performing some actions on the query params values before searching
	if err := values.DoBeforeSearch(bloCtx); err != nil {
		return nil, ErrorC(err, "Could not search for business objects since the query params values got an error")
	}

	// first we need to make sure the query params values are valid
	if err := values.IsModelValid(); err != nil {
		return nil, ErrorC(err, "Could not search for business objects since the query is not valid")
	}

	// doing the selection of the business objects from the database
	results, err := dbSelect(bloCtx.daoFor(query.getModel()), query.getName(), values)
	if err != nil {
		return nil, err
	}

	// TODO access control
	// // we check that this entity can be read, given the context
	// if err = loadedEntity.CanBeRead(biContext, loadingID); err != nil {
	// 	return nil, NewErrC(err, "Could not read entity '%s' given the context", fullReference)
	// }

	// // now that we're okey with our loading, me might have to perform some specific actions
	// if err = loadedEntity.DoAfterRead(biContext, loadingID); err != nil {
	// 	return nil, NewErrC(err, "Error occurring after reading entity '%s'", referenceOrID)
	// }

	// not much more logic for now
	return results, nil
}

// func ReadBO(bloCtx BloContext, idProp IField, idPropVal string, loadingType LoadingType) (IBusinessObject, error) {
// 	loadedBOs, errLoad := dbLoadOne(bloCtx.GetDaoContext(), idProp, idPropVal)

// 	if errLoad != nil {
// 		return nil, ErrorC(errLoad, "error while loading one instance of '%s' (%s)", idProp.ownerModel().base().name, idPropVal)
// 	}

// 	// TODO add post read

// 	return loadedBOs, nil
// }

// func LoadBOs[ResourceType IBusinessObject](bloCtx BloContext, model IBusinessObjectModel, loadingType LoadingType) ([]ResourceType, error) {
// 	// func LoadBOs(bloCtx BloContext, model IBusinessObjectModel, loadingType LoadingType) ([]ResourceType, error) {
// 	loadedBOs, errLoad := dbLoadList[ResourceType](bloCtx.GetDaoContext(), model)
// 	// loadedBOs, errLoad := dbLoadList(bloCtx.GetDaoContext(), model)

// 	if errLoad != nil {
// 		return nil, ErrorC(errLoad, "error while loading a list of '%s'", model.base().name)
// 	}

// 	// TODO add post read, i.e.:
// 	// - reading the links, using the LoadingType
// 	// - on each BO: setting the loadingID + check if reading is ok, then do after read changes

// 	return loadedBOs, nil
// }

// func DeleteBO(bloCtx BloContext, idProp IField, idPropVal string) (IBusinessObject, error) {
// 	loadedBOs, errLoad := dbRemoveOne(bloCtx.GetDaoContext(), idProp, idPropVal)

// 	if errLoad != nil {
// 		return nil, ErrorC(errLoad, "error while deleting one instance of '%s' (%s)", idProp.ownerModel().base().name, idPropVal)
// 	}

// 	// TODO add post read

// 	return loadedBOs, nil
// }

// func UpdateBO(bloCtx BloContext, input IBusinessObject, loadingType LoadingType) error {
// 	if errUpd := dbUpdate(bloCtx.GetDaoContext(), input); errUpd != nil {
// 		return ErrorC(errUpd, "error while updating one instance of '%T' (ID = %d)", input, input.GetID())
// 	}

// 	return nil
// }
