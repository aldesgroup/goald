// Generated file, do not edit!
package i18n

import (
	"github.com/aldesgroup/goald"
)

// getting a property's value as a string, without using reflection
func (thisTranslation *Translation) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "ID":
		return string(thisTranslation.ID)
	case "Key":
		return thisTranslation.Key
	case "Lang":
		return thisTranslation.Lang
	case "Part":
		return thisTranslation.Part
	case "Route":
		return thisTranslation.Route
	case "Value":
		return thisTranslation.Value
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTranslation *Translation) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		thisTranslation.ID = goald.BObjID(valueAsString)
	case "Key":
		thisTranslation.Key = valueAsString
	case "Lang":
		thisTranslation.Lang = valueAsString
	case "Part":
		thisTranslation.Part = valueAsString
	case "Route":
		thisTranslation.Route = valueAsString
	case "Value":
		thisTranslation.Value = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", thisTranslation, propertyName)
}
