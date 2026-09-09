// ------------------------------------------------------------------------------------------------
// BloContext is a context that should provide the necessary info for Business LOgic processing
// ------------------------------------------------------------------------------------------------

package goald

import (
	"database/sql"
	"sync"

	"github.com/aldesgroup/goald/features/utils"
)

type BloContext interface {
	AppContext                                                 // a particular AppContext dedicated to Business LOgic processing
	BeginTransaction(model IBusinessObjectModel) (bool, error) // starts a new transaction if none is already started; returns true if a new transaction was started, false if there was already one
	EndTransaction(err error) error                            // ends the current transaction if it was started by this BloContext, and commits or rollbacks depending on the given error
	daoFor(modelName utils.ModelName) IBusinessObjectDAO       // returns a new DAO from a given business object
	bObjCache() *BObjCache                                     // returns a cache of business objects associated with this context
}

// ------------------------------------------------------------------------------------------------
// Base BLO context
// ------------------------------------------------------------------------------------------------

type baseBloContextImpl struct {
	currentTx          *sql.Tx    // the transaction currently bearing all the changes that we want to bring to the DB
	currentTxClientsNb int        // the current number of "clients" currently relyong on the current transaction
	mx                 sync.Mutex // safely handling the changes on the current transaction
	boCache            *BObjCache // the cache of business objects associated with this context
}

// BeginTransaction implements [BloContext].
func (baseBloCtx *baseBloContextImpl) BeginTransaction(model IBusinessObjectModel) (newTxStarted bool, err error) {
	// not messing around with concurrent routines
	baseBloCtx.mx.Lock()
	defer baseBloCtx.mx.Unlock()

	// in any case, there's one more client
	baseBloCtx.currentTxClientsNb++

	// we create a new transaction only if there is none yet
	if baseBloCtx.currentTx == nil {
		// let's start a transaction at the DB level
		baseBloCtx.currentTx, err = model.getDB().do.Begin()

		// handling a potential error
		if err != nil {
			return false, ErrorC(err, "Could not begin BI transaction")
		}

		newTxStarted = true
		return
	}

	// no new transaction, and no error
	return
}

// EndTransaction implements [BloContext].
func (baseBloCtx *baseBloContextImpl) EndTransaction(previousErr error) error {
	// not messing around with concurrent routines
	baseBloCtx.mx.Lock()
	defer baseBloCtx.mx.Unlock()

	// are we trying to end a transaction that does not exist ?
	if baseBloCtx.currentTx == nil {
		return Error("Could not end a BI transaction, since there is none in the context")
	}

	// ok now we're sure there's a transaction going on ! Whatever happens, we try to remove it at the end
	defer func() {
		// in any case, there's one less client
		baseBloCtx.currentTxClientsNb--

		// if there's no client left that uses the current transaction let's remove it
		if baseBloCtx.currentTxClientsNb == 0 {
			baseBloCtx.currentTx = nil
		}
	}()

	// if at least one of the operations that occurred during this transaction has failed, we have to rollback
	if previousErr != nil {
		if errRb := baseBloCtx.currentTx.Rollback(); errRb != nil {
			return ErrorC(errRb, "Have to roll back current transaction, because one operation during it has failed (%s),"+
				" but rolling back the transaction has also failed", previousErr)
		}

		return ErrorC(previousErr, "Have to roll back current transaction, because one operation during it has failed")
	}

	// we can commit the transaction only if we are the only client (left) of it
	if baseBloCtx.currentTxClientsNb == 1 {
		if txErr := baseBloCtx.currentTx.Commit(); txErr != nil {
			// rolling the transaction back since we could not commit it
			if errRb := baseBloCtx.currentTx.Rollback(); errRb != nil {
				return ErrorC(errRb, "Have to roll back current transaction, since committing it has failed, but rolling back the transaction has also failed")
			}

			return ErrorC(txErr, "Have to roll back current transaction, since committing it has failed")
		}
	}

	return nil
}

// bObjCache implements [BloContext].
func (baseBloCtx *baseBloContextImpl) bObjCache() *BObjCache {
	if baseBloCtx.boCache == nil {
		baseBloCtx.boCache = NewBObjCache()
	}
	return baseBloCtx.boCache
}

// ------------------------------------------------------------------------------------------------
// BLO context associated with a HTTP request
// ------------------------------------------------------------------------------------------------

// default implementation for business logic context
type httpBloContextImpl struct {
	restContext         // inheriting the REST context
	*baseBloContextImpl // making this a base business logic context at least
	// *httpRequestContext // wrapping one of the server's children handling 1 request
}

