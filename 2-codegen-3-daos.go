// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the DAO files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"sort"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/utils"
)

const daoINITxTEMPLATE = `// Generated file, do not edit!
package %[1]s

import (
	"database/sql"
	"fmt"
	"sync"

	"%[2]s"
	"github.com/aldesgroup/goald"%[7]s
)

// ------------------------------------------------------------------------------------------------
// %[3]sDAO declaration and registration
// ------------------------------------------------------------------------------------------------

type %[3]sDAO struct {
	goald.BusinessObjectDAO
}

// type check
var _ goald.IBusinessObjectDAO = (*%[3]sDAO)(nil)

// registration
func init() {
	goald.RegisterDAO("%[3]s", &%[3]sDAO{})
}

// NewDAO implements [goald.IBusinessObjectDAO].
func (thisDAO *%[3]sDAO) NewDAO() goald.IBusinessObjectDAO {
	%[4]sDBOnce.Do(func() {
		%[4]sDB = string(goald.GetDB("%[5]s").Name())
	})

	return &%[3]sDAO{}
}

// ------------------------------------------------------------------------------------------------
// Local state
// ------------------------------------------------------------------------------------------------

var (
	%[4]sDB     string
	%[4]sDBOnce sync.Once
	%[4]sMask   = []bool{%[6]s}
)

// ------------------------------------------------------------------------------------------------
// Insertion
// ------------------------------------------------------------------------------------------------

`

// used to insert a literal backtick inside a raw string template, since raw strings can't escape backticks
const bq = "`"

const daoFILExSUFFIX = "--dao.go"
const daoFILExSUFFIXxLEN = len(daoFILExSUFFIX)

// ------------------------------------------------------------------------------------------------
// Main DAO files generation methods
// ------------------------------------------------------------------------------------------------

func (thisServer *server) generateAllObjectDAOs(srcdir string, regen bool) (codeChanged bool) {
	// a type just used here
	type daoFile struct {
		modTime  time.Time
		filename string
	}

	// first, getting all the DB types we have to deal with
	dbTypes := map[dbconn.DatabaseType]bool{}
	for _, db := range dbRegistry.databases {
		if db.schema != nil {
			dbTypes[db.schema.DbServer.Type] = true
		}
	}

	// we're going to keep track of all the DAO folders to import in all the registry files
	daoFolders := map[string]map[dbconn.DatabaseType]bool{}

	// we'll gather all the existing DAO files, per DB type
	existingDAOFiles := map[dbconn.DatabaseType]map[utils.ModelName]*daoFile{}
	for dbType := range dbTypes {
		// some init
		existingDAOFiles[dbType] = map[utils.ModelName]*daoFile{}

		// making sure the DAO folder exists for this DB type
		daoDir := core.EnsureDir(srcdir, includePATH, dbFOLDERNAME, string(dbType))

		// so, let's read this folder now
		for _, daoEntry := range core.EnsureReadDir(daoDir) {
			daoEntryInfo, errInfo := daoEntry.Info()
			core.PanicMsgIfErr(errInfo, "Could not read info for file '%s'", daoEntry.Name())
			modelName := utils.ModelName(core.KebabToPascal(daoEntry.Name()[:len(daoEntry.Name())-daoFILExSUFFIXxLEN]))
			existingDAOFiles[dbType][modelName] = &daoFile{
				modTime:  daoEntryInfo.ModTime(),
				filename: daoEntry.Name(),
			}
		}
	}

	// let's now generate all the DAOs we need in these DB folders
	for name, boModel := range modelRegistry.items {
		// we only consider the concrete business objects that are persisted in this project
		if !boModel.isInterface() && !boModel.isAbstract() && boModel.isPersistedHere() {
			// what's the DB type involved here?
			dbType := boModel.getDB().schema.DbServer.Type

			// do we meed to generate a DAO for this business object ?
			if existingDAO := existingDAOFiles[dbType][name]; regen ||
				existingDAO == nil || existingDAO.modTime.Before(boModel.getLastBOMod()) {

				// generating the missing or outdated DAO
				thisServer.generateOneDAO(srcdir, string(dbType), boModel)

				// code has been changed
				codeChanged = true

				// keeping track of the package needed to import the DAO package in its registry file
				if _, ok := daoFolders[boModel.getSrcPath()]; !ok {
					daoFolders[boModel.getSrcPath()] = map[dbconn.DatabaseType]bool{}
				}
				daoFolders[boModel.getSrcPath()][dbType] = true
			}

			// flagging this business object DAO as NOT unneeded (i.e. needed)
			delete(existingDAOFiles[dbType], name)
		}

	}

	// // iterating over each package for which we've already got a registry
	// for _, dbtype := range core.EnsureReadDir(srcdir, includePATH, dbFOLDERNAME) {
	// 	codeChangedHere := false

	// 	// if it's not a directory, we skip it
	// 	if !dbtype.IsDir() {
	// 		continue
	// 	}

	// 	// where the DAO files will be generated
	// 	daoDir := core.EnsureDir(srcdir, includePATH, includeDirEntry.Name(), daoFOLDERxNAME)

	// 	// removing the unneeded DAOs
	// 	for _, unneededDAO := range existingDAOFiles {
	// 		thisServer.Info(fmt.Sprintf("removing %s", unneededDAO.filename))
	// 		if errRem := os.Remove(path.Join(daoDir, unneededDAO.filename)); errRem != nil {
	// 			core.PanicMsgIfErr(errRem, "Could not delete DAO file '%s'", unneededDAO.filename)
	// 		}
	// 	}

	return
}

