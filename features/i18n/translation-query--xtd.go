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

func NewTranslationQuery(id goald.BObjID) *TranslationQuery {
	// TODO use sync.Pool?
	newTranslationQuery := &TranslationQuery{}
	newTranslationQuery.ID = id

	return newTranslationQuery
}

func GetTranslationQueryFrom(cache *goald.BObjCache, id goald.BObjID) *TranslationQuery {
	if cachedTranslationQuery := cache.Get("TranslationQuery", id); cachedTranslationQuery != nil {
		return cachedTranslationQuery.(*TranslationQuery)
	}

	return nil
}

func CachedOrNewTranslationQuery(cache *goald.BObjCache, id goald.BObjID) *TranslationQuery {
	if cachedTranslationQuery := GetTranslationQueryFrom(cache, id); cachedTranslationQuery != nil {
		return cachedTranslationQuery
	}

	return cache.Set(NewTranslationQuery(id))
}

// ------------------------------------------------------------------------------------------------
// Cloning
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a TranslationQuery, without using reflection
func (bo *TranslationQuery) Clone(withFields, withRelationships bool) goald.IBusinessObject {
	clone := &TranslationQuery{}
	clone.ID = bo.ID

	if withFields {
		clone.Key = bo.Key
		clone.Namespace = bo.Namespace
		clone.Page = bo.Page
		clone.PageSize = bo.PageSize
	}

	if withRelationships {

	}

	return clone
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a TranslationQuery, without using reflection
func (bo *TranslationQuery) GetModelName() utils.ModelName {
	return "TranslationQuery"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *TranslationQuery) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Key":
		return bo.Key
	case "Namespace":
		return bo.Namespace
	case "Page":
		return core.IntToString(bo.Page)
	case "PageSize":
		return core.IntToString(bo.PageSize)
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *TranslationQuery) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Key":
		bo.Key = valueAsString
	case "Namespace":
		bo.Namespace = valueAsString
	case "Page":
		bo.Page = core.StringToInt(valueAsString, "Page")
	case "PageSize":
		bo.PageSize = core.StringToInt(valueAsString, "PageSize")
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
func (bo *TranslationQuery) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *TranslationQuery) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *TranslationQuery) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *TranslationQuery) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *TranslationQuery) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *TranslationQuery) IsModelValid() error {

	return nil
}

// ------------------------------------------------------------------------------------------------
// Diffing
// ------------------------------------------------------------------------------------------------

// Creates 2 synthetic instances gathering the added and removed relationships
func (bo *TranslationQuery) DiffWith(other goald.IBusinessObject, forLinks map[string]bool) (goald.IBusinessObject, goald.IBusinessObject) {
	added := NewTranslationQuery(bo.ID)
	removed := NewTranslationQuery(bo.ID)

	return added, removed
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *TranslationQuery) RemoveCycles() {
}
