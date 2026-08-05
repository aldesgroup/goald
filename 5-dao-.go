package goald

import (
	"database/sql"
	"strings"

	core "github.com/aldesgroup/corego"
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

	// proxy functions to the DB adapter associated with this DAO's business object model
	HandleDbError(err error) error

	// these are the generic queries any DAO should be able to perform for its associated business object ModelName
	ExecCreateQuery(bObjs ...IBusinessObject) (map[int]int64, error)                           // executes the insert query for the given BOs, returning a map of the BO's rowID to the DB's ID
	ExecCreateLinksQueries(bObjs ...IBusinessObject) error                                     // executes the insert links queries for the given business objects
	ExecSearchQuery(queryName queryName, values ISearchParamValues) ([]IBusinessObject, error) // executes the select query for the given query name and values, returning the resulting business objects
	ExecReadQuery(bObjs map[BObjID]IBusinessObject, bObjIDs []any) error                       // executes the read query for the given business objects
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

// ExecCreateQuery implements [IBusinessObjectDAO].
func (baseDAO *BusinessObjectDAO) ExecCreateQuery(bObjs ...IBusinessObject) (map[int]int64, error) {
	panic("unimplemented")
}

// ExecCreateLinksQueries implements [IBusinessObjectDAO].
func (baseDAO *BusinessObjectDAO) ExecCreateLinksQueries(bObjs ...IBusinessObject) error {
	panic("unimplemented")
}

// ExecSearchQuery implements [IBusinessObjectDAO].
func (baseDAO *BusinessObjectDAO) ExecSearchQuery(queryName queryName, values ISearchParamValues) ([]IBusinessObject, error) {
	panic("unimplemented")
}

// ExecReadQuery implements [IBusinessObjectDAO].
func (baseDAO *BusinessObjectDAO) ExecReadQuery(bObjs map[BObjID]IBusinessObject, bObjIDs []any) error {
	panic("unimplemented")
}

// ------------------------------------------------------------------------------------------------
// DB error handling
// ------------------------------------------------------------------------------------------------

func (baseDAO *BusinessObjectDAO) HandleDbError(err error) error {
	// using the DB adapter to parse the error, so that we can handle it in a DB-agnostic way
	dbError, content := baseDAO.model.getDB().get.ParseDbError(err)

	switch dbError {

	case DbErrorUNIDENTIFIED:
		return err

	case DbErrorDUPLICATExENTRY:
		if strings.HasPrefix(content, prefixUK) {
			if parts := strings.Split(content[len(prefixUK):], "__"); len(parts) == 2 {
				modelName := core.SnakeToPascal(parts[0])
				relationshipName := core.SnakeToCamel(parts[1])
				return Error("%s: there's already a '%s' instance with the same '%s' value!", dbError, modelName, relationshipName)
			}
		}
		return ErrorC(err, "Duplicate entry detected!")

	case DbErrorMISSINGxVALUE:
		if parts := strings.Split(content, "."); len(parts) == 2 {
			modelName := core.SnakeToPascal(parts[0])
			columnName := parts[1]
			if strings.HasSuffix(parts[1], suffixID) {
				columnName = core.SnakeToCamel(core.Before(columnName, suffixID))
			} else {
				columnName = core.SnakeToCamel(columnName)
			}
			return Error("%s: a '%s' instance cannot have its '%s' missing!", dbError, modelName, columnName)
		}
		return ErrorC(err, "Missing value detected!")

	case DbErrorINVALIDxENTRY:
		if strings.HasPrefix(content, prefixFK) {
			if parts := strings.Split(content[len(prefixFK):], "__"); len(parts) == 3 {
				modelName := core.SnakeToPascal(parts[0])
				relationshipName := core.SnakeToCamel(parts[1])
				return Error("%s: this '%s' relationship associated with this '%s' instance does not seem to exist in the database!", dbError,
					relationshipName, modelName)
			}
		}
		return ErrorC(err, "Invalid entry detected!")

	default:
		return ErrorC(err, "Unhandled DB error type (%d): %s", dbError, content)
	}
}

// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

// GetIDOrNilPlm returns nil if bo is nil, otherwise bo.GetID() - use this for relationship fields that
// are already typed as an interface (polymorphic relationships, e.g. domain.IContact): a nil value
// there stays a plain nil interface once passed as an IBusinessObject, so a simple nil check is safe
func GetIDOrNilPlm(bo IBusinessObject) any {
	if bo == nil {
		return nil
	}
	return bo.GetID()
}

// GetModelOrNil returns nil if bo is nil, otherwise bo.GetModelName() - see [GetIDOrNil]
func GetModelOrNil(bo IBusinessObject) any {
	if bo == nil {
		return nil
	}
	return bo.GetModelName()
}

// businessObjectPtr constrains a type parameter to be a concrete pointer type (e.g. *domain.Company)
// that implements IBusinessObject; this lets [GetIDOrNilP] below compare bo to nil while bo still has
// its own, concrete pointer type - which is what makes the nil check reliable
type businessObjectPtr[T any] interface {
	*T
	IBusinessObject
}

// GetIDOrNilP is the pointer-typed equivalent of [GetIDOrNil], for relationship fields declared with
// a concrete pointer type (monomorphic relationships, e.g. *domain.Company); it must be used instead
// of GetIDOrNil there, since converting a nil *domain.Company to the IBusinessObject interface first
// would produce a non-nil interface value wrapping a nil pointer, making a plain "bo == nil" check fail
func GetIDOrNil[T any, PT businessObjectPtr[T]](bo PT) any {
	if bo == nil {
		return nil
	}
	return bo.GetID() //
}
