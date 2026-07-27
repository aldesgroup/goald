// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the class files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
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

	// we'll gather all the existing class files, per DB type
	existingDAOFiles := map[dbconn.DatabaseType]map[className]*daoFile{}
	for dbType := range dbTypes {
		// some init
		existingDAOFiles[dbType] = map[className]*daoFile{}

		// making sure the DAO folder exists for this DB type
		daoDir := core.EnsureDir(srcdir, includePATH, dbFOLDERNAME, string(dbType))

		// so, let's read this folder now
		for _, daoEntry := range core.EnsureReadDir(daoDir) {
			daoEntryInfo, errInfo := daoEntry.Info()
			core.PanicMsgIfErr(errInfo, "Could not read info for file '%s'", daoEntry.Name())
			daoClassName := className(core.KebabToPascal(daoEntry.Name()[:len(daoEntry.Name())-daoFILExSUFFIXxLEN]))
			existingDAOFiles[dbType][daoClassName] = &daoFile{
				modTime:  daoEntryInfo.ModTime(),
				filename: daoEntry.Name(),
			}
		}
	}

	// let's now generate all the DAOs we need in these DB folders
	for name, class := range classRegistry.items {
		// we only consider the business objects
		if !class.isInterface() {
			// obviously we don't persist abstract business objects, so we skip them
			if boModel := modelForName(name); !boModel.base().abstract && boModel.base().isPersistedHere() {
				// what's the DB type involved here?
				dbType := boModel.base().db.schema.DbServer.Type

				// do we meed to generate a DAO for this business object ?
				if existingDAO := existingDAOFiles[dbType][name]; regen ||
					existingDAO == nil || existingDAO.modTime.Before(getClass(boModel).getLastBOMod()) {

					// generating the missing or outdated class
					thisServer.generateOneDAO(srcdir, string(dbType), boModel)

					// code has been changed
					codeChanged = true

					// keeping track of this class package needing to import the DAO package in its registry file
					if _, ok := daoFolders[class.getSrcPath()]; !ok {
						daoFolders[class.getSrcPath()] = map[dbconn.DatabaseType]bool{}
					}
					daoFolders[class.getSrcPath()][dbType] = true
				}

				// flagging this business object class as NOT unneeded (i.e. needed)
				delete(existingDAOFiles[dbType], name)
			}
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

	// 	// removing the unneeded classes
	// 	for _, unneededDAO := range existingDAOFiles {
	// 		thisServer.Info(fmt.Sprintf("removing %s", unneededDAO.filename))
	// 		if errRem := os.Remove(path.Join(daoDir, unneededDAO.filename)); errRem != nil {
	// 			core.PanicMsgIfErr(errRem, "Could not delete class file '%s'", unneededDAO.filename)
	// 		}
	// 	}

	// let's make the registry file import the DAO package, if it doesn't already
	if codeChanged {
		for srcPath := range daoFolders {
			for dbType := range daoFolders[srcPath] {
				filename := path.Join(srcdir, includePATH, srcPath, sourceREGISTRYxNAME)
				if core.FileExists(filename) {
					daoFolderPath := "_ \"" + path.Join(getCurrentModule(), includePATH, dbFOLDERNAME, string(dbType)) + "\""
					if _, line := core.FindLineInFile(filename, func(line string) bool { return strings.Contains(line, "/"+daoFolderPath) }, false); line == 0 {
						core.ReplaceInFile(filename, map[string]string{goaldIMPORT: goaldIMPORT + newline + daoFolderPath})
					}
				}
			}
		}
	}

	return
}

func (thisServer *server) generateOneDAO(srcdir string, dbType string, model IBusinessObjectModel) bool {
	// if model.base().name != "StaffMember" {
	// 	return false
	// }

	////

	// where the DAO should end up
	daoDir := core.EnsureDir(srcdir, includePATH, dbFOLDERNAME, dbType)

	// trivial filling of the template
	class := getClass(model)
	classCamel := core.PascalToCamel(string(class.getClassName()))
	importForClass := getImportPackageLine(class)
	content := fmt.Sprintf(daoINITxTEMPLATE, dbType, importForClass, class.getClassName(), classCamel, model.base().db.name)

	// adding all the needed DAO methods
	content += thisServer.generateExecInsertQuery(model)

	// writing to file
	core.WriteToFile(content, daoDir, core.PascalToKebab(string(class.getClassName()))+daoFILExSUFFIX)

	thisServer.Info(fmt.Sprintf("(Re-)generated DAO for %s", class.getClassName()))

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
	varName := core.PascalToCamel(string(model.base().name))

	for _, prop := range model.base().getPersistedProperties() {
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
			argAssignments += fmt.Sprintf("%sargs[base+%d] = %s.%s.GetID()\n", space, colIndex, varName, relationship.GetName())

			if relationship.IsPolymorphic() {
				columns += ", " + relationship.getColumnNameForTargetClass()
				maskPattern += ", " + core.IfThenElse(relationship.isSecret(), "true", "false")
				polyIndex := nbCols
				nbCols++
				argAssignments += fmt.Sprintf("%[1]sargs[base+%[2]d] = %[3]s.%[4]s.GetClassName(%[3]s.%[4]s)\n", space, polyIndex, varName, relationship.GetName())
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
	className := model.base().name
	varName := core.PascalToCamel(string(className))

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
		className,                    // 1
		varName,                      // 2
		getClass(model).getPackage(), // 3
		model.getTableName(false),    // 4
		columns,                      // 5
		argAssignments,               // 6
		maskPattern,                  // 7
		nbCols,                       // 8
	)
}

//
