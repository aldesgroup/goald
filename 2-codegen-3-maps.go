// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the VMAP (value mapper) files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

const vmapFileTEMPLATE = `// Generated file, do not edit!
package $$package$$

import (
	$$otherimports$$
)

// getting a property's value as a string, without using reflection
func (thisClass *$$Upper$$Class) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
$$getcases$$
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *$$Upper$$Class) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
$$setcases$$
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
`

const valueMapperFILExSUFFIX = "--map.go"

func (thisServer *server) generateAllObjectValueMappers(srcdir, currentPath string, regen bool) (codeChanged bool) {
	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// going through the resources found withing the current directory
	// we got the BO & class registries, but we still need to browse the filesystem since we're updating it with files
	for _, entry := range core.EnsureReadDir(readingPath) {
		if entry.IsDir() {
			// not going into the vendor
			if entry.Name() != "vendor" && entry.Name() != ".git" {
				// found another directory, let's dive deeper!
				codeChanged = thisServer.generateAllObjectValueMappers(srcdir, path.Join(currentPath, entry.Name()), regen) || codeChanged
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry within this file, then the registred entry in the code
				classCore := getClassFromFile(srcdir, currentPath, entry.Name())
				class := classRegistry.items[classCore.class]

				if class == nil {
					core.PanicMsg("It looks like class '%s' has never been imported and thus not initialized and registered. "+
						"\nMake sure its module is imported in the main package: "+
						"\nimport _ \"%s/_include/%s\"",
						classCore.getClassName(), getCurrentModule(), currentPath)
				}

				// the corresponding Value Mapper file, if it exist
				vmapFilepath := path.Join(srcdir, class.getSrcPath(), sourceCLASSxDIR,
					strings.Replace(entry.Name(), sourceFILExSUFFIX, valueMapperFILExSUFFIX, 1))

				// no value mapper for interfaces
				if !class.isInterface() {

					// generating the Value Mapper file, if not existing yet, or too old
					if regen || !core.FileExists(vmapFilepath) || core.EnsureModTime(vmapFilepath).Before(class.getLastBOMod()) {
						generateObjectValueMappersForBO(class, vmapFilepath)
						codeChanged = true
					}
				}
			}
		}
	}

	return
}

func generateObjectValueMappersForBO(class IClass, filepath string) {
	// the corresponding class
	className := class.getClassName()
	boSpecs := specsForName(className)

	// checking the BO code makes use of its class
	// TODO - auto-add this code block to the BO code + the import
	if boSpecs == nil {
		core.PanicMsg("It looks like class '%s' has never been imported and thus not initialized and registered. \n"+
			"Add this - and complete as necessary - to your business object definition code: \n\n"+
			"import (class \"%s/_include/_specs\") \n"+
			"func init() { \n"+
			"	specs.%s().SetNotPersisted() \n"+
			"}",
			className, getCurrentModule(), className)
	}

	// the corresponding package
	classPkg := path.Join(getCurrentModule(), class.getSrcPath())
	shortPkg := path.Base(classPkg)

	// starting the content
	content := strings.ReplaceAll(vmapFileTEMPLATE, "$$package$$", sourceCLASSxDIR)
	content = strings.ReplaceAll(content, "$$Upper$$", string(class.getClassName()))

	getCases := []string{}
	setCases := []string{}

	// need for some imports
	var importsMap = map[string]bool{
		"github.com/aldesgroup/goald": true,
		classPkg:                      true,
	}
	var importUtils bool

	// getting the type of business object
	bObjectType := utils.TypeOf(class.NewObject(), true)

	// browsing the entity's properties to fill the get / set cases in the 2 switch
	for _, field := range core.GetSortedValues(boSpecs.base().fields) {
		// adding to the context, and the class file content
		if typeFamily := field.getTypeFamily(); typeFamily != utils.TypeFamilyUNKNOWN && typeFamily != utils.TypeFamilyRELATIONSHIPxMONOM {
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
					getCase += newline + fmt.Sprintf("\t\treturn core.BoolToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToBool(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilySTRING:
					getBit, setBit, end := getBits(fieldTypeAlias, "string")
					getCase += newline + fmt.Sprintf("\t\treturn %sbo.%s%s", getBit, fieldID, end)
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %svalueAsString%s", fieldID, setBit, end)

				case utils.TypeFamilyINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int")
					getCase += newline + fmt.Sprintf("\t\treturn core.IntToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToInt(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyBIGINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int64")
					getCase += newline + fmt.Sprintf("\t\treturn core.Int64ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToInt64(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyREAL:
					getBit, setBit, end := getBits(fieldTypeAlias, "float32")
					getCase += newline + fmt.Sprintf("\t\treturn core.Float32ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToFloat32(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyDOUBLE:
					getBit, setBit, end := getBits(fieldTypeAlias, "float64")
					getCase += newline + fmt.Sprintf("\t\treturn core.Float64ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToFloat64(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case utils.TypeFamilyDATE:
					getCase += newline + fmt.Sprintf("\t\treturn core.DateToString(bo.%s)", fieldID)
					setCase += newline + fmt.Sprintf("\t\tbo.%s = core.StringToDate(valueAsString, \"%s\")", fieldID, fieldID)

				case utils.TypeFamilyENUM:
					getCase += newline + fmt.Sprintf("\t\treturn core.IntToString(bo.%s.Val())", fieldID)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %s(core.StringToInt(valueAsString, \"%s\"))", fieldID, fieldTypeAlias, fieldID)

					setCase += newline + fmt.Sprintf("\t\tcore.PanicMsgIf(bo.%s.String() == \"\", \"Could not set '%s' to %%s since it's not a listed value\", valueAsString)",
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
		importsMap["github.com/aldesgroup/corego"] = true
	}

	imports := ""
	if len(importsMap) > 0 {
		imports = "\"" + strings.Join(core.GetSortedKeys(importsMap), "\""+newline+"\t"+"\"") + "\""
	}
	content = strings.Replace(content, "$$otherimports$$", imports, 1)

	// write out the file
	core.WriteToFile(content, filepath)
}

func getBits(fieldTypeAlias, getBit string) (string, string, string) {
	if fieldTypeAlias != "" {
		return getBit + "(", fieldTypeAlias + "(", ")"
	}

	return "", "", ""
}

func getNonBuiltInFieldType(bOjbType utils.GoaldType, fieldName string, toBeImported map[string]bool) string {
	fieldType := bOjbType.FieldByName(fieldName).Type()
	fieldPkg := fieldType.PkgPath()

	// this is a built-in field type
	if fieldPkg == "" {
		return ""
	}

	// the field type comes from another package, that we have to import
	if toBeImported != nil {
		toBeImported[fieldPkg] = true
	}

	return fieldType.String() // e.g.: thatpackage.MyEnumType
}
