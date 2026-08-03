package goald

import (
	"database/sql"
	"strconv"
	"strings"
)

// ------------------------------------------------------------------------------------------------
// Batch insert - shared logic reused by every generated DAO's ExecInsertQuery
// ------------------------------------------------------------------------------------------------

// BatchInsertContext gathers everything that's specific to one business object model's insert query,
// so that all the batching/query-building/scanning logic can be shared by every DAO's ExecInsertQuery,
// instead of being duplicated in each generated file.
type BatchInsertContext struct {
	Table       string                                           // the table to insert into (including its schema), e.g. "mydb.staff_member"
	Columns     string                                           // the comma-separated column list, as it should appear in the INSERT statement
	NbCols      int                                              // the number of columns being inserted per row
	MaskPattern []bool                                           // the per-row secret-column mask pattern, e.g. {false, false, true}
	BObjs       []IBusinessObject                                // the business objects to insert
	FillRow     func(bObj IBusinessObject, args []any, base int) // fills in the args slice for one row, starting at index 'base'
}

// ExecBatchInsert performs a batched, multi-row insert of the business objects described by the given
// context, in batches of 1000, and returns a map of each object's pre-ID to its newly assigned DB ID -
// this mapping is then used by dbInsert to consolidate the real DB IDs back onto the business objects.
func (baseDAO *BusinessObjectDAO) ExecBatchInsert(ctx *BatchInsertContext) (map[int]int64, error) {
	// were going to do multiple inserts in batches of 1000, instead of doing one insert per statement
	const batchSize = 1000

	// this will hold the mapping between each business object's pre-ID and its newly assigned DB ID,
	// allowing for the consolidation done later on in dbInsert
	rowToIDMap := make(map[int]int64, len(ctx.BObjs))

	// reusable buffers, sized for the largest possible batch, and reset (not reallocated) at each iteration -
	// this avoids re-allocating them on every batch
	args := make([]any, batchSize*ctx.NbCols)

	// the masked-column pattern never changes row to row, so it's precomputed once for a full batch and
	// simply sliced down for a shorter (last) batch, instead of being rebuilt
	mask := make([]bool, batchSize*ctx.NbCols)
	for i := 0; i < batchSize; i++ {
		copy(mask[i*ctx.NbCols:], ctx.MaskPattern)
	}

	// avoiding a simple string and += operation, which would have the GC work a lot harder
	var queryBuilder strings.Builder

	// reused scan targets: only their address is needed by Scan, no need for a fresh variable on every row
	var preID int
	var newID int64

	for start := 0; start < len(ctx.BObjs); start += batchSize {
		end := start + batchSize
		if end > len(ctx.BObjs) {
			end = len(ctx.BObjs)
		}
		batch := ctx.BObjs[start:end]

		// slicing (not reallocating) the reusable buffers down to this batch's exact size
		batchArgs := args[:len(batch)*ctx.NbCols]
		batchMask := mask[:len(batch)*ctx.NbCols]
		queryBuilder.Reset()

		// actual insert query
		queryBuilder.WriteString("INSERT INTO ")
		queryBuilder.WriteString(ctx.Table)
		queryBuilder.WriteString(" (")
		queryBuilder.WriteString(ctx.Columns)
		queryBuilder.WriteString(") VALUES ")

		// adding the placeholders for each row, and filling the args slice with the actual values to insert
		for i, bObj := range batch {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}

			// building this row's placeholder group, e.g. "($1, $2, ...)", without fmt.Sprintf's
			// reflection-based formatting overhead
			n := i * ctx.NbCols
			queryBuilder.WriteByte('(')
			for c := 1; c <= ctx.NbCols; c++ {
				if c > 1 {
					queryBuilder.WriteString(", ")
				}
				queryBuilder.WriteByte('$')
				queryBuilder.WriteString(strconv.Itoa(n + c))
			}
			queryBuilder.WriteByte(')')

			// writing straight to the known index, rather than appending, since the exact size is known upfront
			ctx.FillRow(bObj, batchArgs, n)
		}

		queryBuilder.WriteString(" RETURNING _pre_id, id")

		// executing the insert query and retrieving the new IDs, along with the pre-IDs they relate to;
		// wrapped in a func so 'defer rows.Close()' runs at the end of *this batch*, not at the end of
		// ExecBatchInsert - otherwise every batch's rows would stay open until the whole function returns
		if errBatch := func() (err error) {
			rows, errQuery := baseDAO.Query(batchMask, queryBuilder.String(), batchArgs...)
			if errQuery != nil {
				return baseDAO.HandleDbError(errQuery)
			}

			defer func() {
				if errClose := rows.Close(); errClose != nil && err == nil {
					err = ErrorC(errClose, "Could not close the rows")
				}
			}()

			for rows.Next() {
				if errScan := rows.Scan(&preID, &newID); errScan != nil {
					return ErrorC(errScan, "Could not scan the row")
				}

				rowToIDMap[preID] = newID
			}

			if errRows := rows.Err(); errRows != nil && err == nil {
				err = ErrorC(errRows, "Error iterating over the rows")
			}

			return
		}(); errBatch != nil {
			return nil, errBatch
		}
	}

	return rowToIDMap, nil
}

