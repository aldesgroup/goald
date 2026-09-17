// ------------------------------------------------------------------------------------------------
// Here we implement the generic business logic involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import (
	"slices"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// SearchBusinessObjects searches for business objects based on the given query and search parameter values, and loads them according to the specified loading configuration.
func SearchBusinessObjects(bloCtx BloContext, query IQuery, values ISearchParamValues, loadingConf ILoadingConfig) ([]IBusinessObject, error) {
	// performing some actions on the query params values before searching
	if err := values.DoBeforeSearch(bloCtx); err != nil {
		return nil, ErrorC(err, "Could not search for business objects since the query params values got an error")
	}

	// first we need to make sure the query params values are valid
	if err := values.IsModelValid(); err != nil {
		return nil, ErrorC(err, "Could not search for business objects since the query is not valid")
	}

	// the relevant DAO for the task
	dao := bloCtx.daoFor(query.getSearchedObjectsModel().GetName())

	// doing the search for the business objects in the database
	results, err := dao.ExecSearchQuery(query.getName(), values, bloCtx.bObjCache())
	if err != nil {
		return nil, err
	}

	// performing some actions on the business objects after reading them
	if err := afterLoad(bloCtx, loadingConf, results...); err != nil {
		return nil, ErrorC(err, "Could not retrieve business objects since the post-load got an error")
	}

	// not much more logic for now
	return results, nil
}

// CreateBusinessObjects creates new business objects in the database
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

// doCreateBusinessObjects does the actual creation of a new business objects in the database
func doCreateBusinessObjects[BOTYPE IBusinessObject](bloCtx BloContext, model IBusinessObjectModel, bObjs ...BOTYPE) error {
	if len(bObjs) == 0 {
		return nil
	}

	// we need this type conversion later on, let's do it now since we're iterating over the given business objects anyway
	iBObjs := make([]IBusinessObject, len(bObjs))

	// performing some checks on the given business objects
	for i, bObj := range bObjs {
		// the business object should not have an ID already
		if bObj.GetID() > 0 {
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

		// let's also keep track of the BO number in the given slice, so that we can

		// setting some tracking info
		bObj.setPreID(i + 1) // used later on for DB ID consolidation
		bObj.setCreation(core.Now())
		// // bObj.SetCreatedByID(biContext.GetCurrentUser().GetID())
		// // bObj.SetCreatedBy(biContext.GetCurrentUser().GetLabel())
		// // bObj.Set__BOBJ__Status(__BOBJ__StatusCREATED)

		// adding the business object to the interface slice
		iBObjs[i] = bObj
	}

	// we need the DAO for the given BO type here
	dao := bloCtx.daoFor(model.GetName())

	// executing the insert request, which should result in all the BOs getting an ID
	rowToIDMap, errCreate := dao.ExecCreateQuery(iBObjs...)
	if errCreate != nil {
		return ErrorC(errCreate, "Error while performing insert query for '%s'", dao.getModel().GetName())
	}

	// consolidating the DB IDs back into the business objects
	for _, bObj := range bObjs {
		bObj.setID(rowToIDMap[bObj.GetPreID()])
	}

	// so far, we've just handle the entities' properties and persisted single links; let's now handle the multiple links
	if errLinks := dao.ExecCreateLinksQueries(iBObjs...); errLinks != nil {
		return ErrorC(errLinks, "Error while performing insert links for '%s'", dao.getModel().GetName())
	}

	// TODO inserting all the links that waited for the current entity to be inserted in DB before getting inserted themselves

	// also, we take the opportunity of creating all the business objects that are linked exclusively
	// to the given business objects (children) if any, in a single transaction
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
				return ErrorC(err, "Could not create child objects for relationship '%s.%s'", model.GetName(), relationship.name)
			}
		}
	}

	return nil
}

