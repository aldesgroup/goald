package goald

import "database/sql"

// ------------------------------------------------------------------------------------------------
// Read relationship - shared logic reused by every generated DAO's ExecReadRelationshipQuery
// ------------------------------------------------------------------------------------------------

// ReadRelationshipContext gathers everything that's specific to reading one relationship's rowsZ
type ReadRelationshipContext struct {
	Query     string                                        // the query to read the relationship's rows (but only the object's ID and sometimes model)
	Args      []any                                         // the query's positional args
	AttachRow func(rows *sql.Rows) (IBusinessObject, error) // scans one row into a target BO (ID, maybe model), and attaches it to the relevant business object it relates to
}

// ExecReadRelationship runs the query described by the given context, and returns the resulting
// "empty" target BOs, which have been attached to their respective parent business objects by the ScanRow function.
func (baseDAO *BusinessObjectDAO) ExecReadRelationship(ctx *ReadRelationshipContext) (result []IBusinessObject, err error) {
	rows, errQuery := baseDAO.Query(nil, ctx.Query, ctx.Args...)
	if errQuery != nil {
		return nil, baseDAO.HandleDbError(errQuery)
	}

	// making sure we don't leak the rows/underlying connection
	defer func() {
		if errClose := rows.Close(); errClose != nil && err == nil {
			err = ErrorC(errClose, "Could not close the rows")
		}
	}()

	// instantiating & attaching the related BOs
	for rows.Next() {
		target, errScan := ctx.AttachRow(rows)
		if errScan != nil {
			return nil, ErrorC(errScan, "Could not scan the relationship row")
		}

		result = append(result, target)
	}

	// handling errors
	if errRows := rows.Err(); errRows != nil {
		return nil, ErrorC(errRows, "Error iterating over the rows")
	}

	return result, nil
}
