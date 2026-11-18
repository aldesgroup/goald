package goald

// import (
// 	"database/sql"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// // ------------------------------------------------------------------------------------------------
// // Batch upsert - shared logic reused by every generated DAO's ExecUpsertQuery
// // ------------------------------------------------------------------------------------------------

// // ExecUpsertContext gathers everything that's specific to one business object model's upsert query
// type ExecUpsertContext struct {
// 	Table           string                                           // the table to insert into / update (including its schema), e.g. "mydb.staff_member"
// 	Columns         string                                           // the comma-separated column list, as it should appear in the INSERT statement - including the pre-ID
// 	ConflictColumns string                                           // the comma-separated columns of the UNIQUE constraint used to detect an already existing row, e.g. "external_id"
// 	NbCols          int                                              // the number of columns being inserted per row
// 	MaskPattern     []bool                                           // the per-row secret-column mask pattern, e.g. {false, false, true}
// 	BObjs           []IBusinessObject                                // the business objects to upsert
// 	FillRow         func(bObj IBusinessObject, args []any, base int) // fills in the args slice for one row, starting at index 'base'
// }

// // UpsertedRow is what the DB tells us about a row that has been inserted or updated by an upsert
// type UpsertedRow struct {
// 	ID       BObjID     // the row's DB ID, whether the row is new or not
// 	Creation *time.Time // the row's actual creation timestamp, which isn't the one of the upserted object if the row already existed
// }

// // ExecBatchUpsert performs a batched, multi-row "INSERT ... ON CONFLICT ... DO UPDATE" of the business objects described by the
// // given context; returns a map of each object's pre-ID to what the DB holds for its row, to consolidate it back onto the objects.
// //
// // All the columns but the conflict ones and the creation timestamp are overwritten on an already existing row, even if nothing
// // changed, since a row left untouched wouldn't be returned. The same unique key must not appear twice in the given objects.
// //
// // TODO WARNING ONLY FOR POSTGRESQL: the "ON CONFLICT (...) DO UPDATE ... RETURNING" syntax is Postgres-specific.
// func (baseDAO *BusinessObjectDAO) ExecBatchUpsert(ctx *ExecUpsertContext) (map[int]UpsertedRow, error) {
// 	// this will hold the mapping between each business object's pre-ID and what the DB holds for its row
// 	rowToUpsertedMap := make(map[int]UpsertedRow, len(ctx.BObjs))

// 	// reusable buffers, sized for the largest possible batch, and reset (not reallocated) at each iteration -
// 	// this avoids re-allocating them on every batch
// 	args := make([]any, batchSize*ctx.NbCols)

// 	// the masked-column pattern never changes row to row, so it's precomputed once for a full batch and
// 	// simply sliced down for a shorter (last) batch, instead of being rebuilt
// 	mask := make([]bool, batchSize*ctx.NbCols)
// 	for i := 0; i < batchSize; i++ {
// 		copy(mask[i*ctx.NbCols:], ctx.MaskPattern)
// 	}

// 	// the "col = EXCLUDED.col, ..." clause doesn't change from batch to batch, so it's built once upfront
// 	updateClause := buildUpsertUpdateClause(ctx.Columns, ctx.ConflictColumns)

// 	// avoiding a simple string and += operation, which would have the GC work a lot harder
// 	var queryBuilder strings.Builder

// 	// reused scan targets: only their address is needed by Scan, no need for a fresh variable on every row
// 	var preID int
// 	var newID BObjID
// 	var creation sql.NullTime

// 	// some DB pecularities
// 	var placeholder = baseDAO.model.getDB().get.QueryPlaceholder()
// 	var placeholderIndexed = baseDAO.model.getDB().is.QueryPlaceholderIndexed()

// 	for start := 0; start < len(ctx.BObjs); start += batchSize {
// 		end := min(start+batchSize, len(ctx.BObjs))
// 		batch := ctx.BObjs[start:end]

// 		// slicing (not reallocating) the reusable buffers down to this batch's exact size
// 		batchArgs := args[:len(batch)*ctx.NbCols]
// 		batchMask := mask[:len(batch)*ctx.NbCols]
// 		queryBuilder.Reset()

// 		// actual insert query
// 		queryBuilder.WriteString("INSERT INTO ")
// 		queryBuilder.WriteString(ctx.Table)
// 		queryBuilder.WriteString(" (")
// 		queryBuilder.WriteString(ctx.Columns)
// 		queryBuilder.WriteString(") VALUES ")

