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

const vmapFileTEMPLATE = `// Generated file, do not edit!
package $$package$$

import (
	$$otherimports$$
)

// getting a property's value as a string, without using reflection
func (thisUtils *$$Upper$$ClassUtils) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
$$getcases$$
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisUtils *$$Upper$$ClassUtils) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
$$setcases$$
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
`

const valueMapperFILExSUFFIX = "--vmap.go"

func (thisServer *server) generateObjectValueMappers(srcdir, currentPath string, regen bool) {
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
				thisServer.generateObjectValueMappers(srcdir, path.Join(currentPath, entry.Name()), regen)
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry within this file, then the registred entry in the code
				clsuCore := getClassUtilsFromFile(srcdir, currentPath, entry.Name())
				classUtils := classUtilsRegistry.content[clsuCore.class]

				// the corresponding Value Mapper file, if it exist
				vmapFilename := path.Join(sourceCLASSxUTILSxDIR, strings.Replace(entry.Name(), sourceFILExSUFFIX, valueMapperFILExSUFFIX, 1))

				// generating the Value Mapper file, if not existing yet, or too old
				if regen || !utils.FileExists(vmapFilename) || utils.EnsureModTime(vmapFilename).Before(classUtils.getLastBOMod()) {
					generateObjectValueMappersForBO(srcdir, classUtils, vmapFilename)
				}
			}
		}
	}
}

func generateObjectValueMappersForBO(srcdir string, classUtils IClassUtils, filename string) {
	// the corresponding class
	className := classUtils.getClass()
	boClass := classForName(className)

	// the corresponding package
	classPkg := path.Join(getCurrentModule(), classUtils.getSrcPath())
	shortPkg := path.Base(classPkg)

	// starting the content
	content := strings.ReplaceAll(vmapFileTEMPLATE, "$$package$$", sourceCLASSxUTILSxDIR)
	content = strings.ReplaceAll(content, "$$Upper$$", string(classUtils.getClass()))

	getCases := []string{}
	setCases := []string{}

	// need for some imports
	var importsMap = map[string]bool{
		"github.com/aldesgroup/goald": true,
		classPkg:                      true,
	}
	var importUtils bool

	// getting the type of business object
	bObjectType := utils.TypeOf(classUtils.NewObject(), true)

	// browsing the entity's properties to fill the get / set cases in the 2 switch
	// for fieldNum := 1; fieldNum < classUtils.bObjType.NumField(); fieldNum++ {
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
				fieldID := fmt.Sprintf("(*%s.%s).%s", shortPkg, className, fieldName)

				switch typeFamily {
				case utils.TypeFamilyBOOL:
					getBit, setBit, end := getBits(fieldTypeAlias, "bool")
					getCase += newline + fmt.Sprintf("\t\treturn utils.BoolToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %sutils.StringToBool(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilySTRING:
					getBit, setBit, end := getBits(fieldTypeAlias, "string")
					getCase += newline + fmt.Sprintf("\t\treturn %sbo.%s%s", getBit, fieldID, end)
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %svalueAsString%s", fieldID, setBit, end)

				case utils.TypeFamilyINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int")
					getCase += newline + fmt.Sprintf("\t\treturn utils.IntToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %sutils.StringToInt(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyBIGINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int64")
					getCase += newline + fmt.Sprintf("\t\treturn utils.Int64ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %sutils.StringToInt64(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyREAL:
					getBit, setBit, end := getBits(fieldTypeAlias, "float32")
					getCase += newline + fmt.Sprintf("\t\treturn utils.Float32ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %sutils.StringToFloat32(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyDOUBLE:
					getBit, setBit, end := getBits(fieldTypeAlias, "float64")
					getCase += newline + fmt.Sprintf("\t\treturn utils.Float64ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %sutils.StringToFloat64(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyDATE:
					getCase += newline + fmt.Sprintf("\t\treturn utils.DateToString(bo.%s)", fieldID)
					setCase += newline + fmt.Sprintf("\t\tbo.%s = utils.StringToDate(valueAsString, \"%s\")", fieldID, fieldID)

				case utils.TypeFamilyENUM:
					getCase += newline + fmt.Sprintf("\t\treturn utils.IntToString(bo.%s.Val())", fieldID)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %s(utils.StringToInt(valueAsString, \"%s\"))", fieldID, fieldTypeAlias, fieldID)

					setCase += newline + fmt.Sprintf("\t\tutils.PanicIff(bo.%s.String() == \"\", \"Could not set '%s' to %%s since it's not a listed value\", valueAsString)",
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
	utils.WriteToFile(content, srcdir, classUtils.getSrcPath(), filename)
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

	// the field type comes from another package, that we have to import
	toBeImported[fieldPkg] = true
	return fieldType.String() // e.g.: thatpackage.MyEnumType
}
