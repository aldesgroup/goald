package goald

import (
	"strconv"
	"strings"
)

// ------------------------------------------------------------------------------------------------
// Batch link delete - shared logic reused by every generated DAO's ExecDeleteLinksQueries
// ------------------------------------------------------------------------------------------------

// DeleteLinkContext gathers everything that's specific to one relationship's link table, for removing rows from it
type DeleteLinkContext struct {
	Table    string                                               // the link table to delete from (including its schema), e.g. "mydb.link__purchase_order__watched_by"
	Columns  string                                               // the comma-separated column list identifying a link row, in the same order as FillRows' addRow args, e.g. "source__purchase_order__id, target__watched_by__id"
	NbCols   int                                                  // the number of columns per row (2, plus 1 more per polymorphic side)
	BObjs    []IBusinessObject                                    // the (source) business objects whose links are being removed
	FillRows func(bObj IBusinessObject, addRow func(args ...any)) // called once per source object; should call addRow(...) once per target, with NbCols values
}

// ExecBatchDeleteLink performs a batched, multi-row delete of all the link rows described by the given context
func (baseDAO *BusinessObjectDAO) ExecBatchDeleteLink(ctx *DeleteLinkContext) error {

	// gathering all the rows to delete first, since we can't know their number upfront - each source
	// object may have any number of targets (including none) for this relationship
	var rows [][]any
	for _, bObj := range ctx.BObjs {
		ctx.FillRows(bObj, func(args ...any) {
			rows = append(rows, args)
		})
	}

	// nothing to do if there are no links to remove
	if len(rows) == 0 {
		return nil
	}

	// without this cast, Postgres would resolve each "__id" column of the row-value list as text, since
	// none of its rows give it any other type info, causing "bigint = text" errors against the real column
	idCol := make([]bool, ctx.NbCols)
	for i, col := range strings.Split(ctx.Columns, ",") {
		idCol[i] = strings.HasSuffix(strings.TrimSpace(col), suffixID)
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

		// actual delete query, matching each row of the batch via its full column tuple
		queryBuilder.WriteString("DELETE FROM ")
		queryBuilder.WriteString(ctx.Table)
		queryBuilder.WriteString(" WHERE (")
		queryBuilder.WriteString(ctx.Columns)
		queryBuilder.WriteString(") IN (")

		// adding the placeholders for each row, and copying the already-extracted values into the args slice
		for i, row := range batch {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}

			n := i * ctx.NbCols
			queryBuilder.WriteByte('(')
			for c := 0; c < ctx.NbCols; c++ {
				if c > 0 {
					queryBuilder.WriteString(", ")
				}
				queryBuilder.WriteByte('$')
				queryBuilder.WriteString(strconv.Itoa(n + c + 1))
				if idCol[c] {
					queryBuilder.WriteString("::bigint")
				}
			}
			queryBuilder.WriteByte(')')

			copy(batchArgs[n:n+ctx.NbCols], row)
		}

		queryBuilder.WriteByte(')')

		// link tables don't have secret columns, so there's no need for a mask here
		if _, err := baseDAO.Exec(nil, queryBuilder.String(), batchArgs...); err != nil {
			return ErrorC(err, "Error while deleting link rows from '%s'", ctx.Table)
		}
	}

	return nil
}