func (thisServer *server) generateOneDAO(srcdir string, dbType string, model IBusinessObjectModel) bool {
	// where the DAO should end up
	daoDir := core.EnsureDir(srcdir, includePATH, dbFOLDERNAME, dbType)

	// trivial filling of the template
	modelNameCamel := core.PascalToCamel(string(model.GetName()))
	importForBOModel := getImportPackageLine(model)
	_, _, maskPattern, _ := thisServer.getColumnsAndInsertData(model, "\t\t\t")

	// gathering the extra imports needed by the select-query-building/scanning code below (e.g. a
	// relationship's target package, or the query params' package if it's not this model's own) -
	// this model's own package is always already imported, so it's never added twice
	extraImports := map[string]bool{}
	selectQueryCode := thisServer.generateExecSearchQuery(model, extraImports)
	delete(extraImports, importForBOModel)

	content := fmt.Sprintf(daoINITxTEMPLATE, dbType, importForBOModel, model.GetName(), modelNameCamel,
		model.getDB().name, maskPattern, buildExtraImportsBlock(extraImports))

	// adding all the needed DAO methods
	content += thisServer.generateExecCreateQuery(model)
	content += "\n\n"
	content += thisServer.generateExecCreateLinksQueries(model)
	content += "\n\n"
	content += thisServer.generateExecReadQuery(model)
	content += "\n\n"
	content += selectQueryCode

	// writing to file
	core.WriteToFile(content, daoDir, core.PascalToKebab(string(model.GetName()))+daoFILExSUFFIX)

	thisServer.Info(fmt.Sprintf("(Re-)generated DAO for %s", model.GetName()))

	return true
}

// buildExtraImportsBlock turns a set of extra import paths (gathered while generating a DAO's
// select-query-building/scanning code) into the source snippet to splice right after the standard
// imports in daoINITxTEMPLATE, e.g. "\n\t\"git-ext.aldes.com/.../domain\"".
func buildExtraImportsBlock(extraImports map[string]bool) string {
	if len(extraImports) == 0 {
		return ""
	}

	paths := make([]string, 0, len(extraImports))
	for importPath := range extraImports {
		paths = append(paths, importPath)
	}
	sort.Strings(paths)

	var block strings.Builder
	for _, importPath := range paths {
		block.WriteString(fmt.Sprintf("\n\t%q", importPath))
	}

	return block.String()
}

// ------------------------------------------------------------------------------------------------
//  DAO utils
// ------------------------------------------------------------------------------------------------