// type check
var _ BloContext = (*httpBloContextImpl)(nil)

// constructors
func newHttpBloContextFromWebCtx(thisWebCtx *webContextImpl) *httpBloContextImpl {
	return &httpBloContextImpl{
		restContext: thisWebCtx,
		// httpRequestContext: thisWebCtx.httpRequestContext,
		baseBloContextImpl: &baseBloContextImpl{},
	}
}

// daoFor implements [BloContext].
func (httpBloCtx *httpBloContextImpl) daoFor(modelName utils.ModelName) IBusinessObjectDAO {
	newDAO := newDaoFor(modelName)
	newDAO.setLogger(httpBloCtx)
	newDAO.setModel(modelFor(modelName, true))
	newDAO.setTx(httpBloCtx.currentTx)
	return newDAO
}

// ------------------------------------------------------------------------------------------------
// Utils - business object cache, used for consistency and efficiency in reading objects
// ------------------------------------------------------------------------------------------------

// BObjCache is a cache for business objects, keyed by their model name and ID.
type BObjCache struct {
	content map[utils.ModelName]map[BObjID]IBusinessObject
}

func NewBObjCache() *BObjCache {
	return &BObjCache{
		content: make(map[utils.ModelName]map[BObjID]IBusinessObject),
	}
}

// Get retrieves a business object from the cache by its model name and ID. Returns nil if not found.
func (thisCache *BObjCache) Get(modelName utils.ModelName, id BObjID) IBusinessObject {
	if _, modelExists := thisCache.content[modelName]; !modelExists {
		return nil
	}

	if _, objectExists := thisCache.content[modelName][id]; !objectExists {
		return nil
	}

	return thisCache.content[modelName][id]
}

// CachedOrNewBusinessObject retrieves a business object from the cache if it exists, or creates a new one and adds it to the cache if it doesn't.
func (thisCache *BObjCache) CachedOrNewBusinessObject(modelName utils.ModelName, id BObjID) IBusinessObject {
	if cachedBObj := thisCache.Get(modelName, id); cachedBObj != nil {
		return cachedBObj
	}

	return thisCache.Set(NewBusinessObject(modelName, id))
}

// CachedOrNewBusinessObjectRaw retrieves a business object from the cache if it exists, or creates a new one
// and adds it to the cache if it doesn't, using raw string and int64 types for the model name and ID.
func (thisCache *BObjCache) CachedOrNewBusinessObjectRaw(modelName string, id int64) IBusinessObject {
	return thisCache.CachedOrNewBusinessObject(utils.ModelName(modelName), BObjID(id))
}

// GetObj retrieves a business object from the cache using the business object itself to determine the model name and ID. Returns nil if not found.
func (thisCache *BObjCache) GetObj(bObj IBusinessObject) IBusinessObject {
	if bObj == nil {
		return nil
	}

	return thisCache.Get(bObj.GetModelName(), bObj.GetID())
}

// GetObjsFor retrieves all business objects of a specific model from the cache. Returns a map of IDs to business objects.
func (thisCache *BObjCache) GetObjsFor(modelName utils.ModelName) map[BObjID]IBusinessObject {
	return thisCache.content[modelName]
}

// Set adds a business object to the cache. If an object with the same model name and ID already exists, it returns the zero value of the type.
func (thisCache *BObjCache) Set[BOTYPE IBusinessObject](bObj BOTYPE) BOTYPE {
	if thisCache.content[bObj.GetModelName()] == nil {
		thisCache.content[bObj.GetModelName()] = map[BObjID]IBusinessObject{}
	}

	if _, alreadyThere := thisCache.content[bObj.GetModelName()][bObj.GetID()]; alreadyThere {
		var zero BOTYPE
		return zero
	}

	thisCache.content[bObj.GetModelName()][bObj.GetID()] = bObj

	return bObj
}

// AddAll adds all the given business objects to the cache and returns a map of all their IDs grouped by model name.
func (thisCache *BObjCache) AddAll(bObjs []IBusinessObject) (allBObjIDs map[utils.ModelName][]any) {
	allBObjIDs = make(map[utils.ModelName][]any)

	// iterating to build the result
	for _, bObj := range bObjs {
		// caching the object
		thisCache.Set(bObj)

		// gathering its ID along with the IDs of the objects of the same model
		allBObjIDs[bObj.GetModelName()] = append(allBObjIDs[bObj.GetModelName()], bObj.GetID())
	}

	return
}
