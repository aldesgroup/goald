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

// ------------------------------------------------------------------------------------------------
// Useful structs & constants
// ------------------------------------------------------------------------------------------------

const daoINITxTEMPLATE = `// Generated file, do not edit!
package %[1]s

import (
	"database/sql"
	"fmt"
	"sync"

	$$extraimports$$
)

// ------------------------------------------------------------------------------------------------
// %[2]sDAO declaration and registration
// ------------------------------------------------------------------------------------------------

type %[2]sDAO struct {
	goald.BusinessObjectDAO
}

// type check
var _ goald.IBusinessObjectDAO = (*%[2]sDAO)(nil)

// registration
func init() {
	goald.RegisterDAO("%[2]s", &%[2]sDAO{})
}

// NewDAO implements [goald.IBusinessObjectDAO].
func (thisDAO *%[2]sDAO) NewDAO() goald.IBusinessObjectDAO {
	%[3]sDBOnce.Do(func() {
		%[3]sDB = string(goald.GetDB("%[4]s").Name())
	})

	return &%[2]sDAO{}
}

// ------------------------------------------------------------------------------------------------
// Local state
// ------------------------------------------------------------------------------------------------

var (
	%[3]sDB     string
	%[3]sDBOnce sync.Once
	%[3]sMask   = []bool{%[5]s}
)
`

// used to insert a literal backtick inside a raw string template, since raw strings can't escape backticks
const bq = "`"

const daoFILExSUFFIX = "--dao.go"
const daoFILExSUFFIXxLEN = len(daoFILExSUFFIX)

// daoGenerator is responsible for generating DAO files for all the business objects.
type daoGenerator struct {
	*server
}

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

	// the generator we instantiate to benefit from its methods
	daoGen := &daoGenerator{
		server: thisServer,
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
				daoGen.generateOneDAO(srcdir, string(dbType), boModel)

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
	// 		thisGen.Info(fmt.Sprintf("removing %s", unneededDAO.filename))
	// 		if errRem := os.Remove(path.Join(daoDir, unneededDAO.filename)); errRem != nil {
	// 			core.PanicMsgIfErr(errRem, "Could not delete DAO file '%s'", unneededDAO.filename)
	// 		}
	// 	}

	return
}

func (thisGen *daoGenerator) generateOneDAO(srcdir string, dbType string, model IBusinessObjectModel) bool {
	// where the DAO should end up
	daoDir := core.EnsureDir(srcdir, includePATH, dbFOLDERNAME, dbType)

	// trivial filling of the template
	modelNameCamel := core.PascalToCamel(string(model.GetName()))
	_, _, maskPattern, _ := thisGen.getColumnsAndInsertData(model, "\t\t\t")

	// dealing with extra imports, and making sure the model at least is imported
	extraImports := map[string]bool{}
	extraImports["github.com/aldesgroup/goald"] = true
	extraImports[getImportPackageLine(model)] = true

	// init of the content
	content := fmt.Sprintf(daoINITxTEMPLATE,
		dbType,
		model.GetName(),
		modelNameCamel,
		model.getDB().name,
		maskPattern)

	// adding all the needed DAO methods
	content += thisGen.generateExecCreateQuery(model)
	content += "\n\n"
	content += thisGen.generateExecCreateLinksQueries(model)
	content += "\n\n"
	content += thisGen.generateExecReadQuery(model)
	content += "\n\n"
	content += thisGen.generateExecReadRelationshipQuery(model, extraImports)
	content += "\n"
	content += thisGen.generateExecSearchQuery(model, extraImports)
	content += "\n\n"
	content += thisGen.generateScanRowFunc(model, extraImports)

	// adding the imports
	extraImportsString := ""
	for _, importPath := range core.GetSortedKeys(extraImports) {
		extraImportsString += fmt.Sprintf("\n\t%q", importPath)
	}
	content = strings.Replace(content, "$$extraimports$$", extraImportsString, 1)

	// writing to file
	core.WriteToFile(content, daoDir, core.PascalToKebab(string(model.GetName()))+daoFILExSUFFIX)

	thisGen.Info(fmt.Sprintf("(Re-)generated DAO for %s", model.GetName()))

	return true
}

// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

// getColumnsAndInsertData gathers the necessary information for generating an insert statement for the given model.
func (thisGen *daoGenerator) getColumnsAndInsertData(model IBusinessObjectModel, space string) (columns, argAssignments, maskPattern string, nbCols int) {
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
// Generating the CRUD(S) methods - the CREATION part
// ------------------------------------------------------------------------------------------------

func (thisGen *daoGenerator) generateExecCreateQuery(model IBusinessObjectModel) string {

	columns, argAssignments, _, nbCols := thisGen.getColumnsAndInsertData(model, "\t\t\t")
	varName := core.PascalToCamel(string(model.GetName()))

	return fmt.Sprintf(`
// ------------------------------------------------------------------------------------------------
// Insertion
// ------------------------------------------------------------------------------------------------

// ExecCreateQuery implements [goald.IBusinessObjectDAO].
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

// generateExecCreateLinksQueries builds the ExecCreateLinksQueries method for the given model
func (thisGen *daoGenerator) generateExecCreateLinksQueries(model IBusinessObjectModel) string {
	varName := core.PascalToCamel(string(model.GetName()))
	pkg := model.getPackage()

	var blocks string
	for _, relationship := range core.GetSortedValues(model.getRelationships()) {
		// the relationship that actually owns the link table
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
			relationship.GetName(),     // 4
			linkRel.getLinkTableName(), // 5
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

// ------------------------------------------------------------------------------------------------
// Generating the CRUD(S) methods - the READING part
// ------------------------------------------------------------------------------------------------

// generateExecReadQuery builds the ExecSearchQuery method
func (thisGen *daoGenerator) generateExecReadQuery(model IBusinessObjectModel) string {
	varName := core.PascalToCamel(string(model.GetName()))
	columns := getColumnsForSelect(model)

	return fmt.Sprintf(`// ------------------------------------------------------------------------------------------------
// Reading
// ------------------------------------------------------------------------------------------------

// ExecReadQuery implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecReadQuery(bObjIDs []any, cache *goald.BObjCache) error {
	return thisDAO.ExecRead(&goald.ReadContext{
		Table:     %[2]sDB + `+bq+`.%[3]s`+bq+`,
		Columns:   `+bq+`%[4]s`+bq+`,
		ObjectIDs: bObjIDs,
		BoCache:   cache,
		ScanRow:   scan%[1]sRow,
	})
}`,
		model.GetName(),           // 1
		varName,                   // 2
		model.getTableName(false), // 3
		columns,                   // 4
	)
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

// ------------------------------------------------------------------------------------------------
// Generating the CRUD(S) methods - the READING RELATIONSHIPS part
// ------------------------------------------------------------------------------------------------

const execReadRelationshipQueryTpl = `// ------------------------------------------------------------------------------------------------
// Reading relationships
// ------------------------------------------------------------------------------------------------

// ExecReadRelationshipQuery implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecReadRelationshipQuery(bObjIDs []any, relName string, cache *goald.BObjCache) ([]goald.IBusinessObject, error) {
	switch relName {
%[2]s
	default:
		panic(fmt.Sprintf("%[1]sDAO.ExecReadRelationshipQuery: unhandled relationship name '%%s'", relName))
	}
}
%[3]s`

const execReadOneRelationshipTpl = `
// execRead%[2]s reads the "%[2]s" relationship
func (thisDAO *%[1]sDAO) execRead%[2]s(bObjIDs []any, cache *goald.BObjCache) ([]goald.IBusinessObject, error) {
	inClause, queryArgs := goald.NewQueryArgsWithINClause("%[3]s", %[4]t, bObjIDs)

	return thisDAO.ExecReadRelationship(&goald.ReadRelationshipContext{
		Query: "SELECT %[5]s" +
			" FROM " + %[6]sDB + ".%[7]s" +
			" WHERE %[8]s" + inClause,
		Args: queryArgs.Args,
		AttachRow: func(rows *sql.Rows) (goald.IBusinessObject, error) {
			var sourceID, targetID goald.BObjID%[9]s
			if errScan := rows.Scan(%[10]s); errScan != nil {
				return nil, errScan
			}
			return %[11]s.Get%[1]sFrom(cache, sourceID).WithAdded%[2]s(%[12]s.CachedOrNew%[13]s(%[14]s, targetID)%[15]s), nil
		},
	})
}
`

// Implement the logic to generate read relationship cases and functions
func (thisGen *daoGenerator) generateExecReadRelationshipQuery(model IBusinessObjectModel, extraImports map[string]bool) string {
	var readRelCases, readRelFuncs string

	// iterating over the relationships...
	for i, rel := range core.GetSortedValues(model.getRelationships()) {
		// ... but only the multiple ones, since the single-valued ones are read through the main ExecReadQuery method
		if rel.IsMultiple() {
			// name of the model in camel case
			modelNameCamel := core.PascalToCamel(string(model.GetName()))

			// declaring some variables that will be used in the template
			var sourceColName, targetColName, targetMdlColName, selTable string

			// is this link indirectly persisted, from this model's perspective?
			if rel.isDirectlyPersisted() {
				sourceColName, _ = rel.getLinkTableSourceColumn()
				targetColName, targetMdlColName = rel.getLinkTableTargetColumn()
				selTable = rel.getLinkTableName()
			} else {
				// in this case, the backref is the directly persisted one, so things are a bit reversed here
				if dirRel := rel.backRef; dirRel.IsMultiple() {
					sourceColName, _ = dirRel.getLinkTableTargetColumn() // the backref's target is this model, so it's the source for the link table
					targetColName, _ = dirRel.getLinkTableSourceColumn() // the backref's source is the other model, so it's the target for the link table
					selTable = dirRel.getLinkTableName()
				} else {
					sourceColName = dirRel.getColumnName() // the backref's target is this model, so it's the source for this relationship
					targetColName = "id"                   // the backref's source is the other model, so it's the target for this relationship
					selTable = dirRel.owner.getTableName(false)
				}
			}

			// aggregating the SELECT and the scanning parts
			selColNames := sourceColName + ", " + targetColName
			if targetMdlColName != "" {
				selColNames += ", " + targetMdlColName
			}
			selColVars := "&sourceID, &targetID"
			if rel.IsPolymorphic() {
				selColVars += ", &targetMdl"
			}

			// declaring iother variables that will be used in the template
			var mdlColScan, cacheOrNewReceiver, cacheOrNewObject, cacheOrNewArg, cacheOrNewType string

			// Is the model concrete, or an interface?
			if rel.IsPolymorphic() { // let's say it's an interface
				mdlColScan = newline + `			var targetMdl utils.ModelName`
				cacheOrNewReceiver = "cache"
				cacheOrNewObject = "BusinessObject"
				cacheOrNewArg = "targetMdl"
				cacheOrNewType = fmt.Sprintf(".(%s)", getRelationshipFieldType(model.getType(), rel.GetName(), extraImports))

				// let's not forget this then:
				extraImports["github.com/aldesgroup/goald/features/utils"] = true
			} else { // or if it's a concrete model
				cacheOrNewReceiver = rel.getUniqueTargetModel().getPackage()
				cacheOrNewObject = string(rel.getUniqueTargetName())
				cacheOrNewArg = "cache"
			}

			// Adding the function:
			readRelFuncs += fmt.Sprintf(execReadOneRelationshipTpl,
				model.GetName(),                            //  1: func (thisDAO *%[1]sDAO)
				rel.GetName(),                              //  2: execRead%[2]s reads the "%[2]s" relationship
				model.getDB().get.QueryPlaceholder(),       //  3: NewQueryArgsWithINClause("%[3]s", %[4]s, bObjIDs)
				model.getDB().is.QueryPlaceholderIndexed(), //  4: NewQueryArgsWithINClause("%[3]s", %[4]s, bObjIDs)
				selColNames,                                //  5: SELECT %[5]s
				modelNameCamel,                             //  6: FROM ` + bq + ` + %[6]s + ` + bq + `.%[7]s
				selTable,                                   //  7: FROM ` + bq + ` + %[6]s + ` + bq + `.%[7]s
				sourceColName,                              //  8: WHERE %[8]s` + bq + ` + inClause
				mdlColScan,                                 //  9: var sourceID, targetID goald.BObjID%[9]s
				selColVars,                                 // 10: if errScan := rows.Scan(%[10]s); errScan != nil {
				model.getPackage(),                         // 11: return %[11]s.Get%[1]sFrom(cache, sourceID)
				cacheOrNewReceiver,                         // 12: .WithAdded%[2]s(%[12]s.CachedOrNew%[13]s(%[14]s, targetID)%[15]s)
				cacheOrNewObject,                           // 13: .WithAdded%[2]s(%[12]s.CachedOrNew%[13]s(%[14]s, targetID)%[15]s)
				cacheOrNewArg,                              // 14: .WithAdded%[2]s(%[12]s.CachedOrNew%[13]s(%[14]s, targetID)%[15]s)
				cacheOrNewType,                             // 15: .WithAdded%[2]s(%[12]s.CachedOrNew%[13]s(%[14]s, targetID)%[15]s)
			)

			// Adding the case:
			readRelCases += fmt.Sprintf(`	case "%[1]s":
		return thisDAO.execRead%[1]s(bObjIDs, cache)
`+core.IfThenElse(i+1 == len(model.getRelationships()), "", newline), rel.GetName())
		}
	}

	return fmt.Sprintf(execReadRelationshipQueryTpl, model.GetName(), readRelCases, readRelFuncs)
}

// ------------------------------------------------------------------------------------------------
// Generating the CRUD(S) methods - the SEARCH part
// ------------------------------------------------------------------------------------------------

// generateExecSearchQuery builds the ExecSearchQuery method, the query-args dispatcher
// + 1 builder function per query registered against this model
func (thisGen *daoGenerator) generateExecSearchQuery(model IBusinessObjectModel, extraImports map[string]bool) string {
	varName := core.PascalToCamel(string(model.GetName()))
	columns := getColumnsForSelect(model)
	queries := getQueriesForModel(model)

	execSelect := fmt.Sprintf(`// ------------------------------------------------------------------------------------------------
// Searching
// ------------------------------------------------------------------------------------------------

// ExecSearchQuery implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecSearchQuery(queryName goald.QueryName, values goald.ISearchParamValues, cache *goald.BObjCache) ([]goald.IBusinessObject, error) {
	return thisDAO.ExecSearch(&goald.SearchContext{
		Table:         %[2]sDB + `+bq+`.%[3]s`+bq+`,
		Columns:       `+bq+`%[4]s`+bq+`,
		QueryName:     queryName,
		QueryValues:   values,
		MakeQueryArgs: build%[1]sQueryArgs,
		BoCache:       cache,
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
		builders += "\n\n" + thisGen.generateQueryArgsBuilder(model, q, funcName, extraImports)
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

	return execSelect + "\n\n" + dispatcher + builders
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

// generateQueryArgsBuilder builds the function computing 1 registered query's *goald.QueryArgs (its
// OR-ed/AND-ed condition clauses, args, and secret mask) from that query's param values.
func (thisGen *daoGenerator) generateQueryArgsBuilder(model IBusinessObjectModel, q IQuery, funcName string, extraImports map[string]bool) string {
	// flatten the WHERE clause into disjunctive normal form (DNF)
	var dnf [][]IClause
	if q.getWhere() != nil {
		dnf = q.getWhere().toDNF()
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

	// which query params model this query is using
	queryParamsModel := thisGen.findQueryParamsModel(dnf)

	// making sure we import the package for the query params model
	extraImports[getImportPackageLine(queryParamsModel)] = true

	// making a string representing the query params type
	queryParamsType := fmt.Sprintf("*%s.%s", queryParamsModel.getPackage(), queryParamsModel.GetName())

	// building the clauses
	var body strings.Builder
	nbArgs := 0

	for i, group := range dnf {
		if i > 0 {
			body.WriteString("\n\tqueryArgs.NewAndClause()\n\n")
		}

		for _, leaf := range group {
			line, argCount := thisGen.generateLeafClauseCode(leaf, model, extraImports, mandatory[leaf])
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

// findQueryParamsModel returns the query params model that a query's clause is comparing against.
// found via any comparison leaf's right-hand side - the left-hand side is always the
// business object's own property (e.g. order.OrderRef()), the right-hand side the query param being
// compared against it (e.g. query.OrderRefExact()); an IN clause has no right-hand side (its values are
// literal constants), hence looking across every leaf until one is found.
func (thisGen *daoGenerator) findQueryParamsModel(dnf [][]IClause) IBusinessObjectModel {
	for _, group := range dnf {
		for _, leaf := range group {
			if concreteLeaf, ok := leaf.(*clause); ok && concreteLeaf.right != nil {
				return concreteLeaf.right.ownerModel()
			}
		}
	}

	return nil
}

// generateLeafClauseCode builds the single statement appending 1 leaf clause (a comparison, or an IN)
// to the AND-clause currently being built, and returns how many query args it consumes. isMandatory
// tells it whether this leaf came from a registered mandatory-clause builder (see
// RegisterMandatoryClauseBuilder), rather than from the *--blo.go file's own Where(...) clauses: unlike
// those, a mandatory clause must always be applied unconditionally (never skipped for lacking a query
// param value) - it's enforced regardless of what the caller did or didn't provide.
func (thisGen *daoGenerator) generateLeafClauseCode(cl IClause, model IBusinessObjectModel, extraImports map[string]bool, isMandatory bool) (string, int) {
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

// ------------------------------------------------------------------------------------------------
// Generating the CRUD(S) methods - the main row scanning function
// ------------------------------------------------------------------------------------------------

// generateScanRowFunc builds the function instantiating and filling in 1 business object from a
// *sql.Rows, in the same column order given to ExecSearchQuery above.
func (thisGen *daoGenerator) generateScanRowFunc(model IBusinessObjectModel, extraImports map[string]bool) string {
	varName := core.PascalToCamel(string(model.GetName()))
	pkg := model.getPackage()

	var declarations []string
	var scanTargets []string
	var assignments []string

	for _, prop := range model.getPersistedProperties() {
		if prop.GetName() == boFieldPreID {
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
				assignments = append(assignments, fmt.Sprintf(`		if %[1]s.Valid && %[2]s.Valid {
			%[3]s.%[4]s = cache.CachedOrNewBusinessObjectRaw(%[2]s.String, %[1]s.Int64).(%[5]s)
		}`,
					idVar, mdlVar, varName, fieldName, targetTypeExpr))
			} else {
				targetTypeExprParts := strings.Split(strings.TrimPrefix(getRelationshipFieldType(model.getType(), fieldName, extraImports), "*"), ".")
				pkg := targetTypeExprParts[0]
				targetType := targetTypeExprParts[1]
				assignments = append(assignments, fmt.Sprintf(`		if %[1]s.Valid {
			%[3]s.%[4]s = %[2]s.CachedOrNew%[5]s(cache, goald.BObjID(%[1]s.Int64))
		}`,
					idVar, pkg, varName, fieldName, targetType))
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
			if fieldName == BoFieldID {
				declarations = append(declarations, propVarName+" goald.BObjID")
				// no assignment for the ID
			} else {
				declarations = append(declarations, propVarName+" int64")
				assignments = append(assignments, fmt.Sprintf("%s.%s = %s", varName, fieldName, propVarName))
			}
			scanTargets = append(scanTargets, "&"+propVarName)

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
func scan%[1]sRow(rows *sql.Rows, cache *goald.BObjCache) (goald.IBusinessObject, error) {
	var (
		%[3]s
	)

	if errScan := rows.Scan(
		%[4]s); errScan != nil {
		return nil, errScan
	}

	// retrieving the object from the cache or creating a new one
	%[2]s := %[5]s.CachedOrNew%[1]s(cache, id)

	// filling in the business object with the scanned values, if it was newly created
	if %[2]s.Creation == nil {
	%[6]s
    }

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