// getColumnsAndInsertData gathers everything needed to generate a batched, multi-row insert for the given
// model: the column list, the per-row 'args[base+N] = ...' assignment statements (for use inside a
// ExecCreateContext.FillRow closure), the masked-column pattern (for hiding secret values from the logs),
// and the total number of columns involved.
func (thisServer *server) getColumnsAndInsertData(model IBusinessObjectModel, space string) (columns, argAssignments, maskPattern string, nbCols int) {
	varName := core.PascalToCamel(string(model.GetName()))

	for _, prop := range model.getPersistedProperties() {
		// the ID column is auto-generated by the DB (and retrieved via the RETURNING clause), so it
		// must never be part of the columns/values being inserted
		if prop.GetName() == BoFieldID {
			continue
		}

		if columns != "" {
			columns += ", "
			maskPattern += ", "
		}

		colIndex := nbCols
		nbCols++

		if relationship, ok := prop.(*Relationship); ok {
			columns += relationship.getColumnName()
			maskPattern += core.IfThenElse(relationship.isSecret(), "true", "false")

			// the relationship's Go field is either a concrete pointer (monomorphic) or an interface
			// (polymorphic); either way, comparing the field itself to nil is always correct - no
			// runtime reflection needed, since we know its exact static type here, at codegen time
			fieldExpr := fmt.Sprintf("%s.%s", varName, relationship.GetName())

			if relationship.IsPolymorphic() {
				columns += ", " + relationship.getColumnNameForTargetModel()
				maskPattern += ", " + core.IfThenElse(relationship.isSecret(), "true", "false")
				polyIndex := nbCols
				nbCols++
				//
				argAssignments += fmt.Sprintf(
					"%[1]sargs[base+%[3]d] = goald.GetIDOrNilPlm(%[2]s)\n"+
						"%[1]sargs[base+%[4]d] = goald.GetModelOrNil(%[2]s)\n",
					space, fieldExpr, colIndex, polyIndex)
			} else {
				argAssignments += fmt.Sprintf(
					"%[1]sargs[base+%[3]d] = goald.GetIDOrNil(%[2]s)\n",
					space, fieldExpr, colIndex)
			}
		} else if prop.GetName() == boFieldPreID {
			// the pre-ID field is unexported, so it can only be accessed through its exported getter,
			// unlike the other, regular fields
			columns += prop.getColumnName()
			maskPattern += "false"
			argAssignments += fmt.Sprintf("%sargs[base+%d] = %s.GetPreID()\n", space, colIndex, varName)
		} else {
			columns += prop.getColumnName()
			maskPattern += core.IfThenElse(prop.isSecret(), "true", "false")
			argAssignments += fmt.Sprintf("%sargs[base+%d] = %s.%s\n", space, colIndex, varName, prop.GetName())
		}
	}

	return
}

// ------------------------------------------------------------------------------------------------
//  DAO methods generation
// ------------------------------------------------------------------------------------------------

func (thisServer *server) generateExecCreateQuery(model IBusinessObjectModel) string {

	columns, argAssignments, _, nbCols := thisServer.getColumnsAndInsertData(model, "\t\t\t")
	varName := core.PascalToCamel(string(model.GetName()))

	return fmt.Sprintf(`// ExecCreateQuery implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecCreateQuery(bObjs ...goald.IBusinessObject) (map[int]int64, error) {
	return thisDAO.ExecBatchCreate(&goald.ExecCreateContext{
		Table:       %[2]sDB + `+bq+`.%[4]s`+bq+`,
		Columns:     `+bq+`%[5]s`+bq+`,
		NbCols:      %[7]d,
		MaskPattern: %[2]sMask,
		BObjs:       bObjs,
		FillRow: func(bObj goald.IBusinessObject, args []any, base int) {
			// casting the business object to its actual type
			%[2]s := bObj.(*%[3]s.%[1]s)

%[6]s		},
	})
}`,
		model.GetName(),           // 1
		varName,                   // 2
		model.getPackage(),        // 3
		model.getTableName(false), // 4
		columns,                   // 5
		argAssignments,            // 6
		nbCols,                    // 7
	)
}

