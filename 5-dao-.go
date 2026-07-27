package goald

import (
	"database/sql"

	"github.com/aldesgroup/goald/features/logging"
)

// ------------------------------------------------------------------------------------------------
// DAO interface definition
// ------------------------------------------------------------------------------------------------

type IBusinessObjectDAO interface {
	// a DAO is a complete context that should allow to perform queries for a specific business object class
	logging.ILogger                   // being able to log from the DAO
	setLogger(logger logging.ILogger) // setting the logger to use for this DAO
	getClassName() className          // the business object class name this DAO is for
	setClassName(clsName className)   // setting the business object class name this DAO is for
	getDB() *DB                       // accessing the DB object associated with this DAO
	setDB(db *DB)                     // setting the DB object associated with this DAO
	getTx() *sql.Tx                   // accessing the current transaction associated with this DAO
	setTx(tx *sql.Tx)                 // setting the current transaction associated with this DAO
	NewDAO() IBusinessObjectDAO       // returns a new DAO of the same type

	// proxy functions to the DB object associated with this DAO
	Query(m mask, query string, args ...any) (*sql.Rows, error) // proxying the DB.Query function
	Exec(m mask, query string, args ...any) (sql.Result, error) // proxying the DB.Exec function
	QueryRow(m mask, query string, args ...any) *sql.Row        // proxying the DB.QueryRow function

	// these are the generic queries any DAO should be able to perform for its associated business object class
	ExecInsertQuery(bObjs ...IBusinessObject) (map[int]int64, error) // executes the insert query for the given BOs, returning a map of the BO's rowID to the DB's ID
	ExecInsertLinksQueries(bObjs ...IBusinessObject) error           // executes the insert links queries for the given business objects
}

// ------------------------------------------------------------------------------------------------
// Base DAO implementation
// ------------------------------------------------------------------------------------------------

type BusinessObjectDAO struct {
	logging.ILogger
	clsName className
	db      *DB
	tx      *sql.Tx
}

// NewDAO implements [IBusinessObjectDAO].
func (baseDAO *BusinessObjectDAO) NewDAO() IBusinessObjectDAO {
	panic("unimplemented")
}

func (baseDAO *BusinessObjectDAO) setLogger(logger logging.ILogger) {
	baseDAO.ILogger = logger
}

func (baseDAO *BusinessObjectDAO) getClassName() className {
	return baseDAO.clsName
}

func (baseDAO *BusinessObjectDAO) setClassName(clsName className) {
	baseDAO.clsName = clsName
}

func (baseDAO *BusinessObjectDAO) getDB() *DB {
	return baseDAO.db
}

func (baseDAO *BusinessObjectDAO) setDB(db *DB) {
	baseDAO.db = db
}

func (baseDAO *BusinessObjectDAO) getTx() *sql.Tx {
	return baseDAO.tx
}

func (baseDAO *BusinessObjectDAO) setTx(tx *sql.Tx) {
	baseDAO.tx = tx
}

// ------------------------------------------------------------------------------------------------
// Proxy functions
// ------------------------------------------------------------------------------------------------

func (baseDAO *BusinessObjectDAO) Query(m mask, query string, args ...any) (*sql.Rows, error) {
	return baseDAO.db.query(baseDAO, baseDAO.tx, m, query, args...)
}

func (baseDAO *BusinessObjectDAO) Exec(m mask, query string, args ...any) (sql.Result, error) {
	return baseDAO.db.exec(baseDAO, baseDAO.tx, m, query, args...)
}

func (baseDAO *BusinessObjectDAO) QueryRow(m mask, query string, args ...any) *sql.Row {
	return baseDAO.db.queryRow(baseDAO, baseDAO.tx, m, query, args...)
}

// ------------------------------------------------------------------------------------------------
// Default (empty) implementation
// ------------------------------------------------------------------------------------------------

// type check
var _ IBusinessObjectDAO = (*BusinessObjectDAO)(nil)

// ExecInsertQuery implements [IBusinessObjectDAO].
func (baseDAO *BusinessObjectDAO) ExecInsertQuery(bObjs ...IBusinessObject) (map[int]int64, error) {
	panic("unimplemented")
}

// ExecInsertLinksQueries implements [IBusinessObjectDAO].
func (baseDAO *BusinessObjectDAO) ExecInsertLinksQueries(bObjs ...IBusinessObject) error {
	return nil
}
