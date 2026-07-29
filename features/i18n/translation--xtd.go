// Generated file, do not edit!
package i18n

import (
	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// getting the name of the class for a Translation, without using reflection
func (bo *Translation) ClassName() utils.ClassName {
	return "Translation"
}

// getting a property's value as a string, without using reflection
func (bo *Translation) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "Key":
		return bo.Key
	case "Lang":
		return bo.Lang
	case "Namespace":
		return bo.Namespace
	case "Value":
		return bo.Value
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *Translation) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "Key":
		bo.Key = valueAsString
	case "Lang":
		bo.Lang = valueAsString
	case "Namespace":
		bo.Namespace = valueAsString
	case "Value":
		bo.Value = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *Translation) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *Translation) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *Translation) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *Translation) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// checking a business object's general validity, without using reflection
func (bo *Translation) IsModelValid() error {
	if bo.Lang == "" {
		return goald.Error("'Lang' is mandatory and must have a non-zero value")
	}

	return nil
}
