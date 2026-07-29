// ------------------------------------------------------------------------------------------------
// BloContext is a context that should provide the necessary info for Business LOgic processing
// ------------------------------------------------------------------------------------------------

package goald

import (
	"database/sql"
	"sync"
)

type BloContext interface {
	AppContext                                   // a particular AppContext dedicated to Business LOgic processing
	BeginTransaction(class IClass) (bool, error) // starts a new transaction if none is already started; returns true if a new transaction was started, false if there was already one
	EndTransaction(err error) error              // ends the current transaction if it was started by this BloContext, and commits or rollbacks depending on the given error
	daoFor(class IClass) IBusinessObjectDAO      // returns a new DAO from a given business object
}

// ------------------------------------------------------------------------------------------------
// Base BLO context
// ------------------------------------------------------------------------------------------------

type baseBloContextImpl struct {
	currentTx          *sql.Tx    // the transaction currently bearing all the changes that we want to bring to the DB
	currentTxClientsNb int        // the current number of "clients" currently relyong on the current transaction
	mx                 sync.Mutex // safely handling the changes on the current transaction
}

// BeginTransaction implements [BloContext].
func (baseBloCtx *baseBloContextImpl) BeginTransaction(class IClass) (newTxStarted bool, err error) {
	// not messing around with concurrent routines
	baseBloCtx.mx.Lock()
	defer baseBloCtx.mx.Unlock()

	// in any case, there's one more client
	baseBloCtx.currentTxClientsNb++

	// we create a new transaction only if there is none yet
	if baseBloCtx.currentTx == nil {
		// let's start a transaction at the DB level
		baseBloCtx.currentTx, err = class.getModel().getDB().do.Begin()

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
func (httpBloCtx *httpBloContextImpl) daoFor(class IClass) IBusinessObjectDAO {
	newDAO := newDaoFor(class)
	newDAO.setLogger(httpBloCtx)
	newDAO.setClass(class)
	newDAO.setDB(class.getModel().getDB())
	newDAO.setTx(httpBloCtx.currentTx)
	return newDAO
}