// ReadBusinessObjects reads the specified business objects from the database according to the given loading configuration.
func ReadBusinessObjects[BOTYPE IBusinessObject](bloCtx BloContext, loadingConf ILoadingConfig, bObjs ...BOTYPE) error {
	// logObjs(bloCtx, bObjs, "Reading business objects with loading config: %s", loadingConf.ToString())

	// first, caching the business objects, whilst computing the IDs mapped by their model
	allBObjIDs := bloCtx.bObjCache().AddAndGetIDs(bObjs)

	// then dealing with all the business objects of the same model
	for modelName, bObjIDs := range allBObjIDs {
		// sorting the IDs, that should help the queries to be more efficient
		slices.SortFunc(bObjIDs, func(a, b any) int {
			return core.IfThenElse(a.(BObjID) < b.(BObjID), -1, 1)
		})

		// the DAO for the current model
		dao := bloCtx.daoFor(modelName)

		// logIDs(bloCtx, bObjIDs, "Reading in DB '%s' instances with IDs", modelName)

		// calling the right DAO method
		if errRead := dao.ExecReadQuery(bObjIDs, bloCtx.bObjCache()); errRead != nil {
			return ErrorC(errRead, "Error while performing read query for '%s'", dao.getModel().GetName())
		}

		// checking if all the business objects were read
		for _, bObjID := range bObjIDs {
			bObj := bloCtx.bObjCache().Get(modelName, bObjID.(BObjID))
			if bObj == nil || bObj.GetCreation() == nil {
				return ErrorC(ErrBoNotFound, "Business object '%s' with ID %d was not found in the database", modelName, bObjID)
			}
		}
	}

	// performing some actions on the business objects after reading them
	if err := afterLoad(bloCtx, loadingConf, bObjs...); err != nil {
		return ErrorC(err, "Could not read business objects since the post-load got an error")
	}

	return nil
}

// afterLoad performs some actions on the business objects after reading them
func afterLoad[BOTYPE IBusinessObject](bloCtx BloContext, loadingConf ILoadingConfig, bObjs ...BOTYPE) error {
	// we check that this entity can be read, given the context, and maybe do some changes
	for _, bObj := range bObjs {
		if err := bObj.CheckAndChangeAfterRead(bloCtx); err != nil {
			return err
		}
	}

	// now, loading some relationships, if the loading config requires it
	if err := ReadRelationships(bloCtx, loadingConf, bObjs...); err != nil {
		return ErrorC(err, "Could not read business objects' relationships")
	}

	return nil
}

// ReadRelationships reads the relationships of the given business objects, according to the given loading configuration
func ReadRelationships[BOTYPE IBusinessObject](bloCtx BloContext, loadingConf ILoadingConfig, bObjs ...BOTYPE) error {
	// logObjs(bloCtx, bObjs, "Reading relationships with loading config: %s", loadingConf.ToString())

	if len(loadingConf.getRelationshipsToLoad()) > 0 {
		// first, caching the business objects, whilst computing the IDs mapped by their model
		allBObjIDs := bloCtx.bObjCache().AddAndGetIDs(bObjs)

		// reading the relationships of the given business objects, according to the given loading configuration
		for _, subLoadingConf := range loadingConf.getRelationshipsToLoad() {
			// the relationship concerned by this loading config
			relationship := subLoadingConf.getCurrentRelationship()

			// all the objects that will be loaded for this relationship, whatever their model
			allLoadedTargets := []IBusinessObject{}

			// if the relationship is single & directly persisted in the source business objects, then
			// we don't need to query for the targets' IDs, we already have them in the source business objects, so we can just read the targets directly
			if !relationship.IsMultiple() && relationship.isDirectlyPersisted() {
				// gathering all the targets of this relationship, for the given business objects
				for _, bObj := range bObjs {
					target, errGet := bObj.GetSingleRelationshipValue(relationship.GetName())
					if errGet != nil {
						return errGet
					}
					if target != nil {
						allLoadedTargets = append(allLoadedTargets, target)
					}
				}

			} else {

				// then dealing with all the business objects of the same model
				for modelName, bObjIDs := range allBObjIDs {

					// this relationship is either multi-valued, or not directly persisted (i.e. a backref),
					// so the relationship's *owner* DAO is the one that knows how to read it
					dao := bloCtx.daoFor(modelName)

					// retrieving the "empty" target business objects (ID - and model, for polymorphic targets - only)
					targets, errRead := dao.ExecReadRelationshipQuery(bObjIDs, relationship.GetName(), bloCtx.bObjCache())
					if errRead != nil {
						return ErrorC(errRead, "Could not read relationship '%s.%s'", relationship.owner.GetName(), relationship.GetName())
					}

					allLoadedTargets = append(allLoadedTargets, targets...)
				}
			}

			// now, we just have to "read" these target objects, and maybe some of their relationships as well
			if err := ReadBusinessObjects(bloCtx, subLoadingConf, allLoadedTargets...); err != nil {
				return err
			}
		}
	}

	// now, we can flag all the relationships as loaded for the source business objects
	for _, bObj := range bObjs {
		bObj.setLoaded(loadingConf.getLoadedRelationships())
	}

	return nil
}