// generateExecCreateLinksQueries builds the ExecCreateLinksQueries method for the given model: one
// ExecBatchCreateLink call for every relationship of this model that's persisted through a link table -
// whether this model is on the "source" side (it owns the relationship, e.g. User.MemberOf) or on the
// "target" side (it only has the back-reference, e.g. UserGroup.Members).
//
// Direction matters: the link table's actual column layout (source__..., target__...) is always defined
// by the owning ("source") relationship. So when generating for the target side (e.g. UserGroup), we
// still use the source relationship's table/column names, but flip which side supplies "this business
// object" vs. "each element of the relationship's slice" when building the rows to insert.
//
// Every model gets its own ExecCreateLinksQueries, even if it ends up being just 'return nil' - this
// method is never left to fall back on BusinessObjectDAO's default (which panics), matching how
// ExecCreateQuery is always generated too.
func (thisServer *server) generateExecCreateLinksQueries(model IBusinessObjectModel) string {
	varName := core.PascalToCamel(string(model.GetName()))
	pkg := model.getPackage()

	var blocks string
	for _, relationship := range core.GetSortedValues(model.getRelationships()) {
		// the relationship that actually owns the link table, i.e. the "source" side; and whether we're
		// looking at it from the "target" side (reversed) instead
		linkRel := relationship
		reversed := false

		if !relationship.needsLinkTable() {
			// not persisted at all through a link table, from this model's perspective - either it's a
			// single-valued / column-based relationship, or it's a back-reference whose source side
			// doesn't use a link table (e.g. a plain one-to-many via a foreign key)
			if !relationship.multiple || relationship.backRef == nil || !relationship.backRef.needsLinkTable() {
				continue
			}

			linkRel = relationship.backRef
			reversed = true
		}

		sourceCol, sourceModelCol := linkRel.getLinkTableSourceColumn()
		targetCol, targetModelCol := linkRel.getLinkTableTargetColumn()

		// the columns are always given in the table's actual order: source(s) first, then target(s) -
		// this never changes, regardless of which side we're generating for
		columns := sourceCol
		nbCols := 1
		if sourceModelCol != "" {
			columns += ", " + sourceModelCol
			nbCols++
		}
		columns += ", " + targetCol
		nbCols++
		if targetModelCol != "" {
			columns += ", " + targetModelCol
			nbCols++
		}

		// "own" refers to the business object we're generating this method for; "other" refers to each
		// element of its relationship slice. Depending on the direction, either one can be the link
		// table's source or target
		ownArgs := fmt.Sprintf("%s.GetID()", varName)
		otherArgs := "target.GetID()"
		if !reversed {
			// this business object is the source, the slice elements are the targets
			if sourceModelCol != "" {
				ownArgs += fmt.Sprintf(", %[1]s.GetModelName()", varName)
			}
			if targetModelCol != "" {
				otherArgs += ", target.GetModelName()"
			}
		} else {
			// this business object is the target, the slice elements are the sources
			if sourceModelCol != "" {
				otherArgs += ", target.GetModelName()"
			}
			if targetModelCol != "" {
				ownArgs += fmt.Sprintf(", %[1]s.GetModelName()", varName)
			}
		}

		// building the addRow(...) call's arguments in the same source-then-target order as the columns
		addRowArgs := ownArgs + ", " + otherArgs
		if reversed {
			addRowArgs = otherArgs + ", " + ownArgs
		}

		blocks += fmt.Sprintf(`
	if err := thisDAO.ExecBatchCreateLink(&goald.ExecCreateLinkContext{
		Table:   %[2]sDB + `+bq+`.%[5]s`+bq+`,
		Columns: `+bq+`%[6]s`+bq+`,
		NbCols:  %[7]d,
		BObjs:   bObjs,
		FillRows: func(bObj goald.IBusinessObject, addRow func(args ...any)) {
			%[2]s := bObj.(*%[3]s.%[1]s)
			for _, target := range %[2]s.%[4]s {
				addRow(%[8]s)
			}
		},
	}); err != nil {
		return err
	}
`,
			model.GetName(),            // 1
			varName,                    // 2
			pkg,                        // 3
			relationship.GetName(),     // 4 - the field to iterate, on THIS model
			linkRel.getLinkTableName(), // 5 - the actual link table, always defined by the source side
			columns,                    // 6
			nbCols,                     // 7
			addRowArgs,                 // 8
		)
	}

	return fmt.Sprintf(`// ExecCreateLinksQueries implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecCreateLinksQueries(bObjs ...goald.IBusinessObject) error {%[2]s
	return nil
}`,
		model.GetName(), // 1
		blocks,          // 2
	)
}

