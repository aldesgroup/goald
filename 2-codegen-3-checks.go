// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the CHK (checks) files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"strings"

	core "github.com/aldesgroup/corego"
)

const checkFileTEMPLATE = `// Generated file, do not edit!
package $$package$$

import (
	$$otherimports$$
)

// checking a business object's general validity
func (thisClass *$$Upper$$Class) IsModelValid(bObj goald.IBusinessObject) error {
$$checks$$
	return nil
}
`

const checkFILExSUFFIX = "--chk.go"

func (thisServer *server) generateAllObjectChecks(srcdir, currentPath string, regen bool) (codeChanged bool) {
	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// going through the resources found withing the current directory
	// we got the BO & class registries, but we still need to browse the filesystem since we're updating it with files
	for _, entry := range core.EnsureReadDir(readingPath) {
		if entry.IsDir() {
			// not going into the vendor
			if entry.Name() != "vendor" && entry.Name() != ".git" {
				// found another directory, let's dive deeper!
				codeChanged = thisServer.generateAllObjectChecks(srcdir, path.Join(currentPath, entry.Name()), regen) || codeChanged
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

				class := classForName(classCore.class, true)

				// the corresponding Check file, if it exists
				checkFilepath := path.Join(srcdir, class.getSrcPath(), sourceCLASSxDIR,
					strings.Replace(entry.Name(), sourceFILExSUFFIX, checkFILExSUFFIX, 1))

				// no check file for interfaces
				if !class.isInterface() {

					// generating the Check file, if not existing yet, or too old
					if regen || !core.FileExists(checkFilepath) || core.EnsureModTime(checkFilepath).Before(class.getLastBOMod()) {
						generateObjectChecksForBO(class, checkFilepath)
						codeChanged = true
					}
				}
			}
		}
	}

	return
}

