// Generated file, do not edit!
package i18n

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Instantiation / cache retrieval
// ------------------------------------------------------------------------------------------------

func NewTranslation(id goald.BObjID) *Translation {
	// TODO use sync.Pool?
	newTranslation := &Translation{}
	newTranslation.ID = id

	return newTranslation
}

func GetTranslationFrom(cache *goald.BObjCache, id goald.BObjID) *Translation {
	if cachedTranslation := cache.Get("Translation", id); cachedTranslation != nil {
		return cachedTranslation.(*Translation)
	}

	return nil
}

func CachedOrNewTranslation(cache *goald.BObjCache, id goald.BObjID) *Translation {
	if cachedTranslation := GetTranslationFrom(cache, id); cachedTranslation != nil {
		return cachedTranslation
	}

	return cache.Set(NewTranslation(id))
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a Translation, without using reflection
func (bo *Translation) GetModelName() utils.ModelName {
	return "Translation"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *Translation) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Creation":
		return core.DateToString(bo.Creation)
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
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
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

// ------------------------------------------------------------------------------------------------
// Explicit relationship access
// ------------------------------------------------------------------------------------------------


// ------------------------------------------------------------------------------------------------
// Generic relationship access
// ------------------------------------------------------------------------------------------------

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

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *Translation) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *Translation) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *Translation) IsModelValid() error {
	if bo.Lang == "" {
		return goald.Error("'Lang' is mandatory and must have a non-zero value")
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *Translation) RemoveCycles() {
}
