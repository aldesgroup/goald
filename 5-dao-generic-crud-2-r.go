package goald

import "database/sql"

// ------------------------------------------------------------------------------------------------
// Load - shared logic reused by every generated DAO's ExecReadQuery
// ------------------------------------------------------------------------------------------------

// ReadContext gathers everything that's specific to one business object model's read query
type ReadContext struct {
	Table     string                                                                          // the table to read from (including its schema), e.g. "mydb.purchase_order"
	Columns   string                                                                          // the comma-separated column list, as it should appear in the SELECT statement
	Objects   map[BObjID]IBusinessObject                                                      // the map of business objects to fill with the read data
	ObjectIDs []any                                                                           // the list of business object IDs to read from the database
	ScanRow   func(rows *sql.Rows, bObjs map[BObjID]IBusinessObject) (IBusinessObject, error) // instantiates & fills in one business object, from the current row
}

// ExecRead performs the read query described by the given context and fills the provided business objects.
func (baseDAO *BusinessObjectDAO) ExecRead(ctx *ReadContext) (err error) {
	// building the WHERE clause with the list of business object IDs to read - which depends on the DB type,
	// since the placeholder syntax is different between Postgres and MySQL
	where := " WHERE id IN " + baseDAO.makePlaceholdersString(len(ctx.ObjectIDs))

	// executing the query and retrieving the rows
	rows, errQuery := baseDAO.Query(nil, "SELECT "+ctx.Columns+" FROM "+ctx.Table+where, ctx.ObjectIDs...)
	if errQuery != nil {
		return baseDAO.HandleDbError(errQuery)
	}

	// making sure we don't leak the rows/underlying connection
	defer func() {
		if errClose := rows.Close(); errClose != nil && err == nil {
			err = ErrorC(errClose, "Could not close the rows")
		}
	}()

	// scanning each row into a business object, and returning the list of them
	for rows.Next() {
		if _, errScan := ctx.ScanRow(rows, ctx.Objects); errScan != nil {
			return ErrorC(errScan, "Could not scan the row into a business object")
		}
	}

	// making sure we don't have any error from the rows iteration itself
	if errRows := rows.Err(); errRows != nil {
		return ErrorC(errRows, "Error iterating over the rows")
	}

	return nil
}