func generateObjectChecksForBO(class IClass, filepath string) {
	// the corresponding class
	className := class.getClassName()
	model := modelForName(className)

	// checking the BO code makes use of its class
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
	content := strings.ReplaceAll(checkFileTEMPLATE, "$$package$$", sourceCLASSxDIR)
	content = strings.ReplaceAll(content, "$$Upper$$", string(class.getClassName()))

	// need for some imports - just goald for now; the BO's own package is only added if actually needed below
	var importsMap = map[string]bool{
		"github.com/aldesgroup/goald": true,
	}

	// checking that every required relationship is properly set on the BO
	checks := []string{}
	for _, relationship := range core.GetSortedValues(model.base().relationships) {
		checks = append(checks, buildRequiredRelationshipCheck(relationship, className)...)

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

	checksBody := ""
	if len(checks) > 0 {
		// we only need to cast to the concrete BO type, and import its package, if we actually have checks to run
		importsMap[classPkg] = true

		checksBody = fmt.Sprintf("\tbo := bObj.(*%s.%s)\n\n", shortPkg, className) +
			strings.Join(checks, newline) + newline
	}
	content = strings.ReplaceAll(content, "$$checks$$", checksBody)

	imports := ""
	if len(importsMap) > 0 {
		imports = "\"" + strings.Join(core.GetSortedKeys(importsMap), "\""+newline+"\t"+"\"") + "\""
	}
	content = strings.Replace(content, "$$otherimports$$", imports, 1)

	// write out the file
	core.WriteToFile(content, filepath)
}

// building the checks ensuring that a required relationship (SetRequiredInDb) is properly set on
// the BO: its target must be non-nil, reference an already-persisted BO (ID > 0), and - if this
// is a polymorphic relationship - the target's concrete class must be one of the allowed ones
// (getTargetNames()); returns nil if this relationship isn't required
func buildRequiredRelationshipCheck(relationship *Relationship, className className) []string {
	if !relationship.IsRequiredInDb() {
		return nil
	}

	relName := relationship.GetName()

	checks := []string{
		fmt.Sprintf(
			"\tif bo.%s == nil {\n\t\treturn goald.Error(\"'%s' is required on '%s'\")\n\t}",
			relName, relName, className),
		fmt.Sprintf(
			"\tif bo.%s.GetID() <= 0 {\n\t\treturn goald.Error(\"'%s' must reference an existing, persisted business object\")\n\t}",
			relName, relName),
	}

	if relationship.IsPolymorphic() {
		targetNames := relationship.getTargetClassNames()
		allowed := make([]string, len(targetNames))
		for i, targetName := range targetNames {
			allowed[i] = fmt.Sprintf("%q", string(targetName))
		}

		checks = append(checks, fmt.Sprintf(
			"\tif !core.InSlice([]string{%s}, string(bo.%s.GetClassName(bo.%s))) {\n\t\treturn goald.Error(\"Invalid target class for '%s'\")\n\t}",
			strings.Join(allowed, ", "), relName, relName, relName))
	}

	return checks
}

// building the check ensuring that a mandatory input field (io:"i*") actually has a non-zero
// (from a Go standpoint) value; returns "" if there's nothing to check
func buildMandatoryInputCheck(field IField) string {
	if !field.isMandatoryInput() {
		return ""
	}

	zeroValue := zeroValueLiteral(field.getPropertyType())
	if zeroValue == "" {
		return ""
	}

	fieldName := field.GetName()

	return fmt.Sprintf(
		"\tif bo.%s == %s {\n\t\treturn goald.Error(\"'%s' is mandatory and must have a non-zero value\")\n\t}",
		fieldName, zeroValue, fieldName)
}

// returns the Go zero-value literal to compare a field's value against, given its property type,
// or "" if that property type isn't handled here
func zeroValueLiteral(propType propertyType) string {
	switch propType {
	case propertyTypeBOOL:
		return "false"
	case propertyTypeSTRING:
		return "\"\""
	case propertyTypeINT, propertyTypeBIGINT, propertyTypeREAL, propertyTypeDOUBLE, propertyTypeENUM:
		return "0"
	case propertyTypeDATE:
		return "nil"
	default:
		return ""
	}
}

// building the check ensuring that a float field's value complies with its declared format
// (total digits & decimals), if it has one; returns "" if there's nothing to check
func buildFloatFormatCheck(field IField) string {
	propertyType := field.getPropertyType()
	if propertyType != propertyTypeREAL && propertyType != propertyTypeDOUBLE {
		return ""
	}

	floatFmt, hasFormat := field.(iFloatFormat)
	if !hasFormat || floatFmt.GetTotalDigits() <= 0 {
		return ""
	}

	fieldName := field.GetName()
	valueExpr := fmt.Sprintf("bo.%s", fieldName)
	if propertyType == propertyTypeREAL {
		valueExpr = fmt.Sprintf("float64(%s)", valueExpr)
	}

	return fmt.Sprintf(
		"\tif err := goald.CheckFloatFormat(%s, %d, %d); err != nil {\n\t\treturn goald.ErrorC(err, \"Invalid value for '%s'\")\n\t}",
		valueExpr, floatFmt.GetTotalDigits(), floatFmt.GetDecimals(), fieldName)
}

// implemented by *RealField and *DoubleField (via the embedded floatField), letting us read the
// declared number format of a float field, whatever its exact (32 or 64 bits) concrete type
type iFloatFormat interface {
	GetTotalDigits() int
	GetDecimals() int
}

// building the check ensuring that a string field's value complies with its declared size
// constraints (SetSize), if it has one; returns "" if there's nothing to check
func buildStringSizeCheck(field IField) string {
	if field.getPropertyType() != propertyTypeSTRING {
		return ""
	}

	strField, ok := field.(*StringField)
	if !ok || (strField.GetSize() <= 0 && strField.GetAtLeast() <= 0) {
		return ""
	}

	fieldName := field.GetName()

	return fmt.Sprintf(
		"\tif err := goald.CheckStringSize(bo.%s, %d, %d); err != nil {\n\t\treturn goald.ErrorC(err, \"Invalid value for '%s'\")\n\t}",
		fieldName, strField.GetSize(), strField.GetAtLeast(), fieldName)
}

// building the check ensuring that an int (or bigint) field's value complies with its declared
// Min()/Max() bounds, if any; returns "" if there's nothing to check
func buildIntRangeCheck(field IField) string {
	var min, max int64
	var minSet, maxSet bool

	switch field.getPropertyType() {
	case propertyTypeINT:
		intField, ok := field.(*IntField)
		if !ok {
			return ""
		}
		min, max, minSet, maxSet = int64(intField.min), int64(intField.max), intField.minSet, intField.maxSet

	case propertyTypeBIGINT:
		bigIntField, ok := field.(*BigIntField)
		if !ok {
			return ""
		}
		min, max, minSet, maxSet = bigIntField.min, bigIntField.max, bigIntField.minSet, bigIntField.maxSet

	default:
		return ""
	}

	if !minSet && !maxSet {
		return ""
	}

	fieldName := field.GetName()

	return fmt.Sprintf(
		"\tif err := goald.CheckIntRange(int64(bo.%s), %d, %t, %d, %t); err != nil {\n\t\treturn goald.ErrorC(err, \"Invalid value for '%s'\")\n\t}",
		fieldName, min, minSet, max, maxSet, fieldName)
}

// building the check ensuring that a real (or double) field's value complies with its declared
// Min()/Max() bounds, if any; returns "" if there's nothing to check
func buildFloatRangeCheck(field IField) string {
	var min, max float64
	var minSet, maxSet bool

	switch field.getPropertyType() {
	case propertyTypeREAL:
		realField, ok := field.(*RealField)
		if !ok {
			return ""
		}
		min, max, minSet, maxSet = float64(realField.min), float64(realField.max), realField.minSet, realField.maxSet

	case propertyTypeDOUBLE:
		doubleField, ok := field.(*DoubleField)
		if !ok {
			return ""
		}
		min, max, minSet, maxSet = doubleField.min, doubleField.max, doubleField.minSet, doubleField.maxSet

	default:
		return ""
	}

	if !minSet && !maxSet {
		return ""
	}

	fieldName := field.GetName()

	return fmt.Sprintf(
		"\tif err := goald.CheckFloatRange(float64(bo.%s), %v, %t, %v, %t); err != nil {\n\t\treturn goald.ErrorC(err, \"Invalid value for '%s'\")\n\t}",
		fieldName, min, minSet, max, maxSet, fieldName)
}
