// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

// getting a property's value as a string, without using reflection
func (thisClass *TranslationUrlParamsClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "Key":
		return bo.(*i18n.TranslationUrlParams).Key
	case "Namespace":
		return bo.(*i18n.TranslationUrlParams).Namespace
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *TranslationUrlParamsClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "Key":
		bo.(*i18n.TranslationUrlParams).Key = valueAsString
	case "Namespace":
		bo.(*i18n.TranslationUrlParams).Namespace = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (thisClass *TranslationUrlParamsClass) SetRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (thisClass *TranslationUrlParamsClass) AddRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (thisClass *TranslationUrlParamsClass) ClearRelationshipValue(bo goald.IBusinessObject, relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}
