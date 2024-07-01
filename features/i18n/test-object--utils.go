// Generated file, do not edit!
package i18n

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/other/nested"
	"github.com/aldesgroup/goald/features/utils"
)

// getting a property's value as a string, without using reflection
func (thisTestObject *TestObject) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "BoolProp":
		return utils.BoolToString(thisTestObject.BoolProp)
	case "BoolPropCustom":
		return utils.BoolToString(bool(thisTestObject.BoolPropCustom))
	case "DateProp":
		return utils.DateToString(thisTestObject.DateProp)
	case "EnumProp":
		return utils.IntToString(thisTestObject.EnumProp.Val())
	case "ID":
		return string(thisTestObject.ID)
	case "Int64Prop":
		return utils.Int64ToString(thisTestObject.Int64Prop)
	case "Int64PropCustom":
		return utils.Int64ToString(int64(thisTestObject.Int64PropCustom))
	case "IntProp":
		return utils.IntToString(thisTestObject.IntProp)
	case "IntPropCustom":
		return utils.IntToString(int(thisTestObject.IntPropCustom))
	case "OtherEnumProp":
		return utils.IntToString(thisTestObject.OtherEnumProp.Val())
	case "Real32Prop":
		return utils.Float32ToString(thisTestObject.Real32Prop)
	case "Real32PropCustom":
		return utils.Float32ToString(float32(thisTestObject.Real32PropCustom))
	case "Real64Prop":
		return utils.Float64ToString(thisTestObject.Real64Prop)
	case "Real64PropCustom":
		return utils.Float64ToString(float64(thisTestObject.Real64PropCustom))
	case "StringProp":
		return thisTestObject.StringProp
	case "StringPropCustom":
		return string(thisTestObject.StringPropCustom)
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTestObject *TestObject) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "BoolProp":
		thisTestObject.BoolProp = utils.StringToBool(valueAsString, "TestObject.BoolProp")
	case "BoolPropCustom":
		thisTestObject.BoolPropCustom = BoolPropCustomType(utils.StringToBool(valueAsString, "TestObject.BoolPropCustom"))
	case "DateProp":
		thisTestObject.DateProp = utils.StringToDate(valueAsString, "TestObject.DateProp")
	case "EnumProp":
		thisTestObject.EnumProp = Language(utils.StringToInt(valueAsString, "TestObject.EnumProp"))
		utils.PanicIff(thisTestObject.EnumProp.String() == "", "Could not set 'TestObject.EnumProp' to %s since it's not a listed value", valueAsString)
	case "ID":
		thisTestObject.ID = goald.BObjID(valueAsString)
	case "Int64Prop":
		thisTestObject.Int64Prop = utils.StringToInt64(valueAsString, "TestObject.Int64Prop")
	case "Int64PropCustom":
		thisTestObject.Int64PropCustom = Int64PropCustomType(utils.StringToInt64(valueAsString, "TestObject.Int64PropCustom"))
	case "IntProp":
		thisTestObject.IntProp = utils.StringToInt(valueAsString, "TestObject.IntProp")
	case "IntPropCustom":
		thisTestObject.IntPropCustom = IntPropCustomType(utils.StringToInt(valueAsString, "TestObject.IntPropCustom"))
	case "OtherEnumProp":
		thisTestObject.OtherEnumProp = nested.Origin(utils.StringToInt(valueAsString, "TestObject.OtherEnumProp"))
		utils.PanicIff(thisTestObject.OtherEnumProp.String() == "", "Could not set 'TestObject.OtherEnumProp' to %s since it's not a listed value", valueAsString)
	case "Real32Prop":
		thisTestObject.Real32Prop = utils.StringToFloat32(valueAsString, "TestObject.Real32Prop")
	case "Real32PropCustom":
		thisTestObject.Real32PropCustom = Real32PropCustomType(utils.StringToFloat32(valueAsString, "TestObject.Real32PropCustom"))
	case "Real64Prop":
		thisTestObject.Real64Prop = utils.StringToFloat64(valueAsString, "TestObject.Real64Prop")
	case "Real64PropCustom":
		thisTestObject.Real64PropCustom = Real64PropCustomType(utils.StringToFloat64(valueAsString, "TestObject.Real64PropCustom"))
	case "StringProp":
		thisTestObject.StringProp = valueAsString
	case "StringPropCustom":
		thisTestObject.StringPropCustom = StringPropCustomType(valueAsString)
	}

	return goald.Error("Unknown property: %T.%s", thisTestObject, propertyName)
}