// ------------------------------------------------------------------------------------------------
// Batch link insert - shared logic reused by every generated DAO's ExecInsertLinksQueries
// ------------------------------------------------------------------------------------------------

// LinkInsertContext gathers everything that's specific to one relationship's link table, so that all the
// batching/query-building logic can be shared by every DAO's ExecInsertLinksQueries, instead of being
// duplicated once per relationship, per model.
//
// One LinkInsertContext should be built, and ExecBatchLinkInsert called, for each relationship of the model
// that's persisted through a link table (see IBusinessObjectModel.getRelationshipsWithLinkTable) - i.e. for
// relationships owning the link, on the "source" side. E.g. UserGroup.Members is the back-reference of
// User.MemberOf, so it's User.MemberOf (the source-to-target side) that owns the link table and generates
// a LinkInsertContext, not UserGroup.Members.
type LinkInsertContext struct {
	Table    string                                               // the link table to insert into (including its schema), e.g. "mydb.link__purchase_order__watched_by"
	Columns  string                                               // the comma-separated column list, as it should appear in the INSERT statement
	NbCols   int                                                  // the number of columns per row (2, plus 1 more per polymorphic side)
	BObjs    []IBusinessObject                                    // the (source) business objects being inserted
	FillRows func(bObj IBusinessObject, addRow func(args ...any)) // called once per source object; should call addRow(...) once per target, with NbCols values
}

