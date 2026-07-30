package goald

import (
	"database/sql"

	"github.com/aldesgroup/goald/features/logging"
)

// ------------------------------------------------------------------------------------------------
// DAO interface definition
// ------------------------------------------------------------------------------------------------

type IBusinessObjectDAO interface {
	// a DAO is a complete context that should allow to perform queries for a specific business object ModelName
	logging.ILogger                      // being able to log from the DAO
	setLogger(logger logging.ILogger)    // setting the logger to use for this DAO
	getModel() IBusinessObjectModel      // the business object ModelName name this DAO is for
	setModel(model IBusinessObjectModel) // setting the business object ModelName name this DAO is for
	getTx() *sql.Tx                      // accessing the current transaction associated with this DAO
	setTx(tx *sql.Tx)                    // setting the current transaction associated with this DAO
	NewDAO() IBusinessObjectDAO          // returns a new DAO of the same type

	// proxy functions to the DB object associated with this DAO
	Query(m mask, query string, args ...any) (*sql.Rows, error) // proxying the DB.Query function
	Exec(m mask, query string, args ...any) (sql.Result, error) // proxying the DB.Exec function
	QueryRow(m mask, query string, args ...any) *sql.Row        // proxying the DB.QueryRow function

	// these are the generic queries any DAO should be able to perform for its associated business object ModelName
	ExecInsertQuery(bObjs ...IBusinessObject) (map[int]int64, error) // executes the insert query for the given BOs, returning a map of the BO's rowID to the DB's ID
	ExecInsertLinksQueries(bObjs ...IBusinessObject) error           // executes the insert links queries for the given business objects
}

// ------------------------------------------------------------------------------------------------
// Base DAO implementation
// ------------------------------------------------------------------------------------------------

type BusinessObjectDAO struct {
	logging.ILogger
	model IBusinessObjectModel
	tx    *sql.Tx
}

// NewDAO implements [IBusinessObjectDAO].
func (baseDAO *BusinessObjectDAO) NewDAO() IBusinessObjectDAO {
	panic("unimplemented")
}

func (baseDAO *BusinessObjectDAO) setLogger(logger logging.ILogger) {
	baseDAO.ILogger = logger
}

func (baseDAO *BusinessObjectDAO) getModel() IBusinessObjectModel {
	return baseDAO.model
}

func (baseDAO *BusinessObjectDAO) setModel(model IBusinessObjectModel) {
	baseDAO.model = model
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
	return baseDAO.model.getDB().query(baseDAO, baseDAO.tx, m, query, args...)
}

func (baseDAO *BusinessObjectDAO) Exec(m mask, query string, args ...any) (sql.Result, error) {
	return baseDAO.model.getDB().exec(baseDAO, baseDAO.tx, m, query, args...)
}

func (baseDAO *BusinessObjectDAO) QueryRow(m mask, query string, args ...any) *sql.Row {
	return baseDAO.model.getDB().queryRow(baseDAO, baseDAO.tx, m, query, args...)
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
	panic("unimplemented")
}
