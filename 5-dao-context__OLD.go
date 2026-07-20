package goald

// ------------------------------------------------------------------------------------------------
// DaoContext should contain the necessary info for handling database access
// ------------------------------------------------------------------------------------------------

// type DaoContext interface {
// 	logging.ILogger
// 	// tx() *sql.Tx                                        // returns the current SQL transaction
// 	// setTx(tx *sql.Tx)                                   // sets the current SQL transaction
// 	// delTx()                                             // removes the current SQL transaction
// 	// canCommitTx() bool                                  // tells us if the current transaction can be committed
// 	// BeginDbTransaction() (*sql.Tx, error)               // starts a new SQL transaction if none is already started
// 	// Exec(query string, args ...any) (sql.Result, error) // executes a query within the current transaction if any
// }

// -----------------------------------------------------------------------------
// Base (incomplete) implementation of DbContext
// -----------------------------------------------------------------------------

// // daoContextImpl is a struct that contains some info and methods to implement BiContext
// type daoContextImpl struct {
// 	logging.ILogger
// 	currentTx          *sql.Tx     // the transaction currently bearing all the changes that we want to bring to the DB
// 	currentTxClientsNb int         // the current number of "clients" currently relyong on the current transaction
// 	mx                 *sync.Mutex // safely handling the changes on the current transaction
// 	db                 *DB         // the DB connection to use for this context
// }

// var _ DaoContext = (*daoContextImpl)(nil) // type check

// // constructors
// func newDaoContextFromHttpBloCtx(thisBloCtx *httpBloContextImpl) DaoContext {
// 	return &daoContextImpl{
// 		ILogger: thisBloCtx,
// 		mx:      &sync.Mutex{},
// 		db:      modelForName(thisBloCtx.getTargetResourceClass()).getInDB(),
// 	}
// }

// // implementing DaoContext
// func (thisDbContext *daoContextImpl) tx() *sql.Tx {
// 	thisDbContext.mx.Lock()
// 	defer thisDbContext.mx.Unlock()

// 	return thisDbContext.currentTx
// }

// // setTx sets the current SQL transaction
// func (thisDbContext *daoContextImpl) setTx(tx *sql.Tx) {
// 	// not messing around with concurrent routines: START
// 	thisDbContext.mx.Lock()

// 	// if there's no client yet that has set the transaction let's do it
// 	if thisDbContext.currentTxClientsNb == 0 {
// 		thisDbContext.currentTx = tx
// 	}

// 	// in any case, there's one more client
// 	thisDbContext.currentTxClientsNb++

// 	// not messing around with concurrent routines: END
// 	thisDbContext.mx.Unlock()
// }

// // delTx removes the current SQL transaction
// func (thisDbContext *daoContextImpl) delTx() {
// 	// not messing around with concurrent routines: START
// 	thisDbContext.mx.Lock()

// 	// in any case, there's one less client
// 	thisDbContext.currentTxClientsNb--

// 	// if there's no client left that uses the current transaction let's remove it
// 	if thisDbContext.currentTxClientsNb == 0 {
// 		thisDbContext.currentTx = nil
// 	}

// 	// not messing around with concurrent routines: END
// 	thisDbContext.mx.Unlock()
// }

// // tells us if the current transaction can be committed
// func (thisDbContext *daoContextImpl) canCommitTx() bool {
// 	thisDbContext.mx.Lock()
// 	defer thisDbContext.mx.Unlock()

// 	return thisDbContext.currentTxClientsNb == 1 // only the setter of the transaction can commit
// }

// // BeginDbTransaction implements [DaoContext].
// func (thisDbContext *daoContextImpl) BeginDbTransaction() (*sql.Tx, error) {
// 	if thisDbContext.db == nil {
// 		return nil, ErrorC(nil, "No DB connection available for this context")
// 	}

// 	return thisDbContext.db.do.Begin()
// }

// func (thisDbContext *daoContextImpl) Exec(query string, args ...any) (sql.Result, error) {
// 	if thisDbContext.db == nil {
// 		return nil, ErrorC(nil, "No DB connection available for this context")
// 	}

// 	return thisDbContext.db.exec(thisDbContext, query, args...)
// }

// // InsertReturningID implements [DaoContext]; it hides the DB-specific way of retrieving a freshly inserted row's ID
// func (thisDbContext *daoContextImpl) insertReturningID(query string, args ...any) (int64, error) {
// 	if thisDbContext.db == nil {
// 		return 0, ErrorC(nil, "No DB connection available for this context")
// 	}

// 	// e.g. Postgres: appending a 'RETURNING id' clause, and reading the new ID off of the returned row
// 	if thisDbContext.db.get.SupportsReturningID() {
// 		var id int64
// 		if errScan := thisDbContext.db.queryRow(thisDbContext, query+" RETURNING id", args...).Scan(&id); errScan != nil {
// 			return 0, errScan
// 		}

// 		return id, nil
// 	}

// 	// e.g. MySQL: no 'RETURNING' support, so falling back to the driver's LastInsertId()
// 	res, errExec := thisDbContext.db.exec(thisDbContext, query, args...)
// 	if errExec != nil {
// 		return 0, errExec
// 	}

// 	return res.LastInsertId()
// }
