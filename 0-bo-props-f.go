package goald

import (
	"strconv"
	"strings"
	"unicode/utf8"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Fields (simple properties) of business object models
// ------------------------------------------------------------------------------------------------

type IField interface {
	IBusinessObjectProperty
	isBuiltIn() bool
	getDefaultValue() string
}

type iNumericField interface {
	IField
	isMinSet() bool
	isMaxSet() bool
}

// base implementation
type field struct {
	businessObjectProperty
	defaultStringValue string
}

type numericField struct {
	field
	minSet bool
	maxSet bool
}

func newField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool, propType propertyType) field {
	return field{
		businessObjectProperty: businessObjectProperty{
			owner:       owner,
			declaringBO: declaringBO,
			name:        name,
			propType:    propType,
			multiple:    multiple,
		},
	}
}

func (f *field) SetNotPersisted() *field {
	f.notPersisted = true
	return f
}

func (f *field) isBuiltIn() bool {
	return false
}

func (f *field) SetDefaultValue(val string) *field {
	f.defaultStringValue = val
	return f
}

func (f *field) getDefaultValue() string {
	return f.defaultStringValue
}

func (f *numericField) isMinSet() bool {
	return f.minSet
}

func (f *numericField) isMaxSet() bool {
	return f.maxSet
}

type BoolField struct {
	field
}

type StringField struct {
	field
	size    int
	atLeast int
}

func (sf *StringField) GetSize() int {
	return sf.size
}

func (sf *StringField) GetAtLeast() int {
	return sf.atLeast
}

func (sf *StringField) SetSize(size int, atLeast ...int) *field {
	sf.size = size
	if len(atLeast) > 0 {
		sf.atLeast = atLeast[0]
	}
	return &sf.field
}

// CheckStringSize checks that the given string value's length (in characters) is between
// atLeast and size (inclusive). size <= 0 means there's no maximum length constraint, and
// atLeast <= 0 means there's no minimum length constraint.
func CheckStringSize(value string, size int, atLeast int) error {
	length := utf8.RuneCountInString(value)

	if size > 0 && length > size {
		return Error("Value '%s' is too long (%d character(s), max %d)", value, length, size)
	}

	if atLeast > 0 && length < atLeast {
		return Error("Value '%s' is too short (%d character(s), min %d)", value, length, atLeast)
	}

	return nil
}

type IntField struct {
	numericField
	min        int
	max        int
	equalLenTo IBusinessObjectProperty // if set, this field's value must equal len(bo.<property>)
}

func (f *IntField) Min(min int) *IntField {
	f.min = min
	f.minSet = true
	return f
}

func (f *IntField) Max(max int) *IntField {
	f.max = max
	f.maxSet = true
	return f
}

// SetEqualLen declares that this int field's value must always equal the number of targets of the
// given, multi-valued property (a relationship, e.g. o.NbItems().SetEqualLen(o.Items()), or a
// multi-valued field), ensuring NbItems always equals len(Items). This gets enforced as a generated
// check in IsModelValid(). Panics if the given property isn't itself multi-valued (IsMultiple()).
func (f *IntField) SetEqualLen(property IBusinessObjectProperty) *IntField {
	core.PanicMsgIf(!property.IsMultiple(), "SetEqualLen: '%s' is not a multi-valued property", property.GetName())
	f.equalLenTo = property
	return f
}

type BigIntField struct {
	numericField
	min int64
	max int64
}

func (f *BigIntField) Min(min int64) *BigIntField {
	f.min = min
	f.minSet = true
	return f
}

func (f *BigIntField) Max(max int64) *BigIntField {
	f.max = max
	f.maxSet = true
	return f
}

// CheckIntRange checks that value is within [min, max], considering only the bounds that are
// actually set (minSet / maxSet).
func CheckIntRange(value int64, min int64, minSet bool, max int64, maxSet bool) error {
	if minSet && value < min {
		return Error("Value %d is below the minimum of %d", value, min)
	}

	if maxSet && value > max {
		return Error("Value %d is above the maximum of %d", value, max)
	}

	return nil
}

type floatField struct {
	numericField
	totalDigits int // number of total digits for the numbers, e.g. 8 in 91876.063
	decimals    int // number of digits after the decimal points, e.g. 3 in 91876.063
}

