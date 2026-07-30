// Generated file, do not edit!
package i18n

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// getting the name of the model for a TranslationUrlParams, without using reflection
func (bo *TranslationUrlParams) GetModelName() utils.ModelName {
	return "TranslationUrlParams"
}

// getting a property's value as a string, without using reflection
func (bo *TranslationUrlParams) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Key":
		return bo.Key
	case "Namespace":
		return bo.Namespace
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *TranslationUrlParams) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Key":
		bo.Key = valueAsString
	case "Namespace":
		bo.Namespace = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *TranslationUrlParams) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *TranslationUrlParams) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *TranslationUrlParams) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *TranslationUrlParams) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// checking a business object's general validity, without using reflection
func (bo *TranslationUrlParams) IsModelValid() error {

	return nil
}

// removing any cycles from the business object, without using reflection
func (bo *TranslationUrlParams) RemoveCycles() {
}