// ExecBatchLinkInsert performs a batched, multi-row insert of all the link rows described by the given
// context, in batches of 1000 rows, for a single relationship's link table.
func (baseDAO *BusinessObjectDAO) ExecBatchLinkInsert(ctx *LinkInsertContext) error {
	// were going to do multiple inserts in batches of 1000, instead of doing one insert per statement
	const batchSize = 1000

	// gathering all the rows to insert first, since we can't know their number upfront - each source
	// object may have any number of targets (including none) for this relationship
	var rows [][]any
	for _, bObj := range ctx.BObjs {
		ctx.FillRows(bObj, func(args ...any) {
			rows = append(rows, args)
		})
	}

	// nothing to do if there are no links to persist
	if len(rows) == 0 {
		return nil
	}

	// reusable buffer, sized for the largest possible batch, and reset (not reallocated) at each iteration
	args := make([]any, batchSize*ctx.NbCols)

	// avoiding a simple string and += operation, which would have the GC work a lot harder
	var queryBuilder strings.Builder

	for start := 0; start < len(rows); start += batchSize {
		end := min(start+batchSize, len(rows))
		batch := rows[start:end]

		// slicing (not reallocating) the reusable buffer down to this batch's exact size
		batchArgs := args[:len(batch)*ctx.NbCols]
		queryBuilder.Reset()

		// actual insert query
		queryBuilder.WriteString("INSERT INTO ")
		queryBuilder.WriteString(ctx.Table)
		queryBuilder.WriteString(" (")
		queryBuilder.WriteString(ctx.Columns)
		queryBuilder.WriteString(") VALUES ")

		// adding the placeholders for each row, and copying the already-extracted values into the args slice
		for i, row := range batch {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}

			n := i * ctx.NbCols
			queryBuilder.WriteByte('(')
			for c := 1; c <= ctx.NbCols; c++ {
				if c > 1 {
					queryBuilder.WriteString(", ")
				}
				queryBuilder.WriteByte('$')
				queryBuilder.WriteString(strconv.Itoa(n + c))
			}
			queryBuilder.WriteByte(')')

			copy(batchArgs[n:n+ctx.NbCols], row)
		}

		// link tables don't have secret columns, so there's no need for a mask here
		if _, err := baseDAO.Exec(nil, queryBuilder.String(), batchArgs...); err != nil {
			return ErrorC(err, "Error while inserting link rows into '%s'", ctx.Table)
		}
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Select - shared logic reused by every generated DAO's ExecSelectQuery
// ------------------------------------------------------------------------------------------------

// SelectContext gathers everything that's specific to one business object model's select query, so
// that all the query-building/scanning logic can be shared by every DAO's ExecSelectQuery, instead of
// being duplicated in each generated file. There are 2 model-specific "placeholders" to provide:
//   - BuildWhere, which turns a query's name and param values into a *QueryArgs holding the OR-ed,
//     AND-ed condition clauses mirroring the query built with Find(...).Where(...) in a *--blo.go file
//     (see BuildWhereFromClauses for how these clauses turn into actual WHERE clause text), plus the
//     args needed to fill in their placeholders and a matching per-arg secret mask; as more than one
//     query gets registered for a given model, BuildWhere should switch on its queryName argument to
//     build the right QueryArgs for each. Note the mask here is about the *query params* object's
//     fields (e.g. PurchaseOrderQuery), NOT the selected business object's own fields/columns - these
//     args are what ends up filling in the WHERE clause's placeholders, not the SELECT columns
//   - ScanRow, which instantiates and fills in one business object from the current row
type SelectContext struct {
	Table     string                                                          // the table to select from (including its schema), e.g. "mydb.purchase_order"
	Columns   string                                                          // the comma-separated column list, as it should appear in the SELECT statement
	QueryArgs func(queryName QueryName, values IQueryParamsObject) *QueryArgs // builds the WHERE clause's condition clauses, its args, and their secret mask, for the given query & param values
	ScanRow   func(rows *sql.Rows) (IBusinessObject, error)                   // instantiates & fills in one business object, from the current row
}

// ExecSelect performs the select query described by the given context - applying the query's WHERE
// clause, if any - and returns the resulting business objects.
func (baseDAO *BusinessObjectDAO) ExecSelect(ctx *SelectContext, queryName QueryName, values IQueryParamsObject) (result []IBusinessObject, err error) {
	// let's build the query's WHERE clause parts, from the query name and the given param values
	queryArgs := ctx.QueryArgs(queryName, values)

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
		bObj, errScan := ctx.ScanRow(rows)
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
	OrClauses  [][]string // each entry is a list of AND-ed condition strings; entries are OR-ed together
	Args       []any
	Mask       []bool
	currentAnd []string // the AND-clause currently being accumulated, via AppendAndClause/AppendRawAndClause
}

// NewQueryArgs creates a QueryArgs
func NewQueryArgs(capacity int) *QueryArgs {
	return &QueryArgs{
		Args: make([]any, 0, capacity),
		Mask: make([]bool, 0, capacity),
	}
}

// addSingleClause registers 1 arg value (and whether it's secret) and returns the resulting condition
func (qa *QueryArgs) AddSingleClause(criterion string, value any, secret bool) string {
	qa.Args = append(qa.Args, value)
	qa.Mask = append(qa.Mask, secret)
	return criterion + "$" + strconv.Itoa(len(qa.Args))
}

// AppendAndClause registers 1 arg value (and whether it's secret) and appends the resulting condition
// (prefix + its placeholder) to the AND-clause currently being built - see NewAndClause to commit it
// and move on to the next, OR-ed one.
func (qa *QueryArgs) AppendAndClause(condition bool, prefix string, value any, secret bool) {
	if condition {
		qa.currentAnd = append(qa.currentAnd, qa.AddSingleClause(prefix, value, secret))
	}
}

// AppendRawAndClause appends an already-built condition string to the AND-clause currently being
// built - useful for conditions spanning more than 1 placeholder (e.g. an IN (...) clause), which
// AppendAndClause can't express directly since it only ever registers a single value/placeholder:
// qa.AppendRawAndClause(qa.Add("status IN (", v1, false) + ", " + qa.Add("", v2, false) + ")")
func (qa *QueryArgs) AppendRawAndClause(condition bool, rawClause string) {
	if condition {
		qa.currentAnd = append(qa.currentAnd, rawClause)
	}
}

// NewAndClause commits the AND-clause currently being built (possibly empty/nil, meaning it has no
// active condition and thus takes no part at all in the final WHERE clause) as one more OR-ed clause,
// and starts a fresh, empty one to accumulate the next clause's conditions into.
func (qa *QueryArgs) NewAndClause() {
	qa.OrClauses = append(qa.OrClauses, qa.currentAnd)
	qa.currentAnd = nil
}

// toWhereClause assembles a "WHERE (...) OR (...) ..." clause from this QueryArgs' OR-ed
// clauses, skipping any clause left empty (i.e. with no active condition) - this is the generic,
// model-agnostic part of turning a Find(...).Where(...) query (built with Either/Or in a *--blo.go
// file) into actual SQL: by OR's associativity, any such query flattens into exactly this shape - a
// flat list of AND-ed clauses, OR-ed together, however many there are.
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
