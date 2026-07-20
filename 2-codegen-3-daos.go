// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the class files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
)

const daoINITxTEMPLATE = `// Generated file, do not edit!
package dao

import (
	"%[1]s"
	"github.com/aldesgroup/goald"
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
	return &%[2]sDAO{}
}

// ------------------------------------------------------------------------------------------------
// Queries implementation
// ------------------------------------------------------------------------------------------------

`

// used to insert a literal backtick inside a raw string template, since raw strings can't escape backticks
const bq = "`"

const daoFOLDERxNAME = "dao"

const daoFILExSUFFIX = "--dao.go"
const daoFILExSUFFIXxLEN = len(daoFILExSUFFIX)

// ------------------------------------------------------------------------------------------------
// Main DAO files generation methods
// ------------------------------------------------------------------------------------------------

func (thisServer *server) generateAllObjectDAOs(srcdir string, regen bool) (codeChanged bool) {
	regen = true

	// a type just used here
	type daoFile struct {
		modTime  time.Time
		filename string
	}

	// iterating over each package for which we've already got a registry
	for _, includeDirEntry := range core.EnsureReadDir(srcdir, includePATH) {
		// if it's not a directory, we skip it
		if !includeDirEntry.IsDir() || includeDirEntry.Name() == dbFOLDERNAME {
			continue
		}

		// where the DAO files will be generated
		daoDir := core.EnsureDir(srcdir, includePATH, includeDirEntry.Name(), daoFOLDERxNAME)

		// we'll gather all the existing class files
		existingDAOFiles := map[className]*daoFile{}

		// so, let's read the DAO folders
		for _, daoEntry := range core.EnsureReadDir(daoDir) {
			daoEntryInfo, errInfo := daoEntry.Info()
			core.PanicMsgIfErr(errInfo, "Could not read info for file '%s'", daoEntry.Name())
			daoClassName := className(core.KebabToPascal(daoEntry.Name()[:len(daoEntry.Name())-daoFILExSUFFIXxLEN]))
			existingDAOFiles[daoClassName] = &daoFile{
				modTime:  daoEntryInfo.ModTime(),
				filename: daoEntry.Name(),
			}
		}

		// let's see what we have in terms of business objects
		for name, class := range classRegistry.items {
			// we only consider the business objects of THIS module, and no interface (at least for now)
			if class.isFromDir(includeDirEntry.Name()) && !class.isInterface() {
				// obviously we don't persist abstract business objects, so we skip them
				if boModel := modelForName(name); !boModel.base().abstract && !boModel.isNotPersisted() {
					// do we meed to generate a DAO for this business object ?
					if existingDAO := existingDAOFiles[name]; regen ||
						existingDAO == nil || existingDAO.modTime.Before(getClass(boModel).getLastBOMod()) {
						// generating the missing or outdated class
						codeChanged = thisServer.generateOneDAO(daoDir, boModel) || codeChanged
					}

					// flagging this business object class as NOT unneeded (i.e. needed)
					delete(existingDAOFiles, name)
				}
			}
		}

		// removing the unneeded classes
		for _, unneededDAO := range existingDAOFiles {
			thisServer.Info(fmt.Sprintf("removing %s", unneededDAO.filename))
			if errRem := os.Remove(path.Join(daoDir, unneededDAO.filename)); errRem != nil {
				core.PanicMsgIfErr(errRem, "Could not delete class file '%s'", unneededDAO.filename)
			}
		}

		// let's make the registry file import the DAO package, if it doesn't already
		if codeChanged {
			filename := path.Join(srcdir, includePATH, includeDirEntry.Name(), sourceREGISTRYxNAME)
			daoImportPath := "_ \"" + path.Join(getCurrentModule(), includePATH, includeDirEntry.Name(), daoFOLDERxNAME) + "\""
			if _, line := core.FindLineInFile(filename, func(line string) bool { return strings.Contains(line, "/"+daoFOLDERxNAME) }, false); line == 0 {
				core.ReplaceInFile(filename, map[string]string{goaldIMPORT: goaldIMPORT + newline + daoImportPath})
			}
		}
	}

	return
}

