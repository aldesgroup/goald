package goald

import (
	"strconv"
	"strings"
)

// ------------------------------------------------------------------------------------------------
// Batch update - shared logic reused by every generated DAO's ExecUpdateQuery
// ------------------------------------------------------------------------------------------------

// UpdateContext gathers everything that's specific to one business object model's update query
type UpdateContext struct {
	Table       string                                           // the table to update (including its schema), e.g. "mydb.staff_member"
	Columns     string                                           // the comma-separated list of columns to update, as they should appear in the SET clause, e.g. "creation, description, name"
	ColumnTypes []string                                         // the SQL type of each column in Columns, in the same order, e.g. {"timestamptz", "text"} - used to cast each VALUES entry, since Postgres would otherwise default untyped params to text
	NbCols      int                                              // the number of columns being updated per row (not counting "id")
	MaskPattern []bool                                           // the per-row secret-column mask pattern, matching Columns (not counting "id")
	BObjs       []IBusinessObject                                // the business objects to update
	FillRow     func(bObj IBusinessObject, args []any, base int) // fills in the args slice with this row's updated column values (not "id"), starting at index 'base'
}

// ExecBatchUpdate performs a batched, multi-row update of the business objects described by the given context, using
// one "UPDATE ... FROM (VALUES ...)" statement per batch so every row in a batch is updated in a single round-trip.
//
// TODO WARNING ONLY FOR POSTGRESQL: the "UPDATE ... FROM (VALUES ...) AS v(...)" syntax is Postgres-specific.
func (baseDAO *BusinessObjectDAO) ExecBatchUpdate(ctx *UpdateContext) error {
	// each row holds the object's id (used to match it to the DB row to update), plus its updated columns
	rowWidth := ctx.NbCols + 1

	// reusable buffers, sized for the largest possible batch, and reset (not reallocated) at each iteration -
	// this avoids re-allocating them on every batch
	args := make([]any, batchSize*rowWidth)

	// the masked-column pattern never changes row to row, so it's precomputed once for a full batch and simply
	// sliced down for a shorter (last) batch, instead of being rebuilt; the "id" column is never a secret
	mask := make([]bool, batchSize*rowWidth)
	for i := 0; i < batchSize; i++ {
		copy(mask[i*rowWidth+1:], ctx.MaskPattern)
	}

	// the "col = v.col, ..." SET clause and the "v(id, col, ...)" aliased columns don't change from batch to
	// batch, so they're built once upfront
	setClause, valuesColumns := baseDAO.buildUpdateFromValuesClauses(ctx.Columns)

	// avoiding a simple string and += operation, which would have the GC work a lot harder
	var queryBuilder strings.Builder

	// some DB pecularities
	var placeholder = baseDAO.model.getDB().get.QueryPlaceholder()
	var placeholderIndexed = baseDAO.model.getDB().is.QueryPlaceholderIndexed()

	for start := 0; start < len(ctx.BObjs); start += batchSize {
		end := start + batchSize
		if end > len(ctx.BObjs) {
			end = len(ctx.BObjs)
		}
		batch := ctx.BObjs[start:end]

		// slicing (not reallocating) the reusable buffers down to this batch's exact size
		batchArgs := args[:len(batch)*rowWidth]
		batchMask := mask[:len(batch)*rowWidth]
		queryBuilder.Reset()

		// actual update query, matching each row of the batch to its business object via "id"
		queryBuilder.WriteString("UPDATE ")
		queryBuilder.WriteString(ctx.Table)
		queryBuilder.WriteString(" AS t SET ")
		queryBuilder.WriteString(setClause)
		queryBuilder.WriteString(" FROM (VALUES ")

		// adding the placeholders for each row, and filling the args slice with the actual values to update
		for i, bObj := range batch {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}

			// building this row's placeholder group, e.g. "($1, $2, ...)", without fmt.Sprintf's
			// reflection-based formatting overhead
			n := i * rowWidth
			queryBuilder.WriteByte('(')
			for c := 1; c <= rowWidth; c++ {
				if c > 1 {
					queryBuilder.WriteString(", ")
				}
				queryBuilder.WriteString(placeholder)
				if placeholderIndexed {
					queryBuilder.WriteString(strconv.Itoa(n + c))
				}
				// without these casts, Postgres defaults each VALUES column to text, since none of its rows
				// give it any other type info, causing "<type> = text" or "<column> is of type <type> but
				// expression is of type text" errors
				if c == 1 {
					queryBuilder.WriteString("::bigint")
				} else if ctx.ColumnTypes != nil {
					queryBuilder.WriteString("::")
					queryBuilder.WriteString(ctx.ColumnTypes[c-2])
				}
			}
			queryBuilder.WriteByte(')')

			// the row's id comes first, so it can be matched back to the actual DB row in the WHERE clause
			batchArgs[n] = bObj.GetID()
			ctx.FillRow(bObj, batchArgs, n+1)
		}

		queryBuilder.WriteString(") AS v(")
		queryBuilder.WriteString(valuesColumns)
		queryBuilder.WriteString(") WHERE t.id = v.id")

		if _, errExec := baseDAO.Exec(batchMask, queryBuilder.String(), batchArgs...); errExec != nil {
			return baseDAO.HandleDbError(errExec)
		}
	}

	return nil
}

// buildUpdateFromValuesClauses turns a comma-separated column list, e.g. "creation, description, name", into the
// "SET" clause (e.g. "creation = v.creation, description = v.description, name = v.name") and the aliased column
// list used in the "VALUES (...) AS v(...)" part of the query (e.g. "id, creation, description, name")
func (baseDAO *BusinessObjectDAO) buildUpdateFromValuesClauses(columns string) (setClause string, valuesColumns string) {
	cols := strings.Split(columns, ",")

	var setBuilder, valuesBuilder strings.Builder
	valuesBuilder.WriteString("id")

	for i, col := range cols {
		col = strings.TrimSpace(col)

		if i > 0 {
			setBuilder.WriteString(", ")
		}
		setBuilder.WriteString(col)
		setBuilder.WriteString(" = v.")
		setBuilder.WriteString(col)

		valuesBuilder.WriteString(", ")
		valuesBuilder.WriteString(col)
	}

	return setBuilder.String(), valuesBuilder.String()
}
