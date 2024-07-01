// Generated file, do not edit!
package i18n

import (
	"github.com/aldesgroup/goald"
)

// getting a property's value as a string, without using reflection
func (thisTranslationUrlParams *TranslationUrlParams) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Key":
		return thisTranslationUrlParams.Key
	case "Part":
		return thisTranslationUrlParams.Part
	case "Route":
		return thisTranslationUrlParams.Route
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTranslationUrlParams *TranslationUrlParams) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Key":
		thisTranslationUrlParams.Key = valueAsString
	case "Part":
		thisTranslationUrlParams.Part = valueAsString
	case "Route":
		thisTranslationUrlParams.Route = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", thisTranslationUrlParams, propertyName)
}