func (thisServer *server) generateOneDAO(daoDir string, model IBusinessObjectModel) bool {
	// if model.base().name != "StaffMember" {
	// 	return false
	// }

	// trivial filling of the template
	class := getClass(model)
	importForClass := path.Join(getCurrentModule(), class.getSrcPath())
	content := fmt.Sprintf(daoINITxTEMPLATE, importForClass, class.getClassName())

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

func (thisServer *server) getColumnsValuesAndMask(model IBusinessObjectModel, includeID bool, space string) (columns, placeholders, mask, values string) {

	colNb := 0

	for _, prop := range model.base().getPersistedProperties() {
		if prop.GetName() != BoFieldID || includeID {
			if columns != "" {
				columns += ", " + newline
				placeholders += ", " + newline
				mask += ", " + newline
				values += ", " + newline
			}
			if relationship, ok := prop.(*Relationship); ok {
				columns += space + relationship.getColumnName()
				colNb++
				placeholders += space + fmt.Sprintf("$%d", colNb)
				mask += core.IfThenElse(relationship.isSecret(), "true", "false")
				values += fmt.Sprintf("%s.%s.GetID()", core.PascalToCamel(string(model.base().name)), relationship.GetName())
				if relationship.IsPolymorphic() {
					columns += ", " + newline + space + relationship.getColumnNameForTargetClass()
					colNb++
					placeholders += ", " + newline + space + fmt.Sprintf("$%d", colNb)
					mask += ", " + newline + core.IfThenElse(relationship.isSecret(), "true", "false")
					values += ", " + newline + fmt.Sprintf("%[1]s.%[2]s.GetClassName(%[1]s.%[2]s)", core.PascalToCamel(string(model.base().name)), relationship.GetName())
				}
			} else {
				columns += space + prop.getColumnName()
				colNb++
				placeholders += space + fmt.Sprintf("$%d", colNb)
				mask += core.IfThenElse(prop.isSecret(), "true", "false")
				values += fmt.Sprintf("%s.%s", core.PascalToCamel(string(model.base().name)), prop.GetName())
			}
		}
	}

	values += "," // we need an ending commar here

	return
}

// ------------------------------------------------------------------------------------------------
//  DAO methods generation
// ------------------------------------------------------------------------------------------------

func (thisServer *server) generateExecInsertQuery(model IBusinessObjectModel) string {

	columns, placeholders, mask, values := thisServer.getColumnsValuesAndMask(model, false, "			")

	return fmt.Sprintf(`// ExecInsertQuery implements [goald.IBusinessObjectDAO].
func (thisDAO *%[1]sDAO) ExecInsertQuery(bObj goald.IBusinessObject) (int64, error) {
	// casting the business object to its actual type
	%[2]s := bObj.(*%[3]s.%[1]s)

	// actual insert query for a postgres DB
	insertQuery := `+bq+`
		INSERT INTO %[4]s (
%[5]s
		)
		VALUES (
%[6]s 
		)
		RETURNING id`+bq+`

	// the mask applied to the query arguments, to hide some of them from the logs (like the password)
	m := []bool{%[7]s}

	println(%[2]s)

	// executing the insert query and retrieving the new ID
	var newID int64
	err := thisDAO.QueryRow(m, insertQuery,
%[8]s
	).Scan(&newID)
	if err != nil {
		return 0, goald.ErrorC(err, "Error while performing insert query for %[1]s")
	}

	return newID, nil
}`,
		model.base().name, // 1
		core.PascalToCamel(string(model.base().name)), // 2
		getClass(model).getPackage(),                  // 3
		model.getTableName(true),                      // 4
		columns,                                       // 5
		placeholders,                                  // 6
		mask,                                          // 7
		values,                                        // 8
	)
}