// UpdateBusinessObjects updates existing business objects in the database
// WARNING: all the business objects must be of the same type (same model), this does not handle polymorphism
func UpdateBusinessObjects[BOTYPE IBusinessObject](bloCtx BloContext, bObjs ...BOTYPE) (updateErr error) {
	if len(bObjs) == 0 {
		return nil
	}

	// we know the BOs here are all of the same model
	model := bObjs[0].getModel(bObjs[0])

	// let's start a transaction if none is already started
	beginTransactionHere, errBegin := bloCtx.BeginTransaction(model)
	if errBegin != nil {
		return ErrorC(errBegin, "Could not update object since a transaction could not be started")
	}

	// let's make sure the transaction is always ended, even if a panic occurs below,
	// so that we never leave a dangling transaction / corrupt the BloContext's tx state
	if beginTransactionHere {
		bloCtx.Trace("New transaction started here!")
		defer func() {
			if r := recover(); r != nil {
				bloCtx.Trace("Recovered from panic while updating object(s): %v", r)
				errEnd := bloCtx.EndTransaction(Error("Recovered from panic while updating object(s): %v", r))
				panic(errEnd)
			}

			if errEnd := bloCtx.EndTransaction(updateErr); errEnd != nil {
				updateErr = ErrorC(errEnd, "Could not update object since the current transaction could not be terminated")
			} else {
				bloCtx.Trace("Transaction ended here!")
			}
		}()
	}

	// let's try to update the given entity
	updateErr = doUpdateBusinessObjects(bloCtx, model, bObjs...)

	return
}

