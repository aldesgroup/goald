// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the VMAP (value mapper) files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/reflection"
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

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (thisClass *$$Upper$$Class) SetRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
$$setrelcases$$
	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (thisClass *$$Upper$$Class) AddRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
$$addrelcases$$
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (thisClass *$$Upper$$Class) ClearRelationshipValue(bo goald.IBusinessObject, relationshipName string) error {
	switch relationshipName {
$$clearrelcases$$
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
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

				if classCore == nil {
					core.PanicMsg("It looks like there's no BusinessObject-derived struct in '%s/%s/%s'",
						srcdir, currentPath, entry.Name())
				}

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
	model := modelForName(className)

	// checking the BO code makes use of its class
	// TODO - auto-add this code block into the BO code + the import
	if model == nil {
		core.PanicMsg("It looks like class '%s' has never been imported and thus not initialized and registered. \n"+
			"Add this - and complete as necessary - to your business object definition code: \n\n"+
			"import (class \"%s/_include/%s/model\") \n"+
			"func init() { \n"+
			"	model.%s().SetNotPersisted() \n"+
			"}",
			className, getCurrentModule(), class.getPackage(), className)
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
	bObjectType := reflection.TypeOf(class.NewObject(), true)

	// browsing the entity's properties to fill the get / set cases in the 2 switch
	for _, field := range core.GetSortedValues(model.base().fields) {
		// adding to the context, and the class file content
		if propertyType := field.getPropertyType(); propertyType != propertyTypeUNKNOWN && propertyType != propertyTypeRELATIONSHIPxMONOM {
			// not handling multiple properties for now
			if fieldName := field.GetName(); !field.IsMultiple() {
				// is the field type a type alias, or a built-in type?
				fieldTypeAlias := getNonBuiltInFieldType(bObjectType, fieldName, importsMap)

				// case init
				getCase := fmt.Sprintf("\tcase \"%s\":", fieldName)
				setCase := getCase

				// this is going to come up a lot
				fieldID := fmt.Sprintf("(*%s.%s).%s", shortPkg, className, fieldName)

				switch propertyType {
				case propertyTypeBOOL:
					getBit, setBit, end := getBits(fieldTypeAlias, "bool")
					getCase += newline + fmt.Sprintf("\t\treturn core.BoolToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToBool(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case propertyTypeSTRING:
					getBit, setBit, end := getBits(fieldTypeAlias, "string")
					getCase += newline + fmt.Sprintf("\t\treturn %sbo.%s%s", getBit, fieldID, end)
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %svalueAsString%s", fieldID, setBit, end)

				case propertyTypeINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int")
					getCase += newline + fmt.Sprintf("\t\treturn core.IntToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToInt(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case propertyTypeBIGINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int64")
					getCase += newline + fmt.Sprintf("\t\treturn core.Int64ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToInt64(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case propertyTypeREAL:
					getBit, setBit, end := getBits(fieldTypeAlias, "float32")
					getCase += newline + fmt.Sprintf("\t\treturn core.Float32ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToFloat32(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case propertyTypeDOUBLE:
					getBit, setBit, end := getBits(fieldTypeAlias, "float64")
					getCase += newline + fmt.Sprintf("\t\treturn core.Float64ToString(%sbo.%s%s)", getBit, fieldID, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToFloat64(valueAsString, \"%s\")%s", fieldID, setBit, fieldID, end)

				case propertyTypeDATE:
					getCase += newline + fmt.Sprintf("\t\treturn core.DateToString(bo.%s)", fieldID)
					setCase += newline + fmt.Sprintf("\t\tbo.%s = core.StringToDate(valueAsString, \"%s\")", fieldID, fieldID)

				case propertyTypeENUM:
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

	setRelCases := []string{}
	addRelCases := []string{}
	clearRelCases := []string{}

	// browsing the entity's relationships to fill the set / add cases for the relationship setters
	for _, relationship := range core.GetSortedValues(model.base().relationships) {
		relName := relationship.GetName()

		// the Go type to assert the incoming value against, e.g. "domain.IContact" or "*domain.Employee"
		targetType := getRelationshipFieldType(bObjectType, relName, importsMap)

		relCase := fmt.Sprintf("\tcase \"%s\":", relName)
		relCase += newline + fmt.Sprintf("\t\ttargetValue, ok := value.(%s)", targetType)
		relCase += newline + "\t\tif !ok {"
		relCase += newline + fmt.Sprintf("\t\t\treturn goald.Error(\"Expected a value of type '%s' for '%s.%s', got %%T\", value)", targetType, className, relName)
		relCase += newline + "\t\t}"

		fieldID := fmt.Sprintf("(*%s.%s).%s", shortPkg, className, relName)

		if relationship.IsMultiple() {
			relCase += newline + fmt.Sprintf("\t\tbo.%s = append(bo.%s, targetValue)", fieldID, fieldID)
			relCase += newline + "\t\treturn nil"
			addRelCases = append(addRelCases, relCase)

			clearCase := fmt.Sprintf("\tcase \"%s\":", relName)
			clearCase += newline + fmt.Sprintf("\t\tbo.%s = []%s{}", fieldID, targetType)
			clearCase += newline + "\t\treturn nil"
			clearRelCases = append(clearRelCases, clearCase)
		} else {
			relCase += newline + fmt.Sprintf("\t\tbo.%s = targetValue", fieldID)
			relCase += newline + "\t\treturn nil"
			setRelCases = append(setRelCases, relCase)
		}
	}

	// handling the imports
	content = strings.ReplaceAll(content, "$$getcases$$", strings.Join(getCases, newline))
	content = strings.ReplaceAll(content, "$$setcases$$", strings.Join(setCases, newline))
	content = strings.ReplaceAll(content, "$$setrelcases$$", strings.Join(setRelCases, newline))
	content = strings.ReplaceAll(content, "$$addrelcases$$", strings.Join(addRelCases, newline))
	content = strings.ReplaceAll(content, "$$clearrelcases$$", strings.Join(clearRelCases, newline))

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

func getNonBuiltInFieldType(bOjbType reflection.GoaldType, fieldName string, toBeImported map[string]bool) string {
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

// returns the Go type to use to type-assert a relationship's incoming value against, e.g.
// "domain.IContact" for a polymorphic relationship, or "*domain.Employee" for a monomorphic one -
// also registering the corresponding package for import, if needed
func getRelationshipFieldType(bOjbType reflection.GoaldType, fieldName string, toBeImported map[string]bool) string {
	fieldType := bOjbType.FieldByName(fieldName).Type()

	// for a multi-valued relationship, the Go field is a slice: we need its element type
	elemType := fieldType
	if fieldType.Kind() == reflection.KindSLICE {
		elemType = fieldType.Elem()
	}

	// the package to import is the one declaring the pointed-to type, whether we have a pointer or an interface
	pkgSource := elemType
	if elemType.Kind() == reflection.KindPTR {
		pkgSource = elemType.Elem()
	}

	if pkgPath := pkgSource.PkgPath(); pkgPath != "" && toBeImported != nil {
		toBeImported[pkgPath] = true
	}

	return elemType.String()
}
