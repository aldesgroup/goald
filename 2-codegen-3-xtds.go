// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the XTD (extension) files - these are meant to eventually
// replace both the CHK (checks) and VMAP (value mapper) files: the difference is that the generated
// methods here are attached directly to the Business Object types themselves (instead of to their
// "*Class" companion objects), so there's no more need to pass the BO as an argument, nor to cast it.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/reflection"
	"github.com/aldesgroup/goald/features/utils"
)

const utilsFileTEMPLATE = `// Generated file, do not edit!
package $$package$$

import (
	$$otherimports$$
)

// getting the name of the class for a $$Upper$$, without using reflection
func (bo *$$Upper$$) ClassName() utils.ClassName {
	return "$$Upper$$"
}

// getting a property's value as a string, without using reflection
func (bo *$$Upper$$) GetValueAsString(propertyName string) string {
	switch propertyName {
$$getcases$$
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *$$Upper$$) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
$$setcases$$
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *$$Upper$$) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
$$setrelcases$$
	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *$$Upper$$) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
$$addrelcases$$
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *$$Upper$$) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {
$$clearrelcases$$
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *$$Upper$$) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {
$$getmultiplerelcases$$
	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// checking a business object's general validity, without using reflection
func (bo *$$Upper$$) IsModelValid() error {
$$modelchecks$$
	return nil
}
`

const utilsFILExSUFFIX = "--xtd.go"

func (thisServer *server) generateAllObjectUtils(srcdir, currentPath string, regen bool) (codeChanged bool) {
	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// going through the resources found withing the current directory
	// we got the BO & class registries, but we still need to browse the filesystem since we're updating it with files
	for _, entry := range core.EnsureReadDir(readingPath) {
		if entry.IsDir() {
			// not going into the vendor
			if entry.Name() != "vendor" && entry.Name() != ".git" {
				// found another directory, let's dive deeper!
				codeChanged = thisServer.generateAllObjectUtils(srcdir, path.Join(currentPath, entry.Name()), regen) || codeChanged
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry within this file, then the registred entry in the code
				baseClass := getClassFromFile(srcdir, currentPath, entry.Name())

				if baseClass == nil {
					core.PanicMsg("It looks like there's no BusinessObject-derived struct in '%s/%s/%s'",
						srcdir, currentPath, entry.Name())
				}

				class := classFor(baseClass.class, true)

				// this file lives directly alongside the BO's own source file, within the BO's own package -
				// there's no more need for a separate "class" subdirectory, since there's no more casting/indirection
				utilsFilepath := path.Join(srcdir, class.getSrcPath(),
					strings.Replace(entry.Name(), sourceFILExSUFFIX, utilsFILExSUFFIX, 1))

				// no xtd file for interfaces
				if !class.isInterface() {

					// generating the xtd file, if not existing yet, or too old
					if regen || !core.FileExists(utilsFilepath) || core.EnsureModTime(utilsFilepath).Before(class.getLastBOMod()) {
						generateObjectUtilsForBO(class, utilsFilepath)
						codeChanged = true
					}
				}
			}
		}
	}

	return
}

// generating a BO's "xtd" file, gathering what used to be split between the value mapper (--map.go)
// and the checks (--chk.go) generators - reusing as much as possible of their case/check-building logic
func generateObjectUtilsForBO(class IClass, filepath string) {
	// the corresponding class
	className := class.getClassName()
	model := class.getModel()

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

	// the package this file will actually live in: the BO's own package, since these methods are
	// now attached directly to the BO type - so this package must never end up importing itself
	classPkg := path.Join(getCurrentModule(), class.getSrcPath())
	shortPkg := path.Base(classPkg)

	// need for some imports - goald & utils (for the Class() method) are always needed
	importsMap := map[string]bool{
		"github.com/aldesgroup/goald":                true,
		"github.com/aldesgroup/goald/features/utils": true,
	}

	// getting the type of business object
	bObjectType := reflection.TypeOf(class.NewObject(), true)

	// building the get/set cases, the relationship cases, and the validity checks
	getCases, setCases, importUtils := buildUtilsValueCases(model, bObjectType, shortPkg, importsMap)
	setRelCases, addRelCases, clearRelCases, getMultiRelCases := buildUtilsRelationshipCases(model, bObjectType, className, shortPkg, importsMap)
	modelChecks := buildUtilsModelChecks(model, className, importsMap)

	if importUtils {
		importsMap["github.com/aldesgroup/corego"] = true
	}

	// this file lives within the BO's own package: it must never import that same package
	delete(importsMap, classPkg)

	// starting the content
	content := strings.ReplaceAll(utilsFileTEMPLATE, "$$package$$", shortPkg)
	content = strings.ReplaceAll(content, "$$Upper$$", string(className))
	content = strings.ReplaceAll(content, "$$getcases$$", strings.Join(getCases, newline))
	content = strings.ReplaceAll(content, "$$setcases$$", strings.Join(setCases, newline))
	content = strings.ReplaceAll(content, "$$setrelcases$$", strings.Join(setRelCases, newline))
	content = strings.ReplaceAll(content, "$$addrelcases$$", strings.Join(addRelCases, newline))
	content = strings.ReplaceAll(content, "$$clearrelcases$$", strings.Join(clearRelCases, newline))
	content = strings.ReplaceAll(content, "$$getmultiplerelcases$$", strings.Join(getMultiRelCases, newline))

	checksBody := ""
	if len(modelChecks) > 0 {
		checksBody = strings.Join(modelChecks, newline) + newline
	}
	content = strings.ReplaceAll(content, "$$modelchecks$$", checksBody)

	imports := ""
	if len(importsMap) > 0 {
		imports = "\"" + strings.Join(core.GetSortedKeys(importsMap), "\""+newline+"\t"+"\"") + "\""
	}
	content = strings.Replace(content, "$$otherimports$$", imports, 1)

	// write out the file
	core.WriteToFile(content, filepath)
}

