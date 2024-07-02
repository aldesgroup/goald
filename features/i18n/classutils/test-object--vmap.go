// Generated file, do not edit!
package classutils

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
	"github.com/aldesgroup/goald/features/other/nested"
	"github.com/aldesgroup/goald/features/utils"
)

// getting a property's value as a string, without using reflection
func (thisTestObjectClassUtils *TestObjectClassUtils) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "BoolProp":
		return utils.BoolToString(bo.(*i18n.TestObject).BoolProp)
	case "BoolPropCustom":
		return utils.BoolToString(bool(bo.(*i18n.TestObject).BoolPropCustom))
	case "DateProp":
		return utils.DateToString(bo.(*i18n.TestObject).DateProp)
	case "EnumProp":
		return utils.IntToString(bo.(*i18n.TestObject).EnumProp.Val())
	case "ID":
		return string(bo.(*i18n.TestObject).ID)
	case "Int64Prop":
		return utils.Int64ToString(bo.(*i18n.TestObject).Int64Prop)
	case "Int64PropCustom":
		return utils.Int64ToString(int64(bo.(*i18n.TestObject).Int64PropCustom))
	case "IntProp":
		return utils.IntToString(bo.(*i18n.TestObject).IntProp)
	case "IntPropCustom":
		return utils.IntToString(int(bo.(*i18n.TestObject).IntPropCustom))
	case "OtherEnumProp":
		return utils.IntToString(bo.(*i18n.TestObject).OtherEnumProp.Val())
	case "Real32Prop":
		return utils.Float32ToString(bo.(*i18n.TestObject).Real32Prop)
	case "Real32PropCustom":
		return utils.Float32ToString(float32(bo.(*i18n.TestObject).Real32PropCustom))
	case "Real64Prop":
		return utils.Float64ToString(bo.(*i18n.TestObject).Real64Prop)
	case "Real64PropCustom":
		return utils.Float64ToString(float64(bo.(*i18n.TestObject).Real64PropCustom))
	case "StringProp":
		return bo.(*i18n.TestObject).StringProp
	case "StringPropCustom":
		return string(bo.(*i18n.TestObject).StringPropCustom)
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTestObjectClassUtils *TestObjectClassUtils) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "BoolProp":
		bo.(*i18n.TestObject).BoolProp = utils.StringToBool(valueAsString, "(*i18n.TestObject).BoolProp")
	case "BoolPropCustom":
		bo.(*i18n.TestObject).BoolPropCustom = i18n.BoolPropCustomType(utils.StringToBool(valueAsString, "(*i18n.TestObject).BoolPropCustom"))
	case "DateProp":
		bo.(*i18n.TestObject).DateProp = utils.StringToDate(valueAsString, "(*i18n.TestObject).DateProp")
	case "EnumProp":
		bo.(*i18n.TestObject).EnumProp = i18n.Language(utils.StringToInt(valueAsString, "(*i18n.TestObject).EnumProp"))
		utils.PanicIff(bo.(*i18n.TestObject).EnumProp.String() == "", "Could not set '(*i18n.TestObject).EnumProp' to %s since it's not a listed value", valueAsString)
	case "ID":
		bo.(*i18n.TestObject).ID = goald.BObjID(valueAsString)
	case "Int64Prop":
		bo.(*i18n.TestObject).Int64Prop = utils.StringToInt64(valueAsString, "(*i18n.TestObject).Int64Prop")
	case "Int64PropCustom":
		bo.(*i18n.TestObject).Int64PropCustom = i18n.Int64PropCustomType(utils.StringToInt64(valueAsString, "(*i18n.TestObject).Int64PropCustom"))
	case "IntProp":
		bo.(*i18n.TestObject).IntProp = utils.StringToInt(valueAsString, "(*i18n.TestObject).IntProp")
	case "IntPropCustom":
		bo.(*i18n.TestObject).IntPropCustom = i18n.IntPropCustomType(utils.StringToInt(valueAsString, "(*i18n.TestObject).IntPropCustom"))
	case "OtherEnumProp":
		bo.(*i18n.TestObject).OtherEnumProp = nested.Origin(utils.StringToInt(valueAsString, "(*i18n.TestObject).OtherEnumProp"))
		utils.PanicIff(bo.(*i18n.TestObject).OtherEnumProp.String() == "", "Could not set '(*i18n.TestObject).OtherEnumProp' to %s since it's not a listed value", valueAsString)
	case "Real32Prop":
		bo.(*i18n.TestObject).Real32Prop = utils.StringToFloat32(valueAsString, "(*i18n.TestObject).Real32Prop")
	case "Real32PropCustom":
		bo.(*i18n.TestObject).Real32PropCustom = i18n.Real32PropCustomType(utils.StringToFloat32(valueAsString, "(*i18n.TestObject).Real32PropCustom"))
	case "Real64Prop":
		bo.(*i18n.TestObject).Real64Prop = utils.StringToFloat64(valueAsString, "(*i18n.TestObject).Real64Prop")
	case "Real64PropCustom":
		bo.(*i18n.TestObject).Real64PropCustom = i18n.Real64PropCustomType(utils.StringToFloat64(valueAsString, "(*i18n.TestObject).Real64PropCustom"))
	case "StringProp":
		bo.(*i18n.TestObject).StringProp = valueAsString
	case "StringPropCustom":
		bo.(*i18n.TestObject).StringPropCustom = i18n.StringPropCustomType(valueAsString)
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
