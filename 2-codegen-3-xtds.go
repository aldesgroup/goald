package goald

import (
	"fmt"
	"path"
	"strings"

	core "github.com/aldesgroup/corego"
)

const utilsFileTEMPLATE = `// Generated file, do not edit!
package $$package$$

import (
	$$otherimports$$
)

// getting the name of the model for a $$Upper$$, without using reflection
func (bo *$$Upper$$) GetModelName() utils.ModelName {
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

// removing any cycles from the business object, without using reflection
func (bo *$$Upper$$) RemoveCycles() {
$$removecycles$$}
`

const utilsFILExSUFFIX = "--xtd.go"

func (thisServer *server) generateAllObjectXTDs(srcdir, currentPath string, regen bool) (codeChanged bool) {
	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// going through the resources found withing the current directory
	// we got the BO & model registries, but we still need to browse the filesystem since we're updating it with files
	for _, entry := range core.EnsureReadDir(readingPath) {
		if entry.IsDir() {
			// not going into the vendor
			if entry.Name() != "vendor" && entry.Name() != ".git" {
				// found another directory, let's dive deeper!
				codeChanged = thisServer.generateAllObjectXTDs(srcdir, path.Join(currentPath, entry.Name()), regen) || codeChanged
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), boFILExSUFFIX) {
				// getting the model - which should exist at this stage!
				model := thisServer.getModelFromFile(srcdir, currentPath, entry.Name())

				// this file lives directly alongside the BO's own source file, within the BO's own package
				utilsFilepath := path.Join(srcdir, model.getSrcPath(), strings.Replace(entry.Name(), boFILExSUFFIX, utilsFILExSUFFIX, 1))

				// no xtd file for interfaces
				if !model.isInterface() {

					// generating the xtd file, if not existing yet, or too old
					if regen || !core.FileExists(utilsFilepath) || core.EnsureModTime(utilsFilepath).Before(model.getLastBOMod()) {
						generateObjectXtdForModel(model, utilsFilepath)
						codeChanged = true
					}
				}
			}
		}
	}

	return
}