// 		// adding the placeholders for each row, and filling the args slice with the actual values to insert
// 		for i, bObj := range batch {
// 			if i > 0 {
// 				queryBuilder.WriteString(", ")
// 			}

// 			// building this row's placeholder group, e.g. "($1, $2, ...)", without fmt.Sprintf's
// 			// reflection-based formatting overhead
// 			n := i * ctx.NbCols
// 			queryBuilder.WriteByte('(')
// 			for c := 1; c <= ctx.NbCols; c++ {
// 				if c > 1 {
// 					queryBuilder.WriteString(", ")
// 				}
// 				queryBuilder.WriteString(placeholder)
// 				if placeholderIndexed {
// 					queryBuilder.WriteString(strconv.Itoa(n + c))
// 				}
// 			}
// 			queryBuilder.WriteByte(')')

// 			// writing straight to the known index, rather than appending, since the exact size is known upfront
// 			ctx.FillRow(bObj, batchArgs, n)
// 		}

// 		// what to do when the unique key is already there
// 		queryBuilder.WriteString(" ON CONFLICT (")
// 		queryBuilder.WriteString(ctx.ConflictColumns)
// 		queryBuilder.WriteString(") DO UPDATE SET ")
// 		queryBuilder.WriteString(updateClause)

// 		// TODO WARNING ONLY FOR POSTGRESQL: we need to retrieve the IDs, along with the pre-IDs they relate to, so we can consolidate them back onto the business objects
// 		queryBuilder.WriteString(" RETURNING _pre_id, id, creation")

// 		// executing the query and retrieving the IDs, along with the pre-IDs they relate to;
// 		// wrapped in a func so 'defer rows.Close()' runs at the end of *this batch*, not at the end of
// 		// ExecBatchUpsert - otherwise every batch's rows would stay open until the whole function returns
// 		if errBatch := func() (err error) {
// 			rows, errQuery := baseDAO.Query(batchMask, queryBuilder.String(), batchArgs...)
// 			if errQuery != nil {
// 				return baseDAO.HandleDbError(errQuery)
// 			}

// 			defer func() {
// 				if errClose := rows.Close(); errClose != nil && err == nil {
// 					err = ErrorC(errClose, "Could not close the rows")
// 				}
// 			}()

// 			for rows.Next() {
// 				if errScan := rows.Scan(&preID, &newID, &creation); errScan != nil {
// 					return ErrorC(errScan, "Could not scan the row")
// 				}

// 				upserted := UpsertedRow{ID: newID}
// 				if creation.Valid {
// 					creationTime := creation.Time // a copy, since the scan target is reused
// 					upserted.Creation = &creationTime
// 				}

// 				rowToUpsertedMap[preID] = upserted
// 			}

// 			if errRows := rows.Err(); errRows != nil && err == nil {
// 				err = ErrorC(errRows, "Error iterating over the rows")
// 			}

// 			return
// 		}(); errBatch != nil {
// 			return nil, errBatch
// 		}
// 	}

// 	return rowToUpsertedMap, nil
// }

// // buildUpsertUpdateClause turns the insert column list, e.g. "_pre_id, creation, external_id, name", and the conflict
// // column list, e.g. "external_id", into the "DO UPDATE SET" clause, e.g. "_pre_id = EXCLUDED._pre_id, name = EXCLUDED.name".
// // The pre-ID has to be overwritten too, so that what comes back from the DB matches the objects of the current call.
// func buildUpsertUpdateClause(columns string, conflictColumns string) string {
// 	// these columns must keep the values they already have in the DB
// 	kept := map[string]bool{"creation": true}
// 	for _, conflictColumn := range strings.Split(conflictColumns, ",") {
// 		kept[strings.TrimSpace(conflictColumn)] = true
// 	}

// 	var clause strings.Builder
// 	for _, col := range strings.Split(columns, ",") {
// 		col = strings.TrimSpace(col)
// 		if kept[col] {
// 			continue
// 		}

// 		if clause.Len() > 0 {
// 			clause.WriteString(", ")
// 		}
// 		clause.WriteString(col)
// 		clause.WriteString(" = EXCLUDED.")
// 		clause.WriteString(col)
// 	}

// 	return clause.String()
// }
