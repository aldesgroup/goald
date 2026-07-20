// ------------------------------------------------------------------------------------------------
// Here we implement the generic business logic involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

// CreateBO creates a new business object in the database, and all the other business objects
// that are linked exclusively to it (children) if any, in a single transaction
func CreateBO(bloCtx BloContext, bObj IBusinessObject) error {
	// let's start a transaction if none is already started
	beginTransactionHere, errBegin := bloCtx.BeginTransaction(DbFor(bObj))
	if errBegin != nil {
		return ErrorC(errBegin, "Could not create object since a transaction could not be started")
	}

	// let's try to create the given entity
	createErr := doCreateBO(bloCtx, bObj)

	// let's end the transaction
	if beginTransactionHere {
		if errEnd := bloCtx.EndTransaction(createErr); errEnd != nil {
			return ErrorC(errEnd, "Could not create object since the current transaction could not be terminated")
		}
	}

	return createErr
}

// doCreateBO does the actual creation of a new business object in the database
func doCreateBO(bloCtx BloContext, bObj IBusinessObject) error {
	if bObj == nil {
		return nil
	}

	// the business object should not have an ID already
	if bObj.GetID() != 0 {
		return Error("Could not create object since it already has an ID (%d)", bObj.GetID())
	}

	// do we have stuff to perform on the __BOBJ__ before inserting it ?
	if err := bObj.ChangeBeforeInsert(bloCtx); err != nil {
		return ErrorC(err, "Could not create object since the pre-insert got an error")
	}

	// check of the validity regarding the model constraints
	if err := bObj.getClass(bObj).IsModelValid(bObj); err != nil {
		return ErrorC(err, "Could not create object since it is not valid regarding its model constraints")
	}

	// check of "functional / business" validity
	if err := bObj.IsValid(bloCtx); err != nil {
		return ErrorC(err, "Could not create object since it is not valid")
	}

	// setting some tracking info
	// bObj.SetCreatedByID(biContext.GetCurrentUser().GetID())
	// bObj.SetCreatedBy(biContext.GetCurrentUser().GetLabel())
	// bObj.SetCreation(core.Now())
	// bObj.Set__BOBJ__Status(__BOBJ__StatusCREATED)

	// pushing to the DB ! We're going to add a new line within the __BOBJ__'s table
	if err := dbInsert(bloCtx.DaoFor(bObj), bObj); err != nil {
		return ErrorC(err, "Could not create object because of a problem with the DB")
	}

	// TODO inserting all the links that waited for the current entity to be inserted in DB before getting inserted themselves
	// if err := doCreateChildrenEntities(biContext, entity, ""); err != nil {
	// 	if biContext.IsVerbose() {
	// 		biContext.Log().Trace("Entity we tried to create in DB: (see below)\n%s", ToString(entity, 2, true))
	// 	}

	// 	return NewErrC(err, "Could not completely create entity '%s' since inserting its children crashed", ToEntityFullReference(entity))
	// }

	// we have stuff to do after the insertion ? yeah ? really ? let's do it now !
	if err := bObj.ChangeAfterInsert(bloCtx); err != nil {
		return ErrorC(err, "Could not post-insert object since it got an error")
	}

	// // 'guess everything is alrite here
	// return nilreturn nil
	return nil
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
