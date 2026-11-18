// Generated file, do not edit!
package accessmgt

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Instantiation / cache retrieval
// ------------------------------------------------------------------------------------------------

func NewBaseUserProfile(id goald.BObjID) *BaseUserProfile {
	// TODO use sync.Pool?
	newBaseUserProfile := &BaseUserProfile{}
	newBaseUserProfile.ID = id

	return newBaseUserProfile
}

func GetBaseUserProfileFrom(cache *goald.BObjCache, id goald.BObjID) *BaseUserProfile {
	if cachedBaseUserProfile := cache.Get("BaseUserProfile", id); cachedBaseUserProfile != nil {
		return cachedBaseUserProfile.(*BaseUserProfile)
	}

	return nil
}

func CachedOrNewBaseUserProfile(cache *goald.BObjCache, id goald.BObjID) *BaseUserProfile {
	if cachedBaseUserProfile := GetBaseUserProfileFrom(cache, id); cachedBaseUserProfile != nil {
		return cachedBaseUserProfile
	}

	return cache.Set(NewBaseUserProfile(id))
}

// ------------------------------------------------------------------------------------------------
// Cloning
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a BaseUserProfile, without using reflection
func (bo *BaseUserProfile) Clone(withFields, withRelationships bool) goald.IBusinessObject {
	clone := &BaseUserProfile{}
	clone.ID = bo.ID

	if withFields {
		clone.Creation = bo.Creation
		clone.Modification = bo.Modification
	}

	if withRelationships {
		clone.User = bo.User
	}

	return clone
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a BaseUserProfile, without using reflection
func (bo *BaseUserProfile) GetModelName() utils.ModelName {
	return "BaseUserProfile"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *BaseUserProfile) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Creation":
		return core.DateToString(bo.Creation)
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "Modification":
		return core.DateToString(bo.Modification)
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *BaseUserProfile) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "Modification":
		bo.Modification = core.StringToDate(valueAsString, "Modification")
	}

	return goald.Error("[SetValueAsString] Unknown property: %T.%s", bo, propertyName)
}

// ------------------------------------------------------------------------------------------------
// Explicit relationship access
// ------------------------------------------------------------------------------------------------


// ------------------------------------------------------------------------------------------------
// Generic relationship access
// ------------------------------------------------------------------------------------------------

// setting this BaseUserProfile's parent
func (bo *BaseUserProfile) SetParent(parent goald.IBusinessObject) {
	// no parent for this model
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *BaseUserProfile) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "User":
		if value == nil {
			bo.User = nil
		} else {
			targetValue, ok := value.(*User)
			if !ok {
				return goald.Error("Expected a value of type '*User' for 'BaseUserProfile.User', got %T", value)
			}
			bo.User = targetValue
		}

		return nil
	}

	return goald.Error("[SetRelationshipValue] Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *BaseUserProfile) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("[AddRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *BaseUserProfile) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("[ClearRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *BaseUserProfile) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {
	case "User":
		return bo.User, nil
	}

	return nil, goald.Error("[GetSingleRelationshipValue]Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *BaseUserProfile) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("[GetMultipleRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *BaseUserProfile) IsModelValid() error {

	return nil
}

// ------------------------------------------------------------------------------------------------
// Diffing
// ------------------------------------------------------------------------------------------------

// Creates 2 synthetic instances gathering the added and removed relationships
func (bo *BaseUserProfile) DiffWith(other goald.IBusinessObject, forLinks map[string]bool) (goald.IBusinessObject, goald.IBusinessObject) {
	added := NewBaseUserProfile(bo.ID)
	removed := NewBaseUserProfile(bo.ID)

	return added, removed
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *BaseUserProfile) RemoveCycles() {
}

// setting the models names on all the business objects associated with this one
func (bo *BaseUserProfile) SetModelNames() {
	if bo.User != nil {
		bo.User.SetModelNames()
	}
}