func (f *floatField) SetFormat(totalDigits int, decimals int) {
	f.totalDigits = totalDigits
	f.decimals = decimals
}

func (f *floatField) GetTotalDigits() int {
	return f.totalDigits
}

func (f *floatField) GetDecimals() int {
	return f.decimals
}

// CheckFloatFormat checks that the given value fits within the given number format, i.e. has no
// more than totalDigits significant digits in total, of which no more than decimals after the
// decimal point - mirroring a SQL DECIMAL(totalDigits, decimals) column.
// A totalDigits <= 0 means there's no format constraint to check.
// bitSize should be 32 if the value originally comes from a float32 field, and 64 for a float64 one -
// this matters because a float32 value, once widened to float64, exposes binary rounding noise as a
// long decimal tail (e.g. 8772.652 as float32 becomes 8772.65234375 as float64): passing the correct
// bitSize tells strconv to compute the shortest decimal string that round-trips at that precision,
// rather than at float64's much higher precision.
func CheckFloatFormat(value float64, totalDigits int, decimals int, bitSize int) error {
	if totalDigits <= 0 {
		return nil
	}

	// reasoning on the shortest decimal string that round-trips back to this exact value at the
	// given bit size: cheaper than, and immune to the rounding noise of, redoing the math with Pow/Round
	formatted := strings.TrimPrefix(strconv.FormatFloat(value, 'f', -1, bitSize), "-")
	intPart, decPart, hasDecimals := strings.Cut(formatted, ".")

	if hasDecimals && len(decPart) > decimals {
		return Error("Value %v has more than %d decimal(s)", value, decimals)
	}

	if maxIntDigits := totalDigits - decimals; intPart != "0" && len(intPart) > maxIntDigits {
		return Error("Value %v has more than %d integer digit(s)", value, maxIntDigits)
	}

	return nil
}

type RealField struct {
	floatField
	min float32
	max float32
}

func (f *RealField) Min(min float32) *RealField {
	f.min = min
	f.minSet = true
	return f
}

func (f *RealField) Max(max float32) *RealField {
	f.max = max
	f.maxSet = true
	return f
}

type DoubleField struct {
	floatField
	min float64
	max float64
}

func (f *DoubleField) Min(min float64) *DoubleField {
	f.min = min
	f.minSet = true
	return f
}

func (f *DoubleField) Max(max float64) *DoubleField {
	f.max = max
	f.maxSet = true
	return f
}

// CheckFloatRange checks that value is within [min, max], considering only the bounds that are
// actually set (minSet / maxSet).
func CheckFloatRange(value float64, min float64, minSet bool, max float64, maxSet bool) error {
	if minSet && value < min {
		return Error("Value %v is below the minimum of %v", value, min)
	}

	if maxSet && value > max {
		return Error("Value %v is above the maximum of %v", value, max)
	}

	return nil
}

type DateField struct {
	field
}

type EnumField struct {
	field
	enumName   string
	onlyValues []IEnum
}

func (f *EnumField) Only(values ...IEnum) *EnumField {
	f.onlyValues = values
	return f
}

func AddBoolField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool) *BoolField {
	return owner.addField(&BoolField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeBOOL),
	}).(*BoolField)
}

func AddStringField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool) *StringField {
	return owner.addField(&StringField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeSTRING),
	}).(*StringField)
}

func AddIntField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool) *IntField {
	return owner.addField(&IntField{numericField: numericField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeINT),
	}}).(*IntField)
}

func AddBigIntField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool) *BigIntField {
	return owner.addField(&BigIntField{numericField: numericField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeBIGINT),
	}}).(*BigIntField)
}

func AddRealField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool) *RealField {
	return owner.addField(&RealField{floatField: floatField{numericField: numericField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeREAL),
	}}}).(*RealField)
}

func AddDoubleField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool) *DoubleField {
	return owner.addField(&DoubleField{floatField: floatField{numericField: numericField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeDOUBLE),
	}}}).(*DoubleField)
}

func AddDateField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool) *DateField {
	return owner.addField(&DateField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeDATE),
	}).(*DateField)
}

func AddEnumField(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool, enumName string) *EnumField {
	return owner.addField(&EnumField{
		field:    newField(owner, declaringBO, name, multiple, propertyTypeENUM),
		enumName: enumName,
	}).(*EnumField)
}