// doUpdateBusinessObjects does the actual update of existing business objects in the database
func doUpdateBusinessObjects[BOTYPE IBusinessObject](bloCtx BloContext, model IBusinessObjectModel, bObjs ...BOTYPE) error {
	if len(bObjs) == 0 {
		return nil
	}

	// // we cannot update anything if we do not know who's doing it
	// if bloCtx.GetCurrentUser() == nil {
	// 	return NewErr("Could not update entity '%s' since the current user is unknown", ToEntityFullReference(updated))
	// }

	// preparing to retrieve the current versions of the BOs from the DB
	bObjsInDB, loadingConf, err := cloneBOsWithLoading(model, bObjs)
	if err != nil {
		return err
	}

	// actually reading the current versions of the BOs
	if err := ReadBusinessObjects(bloCtx, loadingConf, bObjsInDB...); err != nil {
		return err
	}

	// we need this type conversion later on, let's do it now since we're iterating over the given business objects anyway
	iBObjs := make([]IBusinessObject, len(bObjs))

	// now, doing a bit of work on the given BOs, thanks to what we just read from the DB
	for i, updated := range bObjs {
		// the current version of the BO in the DB
		instore := bObjsInDB[i]

		// // if it was deleted, we cannot update it from this method (to recover it, for instance), unless we're an admin
		// if instore.GetEntityStatus() == EntityStatusDELETED && !biContext.GetCurrentUser().IsAdmin() {
		// 	return NewErrC(err, "Could not update entity '%s' because it is deleted", ToEntityFullReference(updated))
		// }

		// we make sure the updated entity bears the right technical info, to avoid corruption from the outside
		updated.setCreation(instore.GetCreation())

		// can we update this entity ? does it need any changes before the update ?
		if err := updated.CheckAndChangeBeforeUpdate(bloCtx, instore); err != nil {
			return ErrorC(err, "Could not update entity '%s' because this update is not allowed", KeyFor(updated))
		}

		// check of validity constraints put on the corresponding schema
		if err := updated.IsModelValid(); err != nil {
			return ErrorC(err, "Could not update entity '%s' since it is not valid", KeyFor(updated))
		}

		// check of "functional / business" validity
		if err := updated.IsValid(bloCtx); err != nil {
			return ErrorC(err, "Could not update entity '%s' since it is not valid", KeyFor(updated))
		}

		// updating some tracking info
		// updated.SetModifiedByID(biContext.GetCurrentUser().GetID())
		// updated.SetModifiedBy(biContext.GetCurrentUser().GetLabel())
		updated.setModification(core.Now())
		// updated.SetEntityStatus(EntityStatusUPDATED)

		// // updating the hash, if required
		// if updated.IsHashed() {
		// 	updated.setHash(GetHash(updated))
		// }

		// adding the business object to the interface slice
		iBObjs[i] = updated
	}

	// we need the DAO for the given BO type here
	dao := bloCtx.daoFor(model.GetName())

	// pushing the business objects into the DB !
	if err := dao.ExecUpdateQuery(iBObjs...); err != nil {
		return ErrorC(err, "Could not update the '%s' instances because of a problem with the DB", model.GetName())
	}

	// now, handling the links
	if err := doUpdateBusinessObjectLinks(dao, bObjs, bObjsInDB, loadingConf); err != nil {
		return ErrorC(err, "Could not update the '%s' instances' links because of a problem with the DB", model.GetName())
	}

	// now doing some work on the given BOs, after the update
	for _, bObj := range bObjs {
		// we have stuff to do after the update ? yeah ? really ? let's do it now !
		if err := bObj.ChangeAfterUpdate(bloCtx); err != nil {
			return ErrorC(err, "Could not post-update '%s' because it got an error", KeyFor(bObj))
		}

		// // tracking the diffs brought to the entity through this update
		// if updated.IsDiffed() {
		// 	go trackDiffs(biContext, instore, updated)
		// }
	}

	// no problem updating, so returning
	return nil

}

// doUpdateBusinessObjectLinks updates the links for the given business objects by computing the differences
// with their counterparts in the database and executing the necessary add and remove link queries.
func doUpdateBusinessObjectLinks[BOTYPE IBusinessObject](dao IBusinessObjectDAO, bObjs []BOTYPE, bObjsInDB []BOTYPE, loadingConf ILoadingConfig) error {
	// synthetic BOs bearing links to add to / remove from the DB
	bObjsWithLinksToAdd := make([]IBusinessObject, 0)
	bObjsWithLinksToRemove := make([]IBusinessObject, 0)

	// a map to flag which links to operate on
	forLinks := map[string]bool{}
	for _, loadedLink := range loadingConf.getLoadedRelationships() {
		forLinks[string(loadedLink)] = true
	}

	// handling the links for each business object
	for i, bObj := range bObjs {
		// computing the links to add and remove, and making them borne by synthetic BOs
		boWithLinksToAdd, boWithLinksToRemove := bObj.DiffWith(bObjsInDB[i], forLinks)

		// appending the synthetic BOs to the respective slices
		bObjsWithLinksToAdd = append(bObjsWithLinksToAdd, boWithLinksToAdd)
		bObjsWithLinksToRemove = append(bObjsWithLinksToRemove, boWithLinksToRemove)
	}

	// removing the links
	if err := dao.ExecDeleteLinksQueries(bObjsWithLinksToRemove...); err != nil {
		return ErrorC(err, "Could not remove the links for the business objects")
	}

	// adding the links
	if err := dao.ExecCreateLinksQueries(bObjsWithLinksToAdd...); err != nil {
		return ErrorC(err, "Could not add the links for the business objects")
	}

	return nil
}
