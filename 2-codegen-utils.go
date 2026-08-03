package goald

import (
	"fmt"
	"strings"

	"github.com/aldesgroup/goald/features/reflection"
)

// building the checks ensuring that a required relationship (SetRequiredInDb) is properly set on
// the BO: for a single-valued relationship, its target must be non-nil and reference an
// already-persisted BO (ID > 0); for a multi-valued relationship, there must be at least one
// target, and each of them must reference an already-persisted BO (ID > 0); in both cases, if
// this is a polymorphic relationship, the target's (or each target's) concrete model must be one
// of the allowed ones (getTargetNames()); returns nil if this relationship isn't required
func buildRequiredRelationshipChecks(relationship *Relationship) []string {
	if !relationship.IsRequiredInDb() && !relationship.isMandatoryInput() {
		return nil
	}

	relName := relationship.GetName()

	if relationship.IsMultiple() {
		checks := []string{
			fmt.Sprintf(
				"\tif len(bo.%s) == 0 {\n\t\treturn goald.Error(\"'%s' is required on '%s'\")\n\t}",
				relName, relName, relationship.owner.getName()),
		}

		loopBody := []string{
			fmt.Sprintf("\t\tif target.GetID() <= 0 {\n\t\t\treturn goald.Error(\"'%s' must reference an existing, persisted business object\")\n\t\t}", relName),
		}

		if relationship.IsPolymorphic() {
			targetNames := relationship.getTargetModelNames()
			allowed := make([]string, len(targetNames))
			for i, targetName := range targetNames {
				allowed[i] = fmt.Sprintf("%q", string(targetName))
			}

			loopBody = append(loopBody, fmt.Sprintf(
				"\t\tif !core.InSlice([]string{%s}, string(target.GetModelName())) {\n\t\t\treturn goald.Error(\"Invalid target model for '%s'\")\n\t\t}",
				strings.Join(allowed, ", "), relName))
		}

		checks = append(checks, fmt.Sprintf("\tfor _, target := range bo.%s {\n%s\n\t}", relName, strings.Join(loopBody, "\n")))

		return checks
	}

	checks := []string{
		fmt.Sprintf(
			"\tif bo.%s == nil {\n\t\treturn goald.Error(\"'%s' is required on '%s'\")\n\t}",
			relName, relName, relationship.owner.getName()),
		fmt.Sprintf(
			"\tif bo.%s.GetID() <= 0 {\n\t\treturn goald.Error(\"'%s' must reference an existing, persisted business object\")\n\t}",
			relName, relName),
	}

	if relationship.IsPolymorphic() {
		targetNames := relationship.getTargetModelNames()
		allowed := make([]string, len(targetNames))
		for i, targetName := range targetNames {
			allowed[i] = fmt.Sprintf("%q", string(targetName))
		}

		checks = append(checks, fmt.Sprintf(
			"\tif !core.InSlice([]string{%s}, string(bo.%s.GetModelName())) {\n\t\treturn goald.Error(\"Invalid target model for '%s'\")\n\t}",
			strings.Join(allowed, ", "), relName, relName))
	}

	return checks
}

// building the check ensuring that a mandatory input field (io:"i*") actually has a non-zero
// (from a Go standpoint) value; returns "" if there's nothing to check
func buildMandatoryInputChecks(field IField) string {
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
func buildFloatFormatChecks(field IField) string {
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
	bitSize := 64
	if propertyType == propertyTypeREAL {
		valueExpr = fmt.Sprintf("float64(%s)", valueExpr)
		bitSize = 32
	}

	return fmt.Sprintf(
		"\tif err := goald.CheckFloatFormat(%s, %d, %d, %d); err != nil {\n\t\treturn goald.ErrorC(err, \"Invalid value for '%s'\")\n\t}",
		valueExpr, floatFmt.GetTotalDigits(), floatFmt.GetDecimals(), bitSize, fieldName)
}

// implemented by *RealField and *DoubleField (via the embedded floatField), letting us read the
// declared number format of a float field, whatever its exact (32 or 64 bits) concrete type
type iFloatFormat interface {
	GetTotalDigits() int
	GetDecimals() int
}

// building the check ensuring that a string field's value complies with its declared size
// constraints (SetSize), if it has one; returns "" if there's nothing to check
func buildStringSizeChecks(field IField) string {
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
func buildIntRangeChecks(field IField) string {
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

// building the check ensuring that an int field declared with SetEqualLen actually equals the
// number of targets of the multi-valued property (relationship or field) it's tied to (e.g.
// NbItems == len(Items)); returns "" if this field isn't an int field, or has no SetEqualLen
// property set
func buildEqualLenChecks(field IField) string {
	intField, ok := field.(*IntField)
	if !ok || intField.equalLenTo == nil {
		return ""
	}

	fieldName := field.GetName()
	propName := intField.equalLenTo.GetName()

	return fmt.Sprintf(
		"\tif bo.%s != len(bo.%s) {\n\t\treturn goald.Error(\"'%s' must be equal to the number of '%s'\")\n\t}",
		fieldName, propName, fieldName, propName)
}

// building the check ensuring that a real (or double) field's value complies with its declared
// Min()/Max() bounds, if any; returns "" if there's nothing to check
func buildFloatRangeChecks(field IField) string {
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

// building the check ensuring that an enum field's value is a legit one, i.e. listed in the map
// returned by its IEnum's Values() method (e.g. OrderStatus.Values()); this catches values that
// were set from some arbitrary, unlisted integer (e.g. via SetValueAsString bypassing its own
// guard, or a raw cast); returns "" if this field isn't an enum field
func buildEnumValueChecks(field IField) string {
	if field.getPropertyType() != propertyTypeENUM {
		return ""
	}

	fieldName := field.GetName()

	return fmt.Sprintf(
		"\tif _, isLegitValue := bo.%[1]s.Values()[bo.%[1]s.Val()]; !isLegitValue {\n\t\treturn goald.Error(\"Invalid value '%%d' for '%[2]s.%[1]s'\", bo.%[1]s.Val())\n\t}",
		fieldName, field.ownerModel().getName())
}

// removing this package's own qualification from a generated type expression (e.g. turning
// "accessmgt.Employee" into "Employee"), since the xtd file lives directly within that same
// package, and thus must never reference its own types through a package prefix
func stripSelfPackage(typeExpr, shortPkg string) string {
	if typeExpr == "" || shortPkg == "" {
		return typeExpr
	}

	return strings.ReplaceAll(typeExpr, shortPkg+".", "")
}

func getBits(fieldTypeAlias, getBit string) (string, string, string) {
	if fieldTypeAlias != "" {
		return getBit + "(", fieldTypeAlias + "(", ")"
	}

	return "", "", ""
}

func getNonBuiltInFieldType(bOjbType *reflection.GoaldType, fieldName string, toBeImported map[string]bool) string {
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
func getRelationshipFieldType(bOjbType *reflection.GoaldType, fieldName string, toBeImported map[string]bool) string {
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
