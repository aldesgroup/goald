// Generated file, do not edit!
package i18n

import (
	g "github.com/aldesgroup/goald"
)

// getting a property's value as a string, without using reflection
func (thisTranslationUrlParams *TranslationUrlParams) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Route":
		return thisTranslationUrlParams.Route
	case "Part":
		return thisTranslationUrlParams.Part
	case "Key":
		return thisTranslationUrlParams.Key
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTranslationUrlParams *TranslationUrlParams) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Route":
		thisTranslationUrlParams.Route = valueAsString
	case "Part":
		thisTranslationUrlParams.Part = valueAsString
	case "Key":
		thisTranslationUrlParams.Key = valueAsString
	}

	return g.Error("Unknown property: %T.%s", thisTranslationUrlParams, propertyName)
}
