package goald

import (
	"database/sql"
	"strconv"
	"strings"
)

// ------------------------------------------------------------------------------------------------
// Search - shared logic reused by every generated DAO's ExecSearchQuery
// ------------------------------------------------------------------------------------------------

// SearchContext gathers everything that's specific to one business object model's search query
type SearchContext struct {
	Table         string                                                          // the table to search into (including its schema), e.g. "mydb.purchase_order"
	Columns       string                                                          // the comma-separated column list, as it should appear in the SELECT statement
	QueryName     QueryName                                                       // the name of the search query, e.g. "byClientRef"
	QueryValues   ISearchParamValues                                              // the values of the search query's parameters, e.g. a struct with a ClientRef field
	MakeQueryArgs func(queryName QueryName, values ISearchParamValues) *QueryArgs // builds the WHERE clause's condition clauses, its args, and their secret mask, for the given query & param values
	BoCache       *BObjCache                                                      // the cache for all the objects being read along the way
	ScanRow       func(rows *sql.Rows, cache *BObjCache) (IBusinessObject, error) // instantiates & fills in one business object, from the current row
}

// ExecSearch performs the search query described by the given context and returns the resulting business objects.
func (baseDAO *BusinessObjectDAO) ExecSearch(ctx *SearchContext) (result []IBusinessObject, err error) {
	// let's build the query's WHERE clause parts, from the query name and the given param values
	queryArgs := ctx.MakeQueryArgs(ctx.QueryName, ctx.QueryValues)

	// turning the OR-ed, AND-ed condition clauses into a single WHERE clause string, with its placeholders
	where := queryArgs.toWhereClauseAsString()

	// executing the query and retrieving the rows
	rows, errQuery := baseDAO.Query(queryArgs.Mask, "SELECT "+ctx.Columns+" FROM "+ctx.Table+where, queryArgs.Args...)
	if errQuery != nil {
		return nil, baseDAO.HandleDbError(errQuery)
	}

	// making sure we don't leak the rows/underlying connection
	defer func() {
		if errClose := rows.Close(); errClose != nil && err == nil {
			err = ErrorC(errClose, "Could not close the rows")
		}
	}()

	// scanning each row into a business object, and returning the list of them
	for rows.Next() {
		bObj, errScan := ctx.ScanRow(rows, ctx.BoCache)
		if errScan != nil {
			return nil, ErrorC(errScan, "Could not scan the row into a business object")
		}

		result = append(result, bObj)
	}

	// making sure we don't have any error from the rows iteration itself
	if errRows := rows.Err(); errRows != nil {
		return nil, ErrorC(errRows, "Error iterating over the rows")
	}

	return result, nil
}

// ------------------------------------------------------------------------------------------------
// WHERE-clause building helpers, reused by every generated DAO's BuildWhere function
// ------------------------------------------------------------------------------------------------

// QueryArgs accumulates everything needed to build a query's WHERE clause:
// - the OR-ed clauses,
// - each itself a list of AND-ed condition strings,
// - the positional args filling in their placeholders,
// - and each arg's secret mask
type QueryArgs struct {
	OrClauses          [][]string // each entry is a list of AND-ed condition strings; entries are OR-ed together
	Args               []any
	Mask               []bool
	currentAnd         []string // the AND-clause currently being accumulated, via AppendAndClause/AppendRawAndClause
	placeholder        string   // the placeholder syntax for this DB type, e.g. "$" for Postgres, "?" for MySQL
	placeholderIndexed bool     // tells whether the query placeholders are indexed, e.g. "$1", "$2", etc. for Postgres, or just "?" for MySQL
}

// NewQueryArgs creates a QueryArgs
func NewQueryArgs(capacity int, placeholder string, placeholderIndexed bool) *QueryArgs {
	return &QueryArgs{
		Args:               make([]any, 0, capacity),
		Mask:               make([]bool, 0, capacity),
		placeholder:        placeholder,
		placeholderIndexed: placeholderIndexed,
	}
}

// NewQueryArgsWithINClause creates a QueryArgs and a corresponding "IN" clause for the given column and values.
func NewQueryArgsWithINClause(placeholder string, placeholderIndexed bool, values []any) (string, *QueryArgs) {
	queryArgs := NewQueryArgs(len(values), placeholder, placeholderIndexed)

	var clause strings.Builder
	clause.WriteString(" IN (")
	for i, value := range values {
		if i > 0 {
			clause.WriteString(", ")
		}
		clause.WriteString(queryArgs.AddSingleClause("", value, false))
	}
	clause.WriteString(")")

	return clause.String(), queryArgs
}

// addSingleClause registers 1 arg value (and whether it's secret) and returns the resulting condition
func (qa *QueryArgs) AddSingleClause(criterion string, value any, secret bool) string {
	qa.Args = append(qa.Args, value)
	qa.Mask = append(qa.Mask, secret)
	if qa.placeholderIndexed {
		return criterion + qa.placeholder + strconv.Itoa(len(qa.Args))
	} else {
		return criterion + qa.placeholder
	}
}

// AppendAndClause registers 1 arg value (and whether it's secret) and appends the resulting condition
// (criterion + its placeholder) to the AND-clause currently being built
func (qa *QueryArgs) AppendAndClause(condition bool, criterion string, value any, secret bool) {
	if condition {
		qa.currentAnd = append(qa.currentAnd, qa.AddSingleClause(criterion, value, secret))
	}
}

// AppendRawAndClause appends an already-built condition string to the AND-clause currently being built
func (qa *QueryArgs) AppendRawAndClause(condition bool, rawClause string) {
	if condition {
		qa.currentAnd = append(qa.currentAnd, rawClause)
	}
}

// NewAndClause commits the AND-clause currently being built
func (qa *QueryArgs) NewAndClause() {
	qa.OrClauses = append(qa.OrClauses, qa.currentAnd)
	qa.currentAnd = nil
}

// toWhereClause assembles a "WHERE (...) OR (...) ..." clause from this QueryArgs' OR-ed clauses
func (qa *QueryArgs) toWhereClauseAsString() string {
	// committing whatever AND-clause is still pending (the last one built) before assembling
	qa.NewAndClause()

	var where strings.Builder

	for _, andClauses := range qa.OrClauses {
		if len(andClauses) == 0 {
			continue
		}

		if where.Len() == 0 {
			where.WriteString(" WHERE (")
		} else {
			where.WriteString(") OR (")
		}

		for i, cond := range andClauses {
			if i > 0 {
				where.WriteString(" AND ")
			}
			where.WriteString(cond)
		}
	}

	if where.Len() > 0 {
		where.WriteString(")")
	}

	return where.String()
}
