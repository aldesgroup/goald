// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the DAO files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/utils"
)

const daoINITxTEMPLATE = `// Generated file, do not edit!
package %[1]s

import (
	"sync"

	"%[2]s"
	"github.com/aldesgroup/goald"
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
)

// ------------------------------------------------------------------------------------------------
// Queries implementation
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
	modelNameCamel := core.PascalToCamel(string(model.getName()))
	importForBOModel := getImportPackageLine(model)
	content := fmt.Sprintf(daoINITxTEMPLATE, dbType, importForBOModel, model.getName(), modelNameCamel, model.getDB().name)

	// adding all the needed DAO methods
	content += thisServer.generateExecInsertQuery(model)
	content += "\n\n"
	content += thisServer.generateExecInsertLinksQueries(model)

	// writing to file
	core.WriteToFile(content, daoDir, core.PascalToKebab(string(model.getName()))+daoFILExSUFFIX)

	thisServer.Info(fmt.Sprintf("(Re-)generated DAO for %s", model.getName()))

	return true
}

// ------------------------------------------------------------------------------------------------
//  DAO utils
// ------------------------------------------------------------------------------------------------

// getColumnsAndInsertData gathers everything needed to generate a batched, multi-row insert for the given
// model: the column list, the per-row 'args[base+N] = ...' assignment statements (for use inside a
// BatchInsertContext.FillRow closure), the masked-column pattern (for hiding secret values from the logs),
// and the total number of columns involved.
func (thisServer *server) getColumnsAndInsertData(model IBusinessObjectModel, space string) (columns, argAssignments, maskPattern string, nbCols int) {
	varName := core.PascalToCamel(string(model.getName()))

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
				argAssignments += fmt.Sprintf("%[1]sif %[2]s != nil {\n"+
					"%[1]s\targs[base+%[3]d] = %[2]s.GetID()\n"+
					"%[1]s\targs[base+%[4]d] = %[2]s.GetModelName()\n"+
					"%[1]s} else {\n"+
					"%[1]s\targs[base+%[3]d] = nil\n"+
					"%[1]s\targs[base+%[4]d] = nil\n"+
					"%[1]s}\n",
					space, fieldExpr, colIndex, polyIndex)
			} else {
				argAssignments += fmt.Sprintf("%[1]sif %[2]s != nil {\n"+
					"%[1]s\targs[base+%[3]d] = %[2]s.GetID()\n"+
					"%[1]s} else {\n"+
					"%[1]s\targs[base+%[3]d] = nil\n"+
					"%[1]s}\n",
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

func (thisServer *server) generateExecInsertQuery(model IBusinessObjectModel) string {

	columns, argAssignments, maskPattern, nbCols := thisServer.getColumnsAndInsertData(model, "\t\t\t")
	varName := core.PascalToCamel(string(model.getName()))

	return fmt.Sprintf(`// ExecInsertQuery implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecInsertQuery(bObjs ...goald.IBusinessObject) (map[int]int64, error) {
	return thisDAO.ExecBatchInsert(&goald.BatchInsertContext{
		Table:       %[2]sDB + `+bq+`.%[4]s`+bq+`,
		Columns:     `+bq+`%[5]s`+bq+`,
		NbCols:      %[8]d,
		MaskPattern: []bool{%[7]s},
		BObjs:       bObjs,
		FillRow: func(bObj goald.IBusinessObject, args []any, base int) {
			// casting the business object to its actual type
			%[2]s := bObj.(*%[3]s.%[1]s)

%[6]s		},
	})
}`,
		model.getName(),           // 1
		varName,                   // 2
		model.getPackage(),        // 3
		model.getTableName(false), // 4
		columns,                   // 5
		argAssignments,            // 6
		maskPattern,               // 7
		nbCols,                    // 8
	)
}

// generateExecInsertLinksQueries builds the ExecInsertLinksQueries method for the given model: one
// ExecBatchLinkInsert call for every relationship of this model that's persisted through a link table -
// whether this model is on the "source" side (it owns the relationship, e.g. User.MemberOf) or on the
// "target" side (it only has the back-reference, e.g. UserGroup.Members).
//
// Direction matters: the link table's actual column layout (source__..., target__...) is always defined
// by the owning ("source") relationship. So when generating for the target side (e.g. UserGroup), we
// still use the source relationship's table/column names, but flip which side supplies "this business
// object" vs. "each element of the relationship's slice" when building the rows to insert.
//
// Every model gets its own ExecInsertLinksQueries, even if it ends up being just 'return nil' - this
// method is never left to fall back on BusinessObjectDAO's default (which panics), matching how
// ExecInsertQuery is always generated too.
func (thisServer *server) generateExecInsertLinksQueries(model IBusinessObjectModel) string {
	varName := core.PascalToCamel(string(model.getName()))
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
	if err := thisDAO.ExecBatchLinkInsert(&goald.LinkInsertContext{
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
			model.getName(),            // 1
			varName,                    // 2
			pkg,                        // 3
			relationship.GetName(),     // 4 - the field to iterate, on THIS model
			linkRel.getLinkTableName(), // 5 - the actual link table, always defined by the source side
			columns,                    // 6
			nbCols,                     // 7
			addRowArgs,                 // 8
		)
	}

	return fmt.Sprintf(`// ExecInsertLinksQueries implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecInsertLinksQueries(bObjs ...goald.IBusinessObject) error {%[2]s
	return nil
}`,
		model.getName(), // 1
		blocks,          // 2
	)
}

//
