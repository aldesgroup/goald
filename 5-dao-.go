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
	getDB() *DB                       // accessing the DB object associated with this DAO
	setDB(db *DB)                     // setting the DB object associated with this DAO
	NewDAO() IBusinessObjectDAO       // returns a new DAO of the same type

	// proxy functions to the DB object associated with this DAO
	Query(m mask, query string, args ...any) (*sql.Rows, error) // proxying the DB.Query function
	Exec(m mask, query string, args ...any) (sql.Result, error) // proxying the DB.Exec function
	QueryRow(m mask, query string, args ...any) *sql.Row        // proxying the DB.QueryRow function

	// these are the generic queries any DAO should be able to perform for its associated business object class
	ExecInsertQuery(bObj IBusinessObject) (int64, error) // executes the insert query for a given business object and returns its new ID
}

// ------------------------------------------------------------------------------------------------
// Base DAO implementation
// ------------------------------------------------------------------------------------------------

type BusinessObjectDAO struct {
	logging.ILogger
	db *DB
}

func (baseDAO *BusinessObjectDAO) setLogger(logger logging.ILogger) {
	baseDAO.ILogger = logger
}

func (baseDAO *BusinessObjectDAO) getDB() *DB {
	return baseDAO.db
}

func (baseDAO *BusinessObjectDAO) setDB(db *DB) {
	baseDAO.db = db
}

// ------------------------------------------------------------------------------------------------
// Proxy functions
// ------------------------------------------------------------------------------------------------

func (baseDAO *BusinessObjectDAO) Query(m mask, query string, args ...any) (*sql.Rows, error) {
	return baseDAO.db.query(baseDAO, m, query, args...)
}

func (baseDAO *BusinessObjectDAO) Exec(m mask, query string, args ...any) (sql.Result, error) {
	return baseDAO.db.exec(baseDAO, m, query, args...)
}

func (baseDAO *BusinessObjectDAO) QueryRow(m mask, query string, args ...any) *sql.Row {
	return baseDAO.db.queryRow(baseDAO, m, query, args...)
}
