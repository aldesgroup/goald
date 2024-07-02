// Generated file, do not edit!
package classutils

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

// getting a property's value as a string, without using reflection
func (thisTranslationUrlParamsClassUtils *TranslationUrlParamsClassUtils) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "Key":
		return bo.(*i18n.TranslationUrlParams).Key
	case "Part":
		return bo.(*i18n.TranslationUrlParams).Part
	case "Route":
		return bo.(*i18n.TranslationUrlParams).Route
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTranslationUrlParamsClassUtils *TranslationUrlParamsClassUtils) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "Key":
		bo.(*i18n.TranslationUrlParams).Key = valueAsString
	case "Part":
		bo.(*i18n.TranslationUrlParams).Part = valueAsString
	case "Route":
		bo.(*i18n.TranslationUrlParams).Route = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
