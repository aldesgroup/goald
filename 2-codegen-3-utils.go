// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the class files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aldesgroup/goald/features/utils"
)

const utilsTEMPLATE = `// Generated file, do not edit!
package $$package$$

import (
	$$otherimports$$
)

// getting a property's value as a string, without using reflection
func (this$$Upper$$ *$$Upper$$) GetValueAsString(propertyName string) string {
	switch propertyName {
$$getcases$$
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (this$$Upper$$ *$$Upper$$) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
$$setcases$$
	}

	return goald.Error("Unknown property: %T.%s", this$$Upper$$, propertyName)
}
`

const utilsFILExSUFFIX = "--utils.go"

func (thisServer *server) generateObjectUtils(srcdir, currentPath string, regen bool) {
	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// reading the current directory
	dirEntries, errDir := os.ReadDir(readingPath)
	utils.PanicErrf(errDir, "could not read '%s'", readingPath)

	// going through the resources found withing the current directory
	// we got the BO & class registries, but we still need to browse the filesystem since we're updating it with files
	for _, entry := range dirEntries {
		if entry.IsDir() {
			// not going into the vendor
			if entry.Name() != "vendor" && entry.Name() != ".git" {
				// found another directory, let's dive deeper!
				thisServer.generateObjectUtils(srcdir, path.Join(currentPath, entry.Name()), regen)
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry within this file, then the registred entry in the code
				bObjEntryFile := getEntryFromFile(srcdir, currentPath, entry.Name())
				bObjEntry := boRegistry.content[bObjEntryFile.class]

				// the corresponding utils file, if it exist
				utilsFilename := strings.Replace(entry.Name(), sourceFILExSUFFIX, utilsFILExSUFFIX, 1)

				// generating the utils file, if not existing yet, or too old
				if regen || !utils.FileExists(utilsFilename) || utils.EnsureModTime(utilsFilename).Before(bObjEntry.lastMod) {
					generateObjectUtilsForEntry(srcdir, bObjEntry, utilsFilename)
				}
			}
		}
	}
}

func generateObjectUtilsForEntry(srcdir string, bObjectEntry *businessObjectEntry, filename string) {
	// the corresponding class
	className := bObjectEntry.class
	boClass := classForName(className)

	// starting the content
	content := strings.ReplaceAll(utilsTEMPLATE, "$$package$$", path.Base(bObjectEntry.srcPath))
	content = strings.ReplaceAll(content, "$$Upper$$", string(bObjectEntry.class))

	getCases := []string{}
	setCases := []string{}

	// need for some imports
	var importsMap = map[string]bool{"github.com/aldesgroup/goald": true}
	var importUtils bool

	// getting the type of business object
	bObjectType := utils.TypeOf(bObjectEntry.instanceFn(), true)

	// browsing the entity's properties to fill the get / set cases in the 2 switch
	// for fieldNum := 1; fieldNum < bObjEntry.bObjType.NumField(); fieldNum++ {
	for _, field := range utils.GetSortedValues(boClass.base().fields) {
		// adding to the context, and the class file content
		if typeFamily := field.getTypeFamily(); typeFamily != utils.TypeFamilyUNKNOWN && typeFamily != utils.TypeFamilyRELATIONSHIP {
			// not handling multiple properties for now
			if fieldName := field.getName(); !field.isMultiple() {
				// is the field type a type alias, or a built-in type?
				fieldTypeAlias := getNonBuiltInFieldType(bObjectType, fieldName, importsMap)

				// case init
				getCase := fmt.Sprintf("\tcase \"%s\":", fieldName)
				setCase := getCase

				// this is going to come up a lot
				fieldID := string(className) + "." + fieldName

				switch typeFamily {
				case utils.TypeFamilyBOOL:
					getBit, setBit, end := getBits(fieldTypeAlias, "bool")
					getCase += newline + fmt.Sprintf("\t\treturn utils.BoolToString(%sthis%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tthis%s = %sutils.StringToBool(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilySTRING:
					getBit, setBit, end := getBits(fieldTypeAlias, "string")
					getCase += newline + fmt.Sprintf("\t\treturn %sthis%s%s", getBit, fieldID, end)
					setCase += newline + fmt.Sprintf("\t\tthis%s = %svalueAsString%s", fieldID, setBit, end)

				case utils.TypeFamilyINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int")
					getCase += newline + fmt.Sprintf("\t\treturn utils.IntToString(%sthis%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tthis%s = %sutils.StringToInt(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyBIGINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int64")
					getCase += newline + fmt.Sprintf("\t\treturn utils.Int64ToString(%sthis%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tthis%s = %sutils.StringToInt64(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyREAL:
					getBit, setBit, end := getBits(fieldTypeAlias, "float32")
					getCase += newline + fmt.Sprintf("\t\treturn utils.Float32ToString(%sthis%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tthis%s = %sutils.StringToFloat32(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyDOUBLE:
					getBit, setBit, end := getBits(fieldTypeAlias, "float64")
					getCase += newline + fmt.Sprintf("\t\treturn utils.Float64ToString(%sthis%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tthis%s = %sutils.StringToFloat64(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyDATE:
					getCase += newline + fmt.Sprintf("\t\treturn utils.DateToString(this%s)", fieldID)
					setCase += newline + fmt.Sprintf("\t\tthis%s = utils.StringToDate(valueAsString, \"%s\")", fieldID, fieldID)

				case utils.TypeFamilyENUM:
					getCase += newline + fmt.Sprintf("\t\treturn utils.IntToString(this%s.Val())", fieldID)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tthis%s = %s(utils.StringToInt(valueAsString, \"%s\"))", fieldID, fieldTypeAlias, fieldID)

					setCase += newline + fmt.Sprintf("\t\tutils.PanicIff(this%s.String() == \"\", \"Could not set '%s' to %%s since it's not a listed value\", valueAsString)",
						fieldID, fieldID)
				}

				// appending the case
				getCases = append(getCases, getCase)
				setCases = append(setCases, setCase)
			}
		}
	}

	// handling the imports
	content = strings.ReplaceAll(content, "$$getcases$$", strings.Join(getCases, newline))
	content = strings.ReplaceAll(content, "$$setcases$$", strings.Join(setCases, newline))

	if importUtils {
		importsMap["github.com/aldesgroup/goald/features/utils"] = true
	}

	imports := ""
	if len(importsMap) > 0 {
		imports = "\"" + strings.Join(utils.GetSortedKeys(importsMap), "\""+newline+"\t"+"\"") + "\""
	}
	content = strings.Replace(content, "$$otherimports$$", imports, 1)

	// write out the file
	utils.WriteToFile(content, srcdir, bObjectEntry.srcPath, filename)
}

func getBits(fieldTypeAlias, getBit string) (string, string, string) {
	if fieldTypeAlias != "" {
		return getBit + "(", fieldTypeAlias + "(", ")"
	}

	return "", "", ""
}

func getNonBuiltInFieldType(bOjbType *utils.GoaldType, fieldName string, toBeImported map[string]bool) string {
	fieldType := bOjbType.FieldByName(fieldName).Type()
	fieldPkg := fieldType.PkgPath()

	// this is a built-in field type
	if fieldPkg == "" {
		return ""
	}

	// the type of the field is defined in the same package as the business object
	if fieldPkg == bOjbType.PkgPath() {
		return fieldType.Name() // e.g.: MyEnumType
	}

	// the field type comes from another package, that we have to import
	toBeImported[fieldPkg] = true
	return fieldType.String() // e.g.: thatpackage.MyEnumType
}