// generateExecReadQuery builds the ExecSearchQuery method
func (thisServer *server) generateExecReadQuery(model IBusinessObjectModel) string {
	varName := core.PascalToCamel(string(model.GetName()))
	columns := getColumnsForSelect(model)

	return fmt.Sprintf(`// ------------------------------------------------------------------------------------------------
// Reading
// ------------------------------------------------------------------------------------------------

// ExecReadQuery implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecReadQuery(bObjs map[goald.BObjID]goald.IBusinessObject, bObjIDs []any) error {
	return thisDAO.ExecRead(&goald.ReadContext{
		Table:     %[2]sDB + `+bq+`.%[3]s`+bq+`,
		Columns:   `+bq+`%[4]s`+bq+`,
		Objects:   bObjs,
		ObjectIDs: bObjIDs,
		ScanRow:   scan%[1]sRow,
	})
}`,
		model.GetName(),           // 1
		varName,                   // 2
		model.getTableName(false), // 3
		columns,                   // 4
	)
}

// ------------------------------------------------------------------------------------------------
//  ExecSearchQuery generation: the WHERE-building code mirrors, for every query registered against
//  this model (via Find(...).Where(...) in a *--blo.go file), the OR-ed/AND-ed clause tree built with
//  Either/Or - flattened into its equivalent DNF (a flat list of AND-ed clauses, OR-ed together, by
//  OR's associativity) - plus the row-scanning code for all of the model's persisted columns.
// ------------------------------------------------------------------------------------------------

// generateExecSearchQuery builds the ExecSearchQuery method, the query-args dispatcher + 1 builder
// function per query registered against this model, and the row-scanning function - registering, in
// extraImports, any package (other than this model's own) that this generated code ends up needing.
func (thisServer *server) generateExecSearchQuery(model IBusinessObjectModel, extraImports map[string]bool) string {
	varName := core.PascalToCamel(string(model.GetName()))
	columns := getColumnsForSelect(model)
	queries := getQueriesForModel(model)

	execSelect := fmt.Sprintf(`// ------------------------------------------------------------------------------------------------
// Searching
// ------------------------------------------------------------------------------------------------

// ExecSearchQuery implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecSearchQuery(queryName goald.QueryName, values goald.ISearchParamValues) ([]goald.IBusinessObject, error) {
	return thisDAO.ExecSearch(&goald.SearchContext{
		Table:         %[2]sDB + `+bq+`.%[3]s`+bq+`,
		Columns:       `+bq+`%[4]s`+bq+`,
		QueryName:     queryName,
		QueryValues:   values,
		MakeQueryArgs: build%[1]sQueryArgs,
		ScanRow:       scan%[1]sRow,
	})
}`,
		model.GetName(),           // 1
		varName,                   // 2
		model.getTableName(false), // 3
		columns,                   // 4
	)

	var dispatchCases string
	var builders string

	for i, q := range queries {
		funcName := fmt.Sprintf("build%sQueryArgs%d", model.GetName(), i+1)
		dispatchCases += fmt.Sprintf("\tcase %q:\n\t\treturn %s(values)\n\n", q.getName(), funcName)
		builders += "\n\n" + thisServer.generateQueryArgsBuilder(model, q, funcName, extraImports)
	}

	dispatcher := fmt.Sprintf(`// build%[1]sQueryArgs builds the query args for whichever of %[1]s's registered queries is being run
func build%[1]sQueryArgs(queryName goald.QueryName, values goald.ISearchParamValues) *goald.QueryArgs {
	switch queryName {
%[2]s	default:
		panic(fmt.Sprintf("build%[1]sQueryArgs: unhandled query name '%%s'", queryName))
	}
}`,
		model.GetName(), // 1
		dispatchCases,   // 2
	)

	return execSelect + "\n\n" + dispatcher + builders + "\n\n" + thisServer.generateScanRowFunc(model, extraImports)
}

// getColumnsForSelect returns the comma-separated column list to SELECT for the given model: every
// persisted property, in the same (ID-first, then alphabetical) order used for the table itself -
// except the pre-ID, which is a purely transient, insert-batching helper column with no meaning once
// a row is persisted.
func getColumnsForSelect(model IBusinessObjectModel) string {
	var columns string

	for _, prop := range model.getPersistedProperties() {
		if prop.GetName() == boFieldPreID {
			continue
		}

		if columns != "" {
			columns += ", "
		}

		if relationship, ok := prop.(*Relationship); ok {
			columns += relationship.getColumnName()
			if relationship.IsPolymorphic() {
				columns += ", " + relationship.getColumnNameForTargetModel()
			}
		} else {
			columns += prop.getColumnName()
		}
	}

	return columns
}

