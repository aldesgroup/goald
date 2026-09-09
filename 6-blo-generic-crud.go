// ------------------------------------------------------------------------------------------------
// Here we implement the generic business logic involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"slices"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

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
		bObj.setID(BObjID(rowToIDMap[bObj.GetPreID()]))
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
func ReadBusinessObjects(bloCtx BloContext, loadingConf ILoadingConfig, bObjs ...IBusinessObject) error {
	// logObjs(bloCtx, bObjs, "Reading business objects with loading config: %s", loadingConf.ToString())

	// first, caching the business objects, whilst computing the IDs mapped by their model
	allBObjIDs := bloCtx.bObjCache().AddAll(bObjs)

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
func afterLoad(bloCtx BloContext, loadingConf ILoadingConfig, bObjs ...IBusinessObject) error {
	// we check that this entity can be read, given the context
	for _, bObj := range bObjs {
		if err := bObj.CanBeRead(bloCtx); err != nil {
			return err
		}
	}

	// now that we're okey with our loading, me might have to perform some specific actions
	for _, bObj := range bObjs {
		if err := bObj.ChangeAfterRead(bloCtx); err != nil {
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
func ReadRelationships(bloCtx BloContext, loadingConf ILoadingConfig, bObjs ...IBusinessObject) error {
	// logObjs(bloCtx, bObjs, "Reading relationships with loading config: %s", loadingConf.ToString())

	if len(loadingConf.getRelationshipsToLoad()) > 0 {
		// first, caching the business objects, whilst computing the IDs mapped by their model
		allBObjIDs := bloCtx.bObjCache().AddAll(bObjs)

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

	return nil
}

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

// ------------------------------------------------------------------------------------------------
// Loading configs
// ------------------------------------------------------------------------------------------------

type ILoadingConfig interface {
	withParent(parent ILoadingConfig) ILoadingConfig
	getRelationshipsToLoad() []ILoadingConfig
	getCurrentRelationship() *Relationship
	ToString(indent ...int) string
}

type loadingConfig struct {
	relationshipsOwner  IBusinessObjectModel
	parent              ILoadingConfig
	relationshipsToLoad []ILoadingConfig
	currentRelationship *Relationship
}

func (cfg *loadingConfig) withParent(parent ILoadingConfig) ILoadingConfig {
	cfg.parent = parent
	return cfg
}

func (cfg *loadingConfig) getRelationshipsToLoad() []ILoadingConfig {
	return cfg.relationshipsToLoad
}

func (cfg *loadingConfig) getCurrentRelationship() *Relationship {
	return cfg.currentRelationship
}

func (cfg *loadingConfig) ToString(indent ...int) string {
	ind := 0
	if len(indent) > 0 {
		ind = indent[0]
	}
	prefix := strings.Repeat("  ", ind)

	result := ""

	if cfg.parent == nil {
		result = prefix + fmt.Sprintf("loading '%s' with:", cfg.relationshipsOwner.GetName())
	} else if ind == 0 {
		result = prefix + fmt.Sprintf("loading '%s' (from '%s#%s')%s",
			strings.Join(core.ToStrings(cfg.currentRelationship.getTargetModelNames()), ", "),
			cfg.relationshipsOwner.GetName(),
			cfg.currentRelationship.name,
			core.IfThenElse(len(cfg.relationshipsToLoad) > 0, " with:", "."),
		)
	}

	for _, subConfig := range cfg.relationshipsToLoad {
		result += "\n" + prefix + fmt.Sprintf("  - %s (%v)", subConfig.getCurrentRelationship().name, subConfig.getCurrentRelationship().getTargetModelNames())
		result += subConfig.ToString(ind + 1)
	}

	return result
}

func Load(model IBusinessObjectModel, with ...ILoadingConfig) ILoadingConfig {
	thisConfig := &loadingConfig{
		relationshipsOwner: model,
	}
	for _, subConfig := range with {
		thisConfig.relationshipsToLoad = append(thisConfig.relationshipsToLoad, subConfig.withParent(thisConfig))
	}
	return thisConfig
}

func With(relationship *Relationship, with ...ILoadingConfig) ILoadingConfig {
	thisConfig := &loadingConfig{
		relationshipsOwner:  relationship.owner,
		currentRelationship: relationship,
	}
	for _, subConfig := range with {
		thisConfig.relationshipsToLoad = append(thisConfig.relationshipsToLoad, subConfig.withParent(thisConfig))
	}
	return thisConfig
}

// ReadNoRelationship returns a loading config that loads only the direct relationships of the given business object type
func (thisModel *businessObjectModel) ReadNoRelationship() ILoadingConfig {
	return Load(thisModel)
}

// ReadWithFirstLayer returns a loading config that loads only the direct relationships of the given business object type
func (thisModel *businessObjectModel) ReadWithFirstLayer() ILoadingConfig {
	firstLayer := []ILoadingConfig{}
	for _, relationship := range core.GetSortedValues(thisModel.getRelationships()) {
		firstLayer = append(firstLayer, With(relationship))
	}
	return Load(thisModel, firstLayer...)
}

// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

func logObjs(bloCtx BloContext, bObjs []IBusinessObject, reasonForLogging string, args ...any) {
	if bloCtx.IsTraceEnabled() {
		objStrings := fmt.Sprintf("[%s]", strings.Join(core.MapFn(bObjs, func(bObj IBusinessObject) string {
			return fmt.Sprintf("%s (%p)", KeyFor(bObj), bObj)
		}), ", "))
		bloCtx.Trace("----------------------------------------------")
		bloCtx.Trace(fmt.Sprintf(reasonForLogging+": "+objStrings, args...))
		bloCtx.Trace("----------------------------------------------")
	}
}

func logIDs(bloCtx BloContext, bObjIDs []any, reasonForLogging string, args ...any) {
	if bloCtx.IsTraceEnabled() {
		objStrings := fmt.Sprintf("[%s]", strings.Join(core.MapFn(bObjIDs, func(id any) string {
			return fmt.Sprintf("%d", id)
		}), ", "))
		bloCtx.Trace("----------------------------------------------")
		bloCtx.Trace(fmt.Sprintf(reasonForLogging+": "+objStrings, args...))
		bloCtx.Trace("----------------------------------------------")
	}
}
