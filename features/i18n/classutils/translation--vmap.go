// Generated file, do not edit!
package classutils

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

// getting a property's value as a string, without using reflection
func (thisTranslationClassUtils *TranslationClassUtils) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "ID":
		return string(bo.(*i18n.Translation).ID)
	case "Key":
		return bo.(*i18n.Translation).Key
	case "Lang":
		return bo.(*i18n.Translation).Lang
	case "Part":
		return bo.(*i18n.Translation).Part
	case "Route":
		return bo.(*i18n.Translation).Route
	case "Value":
		return bo.(*i18n.Translation).Value
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisTranslationClassUtils *TranslationClassUtils) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.(*i18n.Translation).ID = goald.BObjID(valueAsString)
	case "Key":
		bo.(*i18n.Translation).Key = valueAsString
	case "Lang":
		bo.(*i18n.Translation).Lang = valueAsString
	case "Part":
		bo.(*i18n.Translation).Part = valueAsString
	case "Route":
		bo.(*i18n.Translation).Route = valueAsString
	case "Value":
		bo.(*i18n.Translation).Value = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