// getQueriesForModel returns every query registered (via Find(...).Where(...)) against the given
// model, sorted by (dynamically generated) name, for deterministic codegen output.
func getQueriesForModel(model IBusinessObjectModel) []IQuery {
	queryRegistry.mx.Lock()
	defer queryRegistry.mx.Unlock()

	var queries []IQuery
	for _, q := range queryRegistry.queries {
		if q.getSearchedObjectsModel().GetName() == model.GetName() {
			queries = append(queries, q)
		}
	}

	sort.Slice(queries, func(i, j int) bool { return queries[i].getName() < queries[j].getName() })

	return queries
}

// flatten turns a clause tree (built with Either/Or/comparisons in a *--blo.go file) into its
// equivalent flat list of AND-ed clauses (leaves), OR-ed together.
func flatten(c IClause) [][]IClause {
	switch c.GetType() {
	case clauseTypeOR:
		var result [][]IClause
		for _, sub := range c.GetSubClauses() {
			result = append(result, flatten(sub)...)
		}
		return result

	case clauseTypeAND:
		product := [][]IClause{{}}
		for _, sub := range c.GetSubClauses() {
			var newProduct [][]IClause
			for _, existingGroup := range product {
				for _, subGroup := range flatten(sub) {
					combined := make([]IClause, 0, len(existingGroup)+len(subGroup))
					combined = append(combined, existingGroup...)
					combined = append(combined, subGroup...)
					newProduct = append(newProduct, combined)
				}
			}
			product = newProduct
		}
		return product

	default:
		// a leaf clause (a comparison, or an IN)
		return [][]IClause{{c}}
	}
}

// findQueryParamsModel returns the query params model (e.g. PurchaseOrderQuery) that a query's clauses
// compare against, found via any comparison leaf's right-hand side - the left-hand side is always the
// business object's own property (e.g. order.OrderRef()), the right-hand side the query param being
// compared against it (e.g. query.OrderRefExact()); an IN clause has no right-hand side (its values are
// literal constants), hence looking across every leaf until one is found.
func findQueryParamsModel(dnf [][]IClause) IBusinessObjectModel {
	for _, group := range dnf {
		for _, leaf := range group {
			if concreteLeaf, ok := leaf.(*clause); ok && concreteLeaf.right != nil {
				return concreteLeaf.right.ownerModel()
			}
		}
	}

	return nil
}

// generateQueryArgsBuilder builds the function computing 1 registered query's *goald.QueryArgs (its
// OR-ed/AND-ed condition clauses, args, and secret mask) from that query's param values.
func (thisServer *server) generateQueryArgsBuilder(model IBusinessObjectModel, q IQuery, funcName string, extraImports map[string]bool) string {
	var dnf [][]IClause
	if q.getWhere() != nil {
		dnf = flatten(q.getWhere())
	}

	// applying any registered mandatory clauses so they can never be bypassed - unlike the clauses
	// coded in a *--blo.go file, these are never optional/conditional on a query param's presence:
	// they must always end up in the generated code as an unconditional AND, on every single group,
	// tracked here so generateLeafClauseCode can tell them apart from the regular, optional ones
	mandatory := map[IClause]bool{}
	for _, buildMandatory := range mandatoryClauseBuilders {
		mandatoryClause := buildMandatory(model, q.getQueryParamsModel())
		if mandatoryClause == nil {
			continue
		}

		mandatory[mandatoryClause] = true

		if len(dnf) == 0 {
			dnf = [][]IClause{{mandatoryClause}}
		} else {
			for i := range dnf {
				dnf[i] = append(dnf[i], mandatoryClause)
			}
		}
	}

	queryParamsModel := findQueryParamsModel(dnf)
	extraImports[getImportPackageLine(queryParamsModel)] = true
	queryParamsType := fmt.Sprintf("*%s.%s", queryParamsModel.getPackage(), queryParamsModel.GetName())

	var body strings.Builder
	nbArgs := 0

	for i, group := range dnf {
		if i > 0 {
			body.WriteString("\n\tqueryArgs.NewAndClause()\n\n")
		}

		for _, leaf := range group {
			line, argCount := thisServer.generateLeafClauseCode(leaf, model, extraImports, mandatory[leaf])
			body.WriteString("\t")
			body.WriteString(line)
			body.WriteString("\n")
			nbArgs += argCount
		}
	}

	return fmt.Sprintf(`// %[1]s builds the query args for the %[2]q query
func %[1]s(values goald.ISearchParamValues) *goald.QueryArgs {
	// casting the query params to their actual type
	queryParams := values.(%[3]s)

	// pre-sizing for the worst case (every condition present), to avoid slice reallocations below
	queryArgs := goald.NewQueryArgs(%[4]d, "%[6]s", %[7]t)

%[5]s
	return queryArgs
}`,
		funcName,                             // 1
		q.getName(),                          // 2
		queryParamsType,                      // 3
		nbArgs,                               // 4
		body.String(),                        // 5
		model.getDB().get.QueryPlaceholder(), // 6
		model.getDB().is.QueryPlaceholderIndexed(), // 7
	)
}

