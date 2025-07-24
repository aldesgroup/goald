// Generated file, do not edit!
package class

import (
	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

// getting a property's value as a string, without using reflection
func (thisClass *TranslationClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "ID":
		return core.Int64ToString(int64(bo.(*i18n.Translation).ID))
	case "Key":
		return bo.(*i18n.Translation).Key
	case "Lang":
		return bo.(*i18n.Translation).Lang
	case "Namespace":
		return bo.(*i18n.Translation).Namespace
	case "Value":
		return bo.(*i18n.Translation).Value
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *TranslationClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.(*i18n.Translation).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*i18n.Translation).ID"))
	case "Key":
		bo.(*i18n.Translation).Key = valueAsString
	case "Lang":
		bo.(*i18n.Translation).Lang = valueAsString
	case "Namespace":
		bo.(*i18n.Translation).Namespace = valueAsString
	case "Value":
		bo.(*i18n.Translation).Value = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
