// Generated file, do not edit!
package i18n

import (
	g "github.com/aldesgroup/goald"
)

// getting a property's value as a string, without using reflection
func (thisTranslation *Translation) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Lang":
		return thisTranslation.Lang
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
		thisTranslation.Lang = valueAsString
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