// generateLeafClauseCode builds the single statement appending 1 leaf clause (a comparison, or an IN)
// to the AND-clause currently being built, and returns how many query args it consumes. isMandatory
// tells it whether this leaf came from a registered mandatory-clause builder (see
// RegisterMandatoryClauseBuilder), rather than from the *--blo.go file's own Where(...) clauses: unlike
// those, a mandatory clause must always be applied unconditionally (never skipped for lacking a query
// param value) - it's enforced regardless of what the caller did or didn't provide.
func (thisServer *server) generateLeafClauseCode(cl IClause, model IBusinessObjectModel, extraImports map[string]bool, isMandatory bool) (string, int) {
	c := cl.(*clause)
	if c.ctype == clauseTypeIN {
		enumTypeExpr := getNonBuiltInFieldType(model.getType(), c.left.GetName(), extraImports)
		columnName := c.left.getColumnName()

		addExprs := make([]string, len(c.values))
		for i, v := range c.values {
			prefix := ""
			if i == 0 {
				prefix = columnName + " IN ("
			}
			addExprs[i] = fmt.Sprintf("queryArgs.AddSingleClause(%q, %s(%d), false)", prefix, enumTypeExpr, v.Val())
		}

		return fmt.Sprintf(`queryArgs.AppendRawAndClause(true, %s+")")`, strings.Join(addExprs, `+", "+`)), len(c.values)
	}

	columnName := c.left.getColumnName()
	fieldName := c.right.GetName()
	secret := c.right.isSecret()
	operator := string(c.ctype)

	condition := "true"
	if !isMandatory {
		condition = fmt.Sprintf("queryParams.%s != %s", fieldName, zeroValueLiteral(c.right.getPropertyType()))
	}

	valueExpr := fmt.Sprintf("queryParams.%s", fieldName)
	if c.right.getPropertyType() == propertyTypeRELATIONSHIPxMONOM {
		valueExpr += ".GetID()"
	}

	return fmt.Sprintf("queryArgs.AppendAndClause(%s, %q, %s, %t)", condition, columnName+" "+operator+" ", valueExpr, secret), 1
}

