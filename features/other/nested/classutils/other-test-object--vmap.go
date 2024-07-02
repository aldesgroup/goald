// Generated file, do not edit!
package classutils

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/other/nested"
	"github.com/aldesgroup/goald/features/utils"
)

// getting a property's value as a string, without using reflection
func (thisOtherTestObjectClassUtils *OtherTestObjectClassUtils) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "BoolProp":
		return utils.BoolToString(bo.(*nested.OtherTestObject).BoolProp)
	case "BoolPropCustom":
		return utils.BoolToString(bool(bo.(*nested.OtherTestObject).BoolPropCustom))
	case "DateProp":
		return utils.DateToString(bo.(*nested.OtherTestObject).DateProp)
	case "ID":
		return string(bo.(*nested.OtherTestObject).ID)
	case "Int64Prop":
		return utils.Int64ToString(bo.(*nested.OtherTestObject).Int64Prop)
	case "Int64PropCustom":
		return utils.Int64ToString(int64(bo.(*nested.OtherTestObject).Int64PropCustom))
	case "IntProp":
		return utils.IntToString(bo.(*nested.OtherTestObject).IntProp)
	case "IntPropCustom":
		return utils.IntToString(int(bo.(*nested.OtherTestObject).IntPropCustom))
	case "OtherEnumProp":
		return utils.IntToString(bo.(*nested.OtherTestObject).OtherEnumProp.Val())
	case "Real32Prop":
		return utils.Float32ToString(bo.(*nested.OtherTestObject).Real32Prop)
	case "Real32PropCustom":
		return utils.Float32ToString(float32(bo.(*nested.OtherTestObject).Real32PropCustom))
	case "Real64Prop":
		return utils.Float64ToString(bo.(*nested.OtherTestObject).Real64Prop)
	case "Real64PropCustom":
		return utils.Float64ToString(float64(bo.(*nested.OtherTestObject).Real64PropCustom))
	case "StringProp":
		return bo.(*nested.OtherTestObject).StringProp
	case "StringPropCustom":
		return string(bo.(*nested.OtherTestObject).StringPropCustom)
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisOtherTestObjectClassUtils *OtherTestObjectClassUtils) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "BoolProp":
		bo.(*nested.OtherTestObject).BoolProp = utils.StringToBool(valueAsString, "(*nested.OtherTestObject).BoolProp")
	case "BoolPropCustom":
		bo.(*nested.OtherTestObject).BoolPropCustom = nested.BoolPropCustomType(utils.StringToBool(valueAsString, "(*nested.OtherTestObject).BoolPropCustom"))
	case "DateProp":
		bo.(*nested.OtherTestObject).DateProp = utils.StringToDate(valueAsString, "(*nested.OtherTestObject).DateProp")
	case "ID":
		bo.(*nested.OtherTestObject).ID = goald.BObjID(valueAsString)
	case "Int64Prop":
		bo.(*nested.OtherTestObject).Int64Prop = utils.StringToInt64(valueAsString, "(*nested.OtherTestObject).Int64Prop")
	case "Int64PropCustom":
		bo.(*nested.OtherTestObject).Int64PropCustom = nested.Int64PropCustomType(utils.StringToInt64(valueAsString, "(*nested.OtherTestObject).Int64PropCustom"))
	case "IntProp":
		bo.(*nested.OtherTestObject).IntProp = utils.StringToInt(valueAsString, "(*nested.OtherTestObject).IntProp")
	case "IntPropCustom":
		bo.(*nested.OtherTestObject).IntPropCustom = nested.IntPropCustomType(utils.StringToInt(valueAsString, "(*nested.OtherTestObject).IntPropCustom"))
	case "OtherEnumProp":
		bo.(*nested.OtherTestObject).OtherEnumProp = nested.Origin(utils.StringToInt(valueAsString, "(*nested.OtherTestObject).OtherEnumProp"))
		utils.PanicIff(bo.(*nested.OtherTestObject).OtherEnumProp.String() == "", "Could not set '(*nested.OtherTestObject).OtherEnumProp' to %s since it's not a listed value", valueAsString)
	case "Real32Prop":
		bo.(*nested.OtherTestObject).Real32Prop = utils.StringToFloat32(valueAsString, "(*nested.OtherTestObject).Real32Prop")
	case "Real32PropCustom":
		bo.(*nested.OtherTestObject).Real32PropCustom = nested.Real32PropCustomType(utils.StringToFloat32(valueAsString, "(*nested.OtherTestObject).Real32PropCustom"))
	case "Real64Prop":
		bo.(*nested.OtherTestObject).Real64Prop = utils.StringToFloat64(valueAsString, "(*nested.OtherTestObject).Real64Prop")
	case "Real64PropCustom":
		bo.(*nested.OtherTestObject).Real64PropCustom = nested.Real64PropCustomType(utils.StringToFloat64(valueAsString, "(*nested.OtherTestObject).Real64PropCustom"))
	case "StringProp":
		bo.(*nested.OtherTestObject).StringProp = valueAsString
	case "StringPropCustom":
		bo.(*nested.OtherTestObject).StringPropCustom = nested.StringPropCustomType(valueAsString)
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