// building the get/set cases for GetValueAsString() / SetValueAsString(), reusing the same
// per-property-type logic as the value mapper generator, but accessing fields directly on the
// receiver (named "bo" in the generated code), since no casting is needed anymore
func buildUtilsValueCases(model IBusinessObjectModel, bObjectType reflection.GoaldType, shortPkg string, importsMap map[string]bool) (getCases []string, setCases []string, importUtils bool) {
	// browsing the entity's properties to fill the get / set cases in the 2 switch
	for _, field := range core.GetSortedValues(model.base().fields) {
		// adding to the context, and the class file content
		if propertyType := field.getPropertyType(); propertyType != propertyTypeUNKNOWN && propertyType != propertyTypeRELATIONSHIPxMONOM {
			// not handling multiple properties for now
			if fieldName := field.GetName(); !field.IsMultiple() && fieldName != boFieldPreID {
				// is the field type a type alias, or a built-in type? - stripping this package's own
				// qualification, since we're generating code that lives directly within that package
				fieldTypeAlias := stripSelfPackage(getNonBuiltInFieldType(bObjectType, fieldName, importsMap), shortPkg)

				// case init
				getCase := fmt.Sprintf("\tcase \"%s\":", fieldName)
				setCase := getCase

				switch propertyType {
				case propertyTypeBOOL:
					getBit, setBit, end := getBits(fieldTypeAlias, "bool")
					getCase += newline + fmt.Sprintf("\t\treturn core.BoolToString(%sbo.%s%s)", getBit, fieldName, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToBool(valueAsString, \"%s\")%s", fieldName, setBit, fieldName, end)

				case propertyTypeSTRING:
					getBit, setBit, end := getBits(fieldTypeAlias, "string")
					getCase += newline + fmt.Sprintf("\t\treturn %sbo.%s%s", getBit, fieldName, end)
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %svalueAsString%s", fieldName, setBit, end)

				case propertyTypeINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int")
					getCase += newline + fmt.Sprintf("\t\treturn core.IntToString(%sbo.%s%s)", getBit, fieldName, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToInt(valueAsString, \"%s\")%s", fieldName, setBit, fieldName, end)

				case propertyTypeBIGINT:
					getBit, setBit, end := getBits(fieldTypeAlias, "int64")
					getCase += newline + fmt.Sprintf("\t\treturn core.Int64ToString(%sbo.%s%s)", getBit, fieldName, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToInt64(valueAsString, \"%s\")%s", fieldName, setBit, fieldName, end)

				case propertyTypeREAL:
					getBit, setBit, end := getBits(fieldTypeAlias, "float32")
					getCase += newline + fmt.Sprintf("\t\treturn core.Float32ToString(%sbo.%s%s)", getBit, fieldName, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToFloat32(valueAsString, \"%s\")%s", fieldName, setBit, fieldName, end)

				case propertyTypeDOUBLE:
					getBit, setBit, end := getBits(fieldTypeAlias, "float64")
					getCase += newline + fmt.Sprintf("\t\treturn core.Float64ToString(%sbo.%s%s)", getBit, fieldName, end)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %score.StringToFloat64(valueAsString, \"%s\")%s", fieldName, setBit, fieldName, end)

				case propertyTypeDATE:
					getCase += newline + fmt.Sprintf("\t\treturn core.DateToString(bo.%s)", fieldName)
					setCase += newline + fmt.Sprintf("\t\tbo.%s = core.StringToDate(valueAsString, \"%s\")", fieldName, fieldName)

				case propertyTypeENUM:
					getCase += newline + fmt.Sprintf("\t\treturn core.IntToString(bo.%s.Val())", fieldName)
					importUtils = true
					setCase += newline + fmt.Sprintf("\t\tbo.%s = %s(core.StringToInt(valueAsString, \"%s\"))", fieldName, fieldTypeAlias, fieldName)
					setCase += newline + fmt.Sprintf("\t\tcore.PanicMsgIf(bo.%s.String() == \"\", \"Could not set '%s' to %%s since it's not a listed value\", valueAsString)",
						fieldName, fieldName)
				}

				// appending the case
				getCases = append(getCases, getCase)
				setCases = append(setCases, setCase)
			}
		}
	}

	return getCases, setCases, importUtils
}

// building the cases for the relationship setters/getters, reusing the same logic as the value
// mapper generator, but without any casting, since these methods are now attached to the BO itself
func buildUtilsRelationshipCases(model IBusinessObjectModel, bObjectType reflection.GoaldType, clsName utils.ClassName, shortPkg string,
	importsMap map[string]bool) (setRelCases, addRelCases, clearRelCases, getMultiRelCases []string) {

	// browsing the entity's relationships to fill the set / add cases for the relationship setters
	for _, relationship := range core.GetSortedValues(model.base().relationships) {
		relName := relationship.GetName()

		// the Go type to assert the incoming value against, e.g. "domain.IContact" or "*domain.Employee" -
		// stripping this package's own qualification, since we're generating code living within that package
		targetType := stripSelfPackage(getRelationshipFieldType(bObjectType, relName, importsMap), shortPkg)

		relCase := fmt.Sprintf("\tcase \"%s\":", relName)
		relCase += newline + fmt.Sprintf("\t\ttargetValue, ok := value.(%s)", targetType)
		relCase += newline + "\t\tif !ok {"
		relCase += newline + fmt.Sprintf("\t\t\treturn goald.Error(\"Expected a value of type '%s' for '%s.%s', got %%T\", value)", targetType, clsName, relName)
		relCase += newline + "\t\t}"

		if relationship.IsMultiple() {
			relCase += newline + fmt.Sprintf("\t\tbo.%s = append(bo.%s, targetValue)", relName, relName)
			relCase += newline + "\t\treturn nil"
			addRelCases = append(addRelCases, relCase)

			clearCase := fmt.Sprintf("\tcase \"%s\":", relName)
			clearCase += newline + fmt.Sprintf("\t\tbo.%s = []%s{}", relName, targetType)
			clearCase += newline + "\t\treturn nil"
			clearRelCases = append(clearRelCases, clearCase)

			getMultiRelCases = append(getMultiRelCases, buildGetMultiRelCase(relationship, relName))
		} else {
			relCase += newline + fmt.Sprintf("\t\tbo.%s = targetValue", relName)
			relCase += newline + "\t\treturn nil"
			setRelCases = append(setRelCases, relCase)
		}
	}

	return setRelCases, addRelCases, clearRelCases, getMultiRelCases
}

// building the corresponding GetMultipleRelationshipValue case; if there's a backref relationship
// on the target side that's single-valued (i.e. each target uniquely points back to us), we
// also make sure it's (re)set, since it might not have been loaded/set that way already - no
// casting is needed here anymore, since "bo" is already of the right, concrete type
func buildGetMultiRelCase(relationship *Relationship, relName string) string {
	resultVar := core.PascalToCamel(relName)
	getMultiCase := fmt.Sprintf("\tcase \"%s\":", relName)
	getMultiCase += newline + fmt.Sprintf("\t\t%s := make([]goald.IBusinessObject, len(bo.%s))", resultVar, relName)
	getMultiCase += newline + fmt.Sprintf("\t\tfor i, target := range bo.%s {", relName)
	if relationship.backRef != nil && !relationship.backRef.IsMultiple() {
		getMultiCase += newline + fmt.Sprintf("\t\t\ttarget.%s = bo // ensuring the unique backref is set", relationship.backRef.GetName())
	}
	getMultiCase += newline + fmt.Sprintf("\t\t\t%s[i] = target", resultVar)
	getMultiCase += newline + "\t\t}"
	getMultiCase += newline + fmt.Sprintf("\t\treturn %s, nil", resultVar)
	return getMultiCase
}

// building the checks for IsModelValid(), simply reusing - unchanged - the check builders from the
// checks (--chk.go) generator: they already produce code referring to a "bo" variable, which is
// exactly the name we're using for this method's receiver, so no adaptation is needed there
func buildUtilsModelChecks(model IBusinessObjectModel, clsName utils.ClassName, importsMap map[string]bool) []string {
	checks := []string{}

	// checking that every required relationship is properly set on the BO
	for _, relationship := range core.GetSortedValues(model.base().relationships) {
		checks = append(checks, buildRequiredRelationshipCheck(relationship, clsName)...)

		// core.InSlice is used for polymorphic relationships' target class check
		if relationship.IsRequiredInDb() && relationship.IsPolymorphic() {
			importsMap["github.com/aldesgroup/corego"] = true
		}
	}

	// going through the entity's fields once, letting each per-field check builder chime in;
	// the mandatory-input check always comes first, before the other, more specific checks
	for _, field := range core.GetSortedValues(model.base().fields) {
		if check := buildMandatoryInputCheck(field); check != "" {
			checks = append(checks, check)
		}
		if check := buildFloatFormatCheck(field); check != "" {
			checks = append(checks, check)
		}
		if check := buildStringSizeCheck(field); check != "" {
			checks = append(checks, check)
		}
		if check := buildIntRangeCheck(field); check != "" {
			checks = append(checks, check)
		}
		if check := buildFloatRangeCheck(field); check != "" {
			checks = append(checks, check)
		}
	}

	return checks
}
