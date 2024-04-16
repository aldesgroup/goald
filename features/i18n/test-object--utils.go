// Generated file, do not edit!
package i18n

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/other/nested"
	"github.com/aldesgroup/goald/features/utils"
	"strconv"
	"time"
)

// getting a property's value as a string, without using reflection
func (thisTestObject *TestObject) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "BoolProp":
		return strconv.FormatBool(thisTestObject.BoolProp)
	case "Int64Prop":
		return strconv.FormatInt(thisTestObject.Int64Prop, 10)
	case "StringProp":
		return thisTestObject.StringProp
	case "IntProp":
		return strconv.Itoa(thisTestObject.IntProp)
	case "Real32Prop":
		return strconv.FormatFloat(float64(thisTestObject.Real32Prop), 'f', -1, 32)
	case "Real64Prop":
		return strconv.FormatFloat(thisTestObject.Real64Prop, 'f', -1, 64)
	case "EnumProp":
		return strconv.Itoa(thisTestObject.EnumProp.Val())
	case "DateProp":
		return thisTestObject.DateProp.Format(utils.RFC3339Milli)
	case "OtherEnumProp":
		return strconv.Itoa(thisTestObject.OtherEnumProp.Val())
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTestObject *TestObject) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "BoolProp":
		value, errConv := strconv.ParseBool(valueAsString)
		utils.PanicErrf(errConv, "Could not set 'TestObject.BoolProp' to %s", valueAsString)
		thisTestObject.BoolProp = value
	case "Int64Prop":
		value, errConv := strconv.ParseInt(valueAsString, 10, 64)
		utils.PanicErrf(errConv, "Could not set 'TestObject.Int64Prop' to %s", valueAsString)
		thisTestObject.Int64Prop = value
	case "StringProp":
		thisTestObject.StringProp = valueAsString
	case "IntProp":
		value, errConv := strconv.Atoi(valueAsString)
		utils.PanicErrf(errConv, "Could not set 'TestObject.IntProp' to %s", valueAsString)
		thisTestObject.IntProp = value
	case "Real32Prop":
		value, errConv := strconv.ParseFloat(valueAsString, 32)
		utils.PanicErrf(errConv, "Could not set 'TestObject.Real32Prop' to %s", valueAsString)
		thisTestObject.Real32Prop = float32(value)
	case "Real64Prop":
		value, errConv := strconv.ParseFloat(valueAsString, 64)
		utils.PanicErrf(errConv, "Could not set 'TestObject.Real64Prop' to %s", valueAsString)
		thisTestObject.Real64Prop = value
	case "EnumProp":
		intValue, errConv := strconv.Atoi(valueAsString)
		utils.PanicErrf(errConv, "Could not set 'TestObject.EnumProp' to %s", valueAsString)
		thisTestObject.EnumProp = (Language)(intValue)
		utils.PanicIff(thisTestObject.EnumProp.String() == "", "Could not set 'TestObject.EnumProp' to %s since it's not a listed value", valueAsString)
	case "DateProp":
		value, errConv := time.Parse(utils.RFC3339Milli, valueAsString)
		utils.PanicErrf(errConv, "Could not set 'TestObject.DateProp' to %s", valueAsString)
		thisTestObject.DateProp = &value
	case "OtherEnumProp":
		intValue, errConv := strconv.Atoi(valueAsString)
		utils.PanicErrf(errConv, "Could not set 'TestObject.OtherEnumProp' to %s", valueAsString)
		thisTestObject.OtherEnumProp = (nested.Origin)(intValue)
		utils.PanicIff(thisTestObject.OtherEnumProp.String() == "", "Could not set 'TestObject.OtherEnumProp' to %s since it's not a listed value", valueAsString)
	}

	return g.Error("Unknown property: %T.%s", thisTestObject, propertyName)
}