// generateScanRowFunc builds the function instantiating and filling in 1 business object from a
// *sql.Rows, in the same column order given to ExecSearchQuery above.
func (thisServer *server) generateScanRowFunc(model IBusinessObjectModel, extraImports map[string]bool) string {
	varName := core.PascalToCamel(string(model.GetName()))
	pkg := model.getPackage()

	var declarations []string
	var scanTargets []string
	var assignments []string

	for _, prop := range model.getPersistedProperties() {
		if prop.GetName() == boFieldPreID {
			continue
		}

		if prop.GetName() == BoFieldID {
			declarations = append(declarations, "id int64")
			scanTargets = append(scanTargets, "&id")
			assignments = append(assignments, fmt.Sprintf("%s.ID = goald.BObjID(id)", varName))
			continue
		}

		fieldName := prop.GetName()
		propVarName := core.PascalToCamel(fieldName)

		if relationship, ok := prop.(*Relationship); ok {
			idVar := propVarName + "ID"
			declarations = append(declarations, fmt.Sprintf("%s sql.NullInt64", idVar))
			scanTargets = append(scanTargets, "&"+idVar)

			if relationship.IsPolymorphic() {
				mdlVar := propVarName + "Mdl"
				declarations = append(declarations, fmt.Sprintf("%s sql.NullString", mdlVar))
				scanTargets = append(scanTargets, "&"+mdlVar)

				targetTypeExpr := getRelationshipFieldType(model.getType(), fieldName, extraImports)
				assignments = append(assignments, fmt.Sprintf(
					"if %[1]s.Valid && %[2]s.Valid {\n\t\t%[3]s.%[4]s = goald.NewBusinessObject(%[2]s.String, %[1]s.Int64).(%[5]s)\n\t}",
					idVar, mdlVar, varName, fieldName, targetTypeExpr))
			} else {
				targetTypeExpr := strings.TrimPrefix(getRelationshipFieldType(model.getType(), fieldName, extraImports), "*")
				assignments = append(assignments, fmt.Sprintf(
					"if %[1]s.Valid {\n\t\t%[2]s.%[3]s = &%[4]s{}\n\t\t%[2]s.%[3]s.ID = goald.BObjID(%[1]s.Int64)\n\t}",
					idVar, varName, fieldName, targetTypeExpr))
			}

			continue
		}

		switch prop.getPropertyType() {
		case propertyTypeBOOL:
			declarations = append(declarations, propVarName+" bool")
			scanTargets = append(scanTargets, "&"+propVarName)
			assignments = append(assignments, fmt.Sprintf("%s.%s = %s", varName, fieldName, propVarName))

		case propertyTypeSTRING:
			declarations = append(declarations, propVarName+" string")
			scanTargets = append(scanTargets, "&"+propVarName)
			assignments = append(assignments, fmt.Sprintf("%s.%s = %s", varName, fieldName, propVarName))

		case propertyTypeINT:
			declarations = append(declarations, propVarName+" int")
			scanTargets = append(scanTargets, "&"+propVarName)
			assignments = append(assignments, fmt.Sprintf("%s.%s = %s", varName, fieldName, propVarName))

		case propertyTypeBIGINT:
			declarations = append(declarations, propVarName+" int64")
			scanTargets = append(scanTargets, "&"+propVarName)
			assignments = append(assignments, fmt.Sprintf("%s.%s = %s", varName, fieldName, propVarName))

		case propertyTypeREAL:
			declarations = append(declarations, propVarName+" float32")
			scanTargets = append(scanTargets, "&"+propVarName)
			assignments = append(assignments, fmt.Sprintf("%s.%s = %s", varName, fieldName, propVarName))

		case propertyTypeDOUBLE:
			declarations = append(declarations, propVarName+" float64")
			scanTargets = append(scanTargets, "&"+propVarName)
			assignments = append(assignments, fmt.Sprintf("%s.%s = %s", varName, fieldName, propVarName))

		case propertyTypeDATE:
			declarations = append(declarations, propVarName+" sql.NullTime")
			scanTargets = append(scanTargets, "&"+propVarName)
			assignments = append(assignments, fmt.Sprintf("if %[1]s.Valid {\n\t\t%[2]s.%[3]s = &%[1]s.Time\n\t}", propVarName, varName, fieldName))

		case propertyTypeENUM:
			enumTypeExpr := getNonBuiltInFieldType(model.getType(), fieldName, extraImports)
			declarations = append(declarations, propVarName+" int")
			scanTargets = append(scanTargets, "&"+propVarName)
			assignments = append(assignments, fmt.Sprintf("%s.%s = %s(%s)", varName, fieldName, enumTypeExpr, propVarName))
		}
	}

	return fmt.Sprintf(`// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

// scan%[1]sRow instantiates and fills in one %[1]s from the current row
func scan%[1]sRow(rows *sql.Rows, bObjs map[goald.BObjID]goald.IBusinessObject) (goald.IBusinessObject, error) {
	var (
		%[3]s
	)

	if errScan := rows.Scan(
		%[4]s); errScan != nil {
		return nil, errScan
	}

	// instantiating the business object, either from the provided map or as a new one
	var %[2]s *%[5]s.%[1]s
	if bObjs != nil {
		%[2]s = bObjs[goald.BObjID(id)].(*%[5]s.%[1]s)
	} else {
		%[2]s := &%[5]s.%[1]s{}
		%[2]s.ID = goald.BObjID(id)
	}

	// filling in the business object with the scanned values
	%[6]s

	// returning the filled-in business object
	return %[2]s, nil
}`,
		model.GetName(),                      // 1
		varName,                              // 2
		strings.Join(declarations, "\n\t\t"), // 3
		strings.Join(scanTargets, ",\n\t\t"), // 4
		pkg,                                  // 5
		strings.Join(assignments, "\n\t"),    // 6
	)
}
