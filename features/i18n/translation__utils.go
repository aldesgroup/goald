// Generated file, do not edit!
package i18n

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
	"strconv"
)

// getting a property's value as a string, without using reflection
func (thisTranslation *Translation) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Lang":
		return strconv.Itoa(thisTranslation.Lang.Val())
	case "LangStr":
		return thisTranslation.LangStr
	case "Route":
		return thisTranslation.Route
	case "Part":
		return thisTranslation.Part
	case "Key":
		return thisTranslation.Key
	case "Value":
		return thisTranslation.Value
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTranslation *Translation) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Lang":
		intValue, errConv := strconv.Atoi(valueAsString)
		utils.PanicErrf(errConv, "Could not set 'Translation.Lang' to %s", valueAsString)
		thisTranslation.Lang = (Language)(intValue)
		utils.PanicIff(thisTranslation.Lang.String() == "", "Could not set 'Translation.Lang' to %s since it's not a listed value", valueAsString)
	case "LangStr":
		thisTranslation.LangStr = valueAsString
	case "Route":
		thisTranslation.Route = valueAsString
	case "Part":
		thisTranslation.Part = valueAsString
	case "Key":
		thisTranslation.Key = valueAsString
	case "Value":
		thisTranslation.Value = valueAsString
	}

	return g.Error("Unknown property: %T.%s", thisTranslation, propertyName)
}