// generating the xtd file for a given model, at the given path
func generateObjectXtdForModel(model IBusinessObjectModel, filepath string) {
	// need for some imports - goald & utils are always needed
	importsMap := map[string]bool{
		"github.com/aldesgroup/goald":                true,
		"github.com/aldesgroup/goald/features/utils": true,
	}

	// building the get/set cases, the relationship cases, and the validity checks
	getCases, setCases, importUtils := buildUtilsValueCases(model, importsMap)
	setRelCases, addRelCases, clearRelCases, getMultiRelCases := buildUtilsRelationshipCases(model, importsMap)
	modelChecks := buildUtilsModelChecks(model, importsMap)
	removeCyclesStatements := buildUtilsRemoveCyclesStatements(model)

	if importUtils {
		importsMap["github.com/aldesgroup/corego"] = true
	}

	// this file lives within the BO's own package: it must never import that same package
	delete(importsMap, path.Join(getCurrentSourceModule(), model.getSrcPath()))

	// starting the content
	content := strings.ReplaceAll(utilsFileTEMPLATE, "$$package$$", model.getPackage())
	content = strings.ReplaceAll(content, "$$Upper$$", string(model.getName()))
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

	removeCyclesBody := ""
	if len(removeCyclesStatements) > 0 {
		removeCyclesBody = strings.Join(removeCyclesStatements, newline) + newline
	}
	content = strings.ReplaceAll(content, "$$removecycles$$", removeCyclesBody)

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
func buildUtilsValueCases(model IBusinessObjectModel, importsMap map[string]bool) (getCases []string, setCases []string, importUtils bool) {
	// browsing the entity's properties to fill the get / set cases in the 2 switch
	for _, field := range core.GetSortedValues(model.getFields()) {
		// adding to the context, and the model file content
		if propertyType := field.getPropertyType(); propertyType != propertyTypeUNKNOWN && propertyType != propertyTypeRELATIONSHIPxMONOM {
			// not handling multiple properties for now
			if fieldName := field.GetName(); !field.IsMultiple() && fieldName != boFieldPreID {
				// is the field type a type alias, or a built-in type? - stripping this package's own
				// qualification, since we're generating code that lives directly within that package
				fieldTypeAlias := stripSelfPackage(getNonBuiltInFieldType(model.getType(), fieldName, importsMap), model.getPackage())

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
func buildUtilsRelationshipCases(model IBusinessObjectModel, importsMap map[string]bool) (setRelCases, addRelCases, clearRelCases, getMultiRelCases []string) {

	// browsing the entity's relationships to fill the set / add cases for the relationship setters
	for _, relationship := range core.GetSortedValues(model.getRelationships()) {
		relName := relationship.GetName()

		// the Go type to assert the incoming value against, e.g. "domain.IContact" or "*domain.Employee" -
		// stripping this package's own qualification, since we're generating code living within that package
		targetType := stripSelfPackage(getRelationshipFieldType(model.getType(), relName, importsMap), model.getPackage())

		relCase := fmt.Sprintf("\tcase \"%s\":", relName)
		relCase += newline + fmt.Sprintf("\t\ttargetValue, ok := value.(%s)", targetType)
		relCase += newline + "\t\tif !ok {"
		relCase += newline + fmt.Sprintf("\t\t\treturn goald.Error(\"Expected a value of type '%s' for '%s.%s', got %%T\", value)", targetType, model.getName(), relName)
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

// building the statements for RemoveCycles(): for every multi-valued relationship whose targets hold a
// single-valued backref pointing back to us (e.g. a PurchaseOrder's Items pointing back via OrderItem.PurchaseOrder),
// we nil out that backref on each target, so that marshaling this BO to JSON doesn't loop forever; polymorphic
// relationships are skipped, since the backref field then lives on an interface, not on a concrete type
func buildUtilsRemoveCyclesStatements(model IBusinessObjectModel) []string {
	statements := []string{}

	for _, relationship := range core.GetSortedValues(model.getRelationships()) {
		if relationship.IsMultiple() && !relationship.IsPolymorphic() && relationship.backRef != nil && !relationship.backRef.IsMultiple() {
			relName := relationship.GetName()
			backRefName := relationship.backRef.GetName()

			statements = append(statements, fmt.Sprintf("\tfor _, target := range bo.%s {", relName))
			statements = append(statements, fmt.Sprintf("\t\ttarget.%s = nil", backRefName))
			statements = append(statements, "\t}")
		}
	}

	return statements
}

// building the checks for IsModelValid(), simply reusing - unchanged - the check builders from the
// checks (--chk.go) generator: they already produce code referring to a "bo" variable, which is
// exactly the name we're using for this method's receiver, so no adaptation is needed there
func buildUtilsModelChecks(model IBusinessObjectModel, importsMap map[string]bool) []string {
	checks := []string{}

	// checking that every required relationship is properly set on the BO
	for _, relationship := range core.GetSortedValues(model.getRelationships()) {
		checks = append(checks, buildRequiredRelationshipChecks(relationship)...)

		// core.InSlice is used for polymorphic relationships' target model check
		if relationship.IsRequiredInDb() && relationship.IsPolymorphic() {
			importsMap["github.com/aldesgroup/corego"] = true
		}
	}

	// going through the entity's fields once, letting each per-field check builder chime in;
	// the mandatory-input check always comes first, before the other, more specific checks
	for _, field := range core.GetSortedValues(model.getFields()) {
		if check := buildMandatoryInputChecks(field); check != "" {
			checks = append(checks, check)
		}
		if check := buildFloatFormatChecks(field); check != "" {
			checks = append(checks, check)
		}
		if check := buildStringSizeChecks(field); check != "" {
			checks = append(checks, check)
		}
		if check := buildIntRangeChecks(field); check != "" {
			checks = append(checks, check)
		}
		if check := buildFloatRangeChecks(field); check != "" {
			checks = append(checks, check)
		}
	}

	return checks
}
