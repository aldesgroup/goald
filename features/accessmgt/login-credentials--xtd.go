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

func NewLoginCredentials(id goald.BObjID) *LoginCredentials {
	// TODO use sync.Pool?
	newLoginCredentials := &LoginCredentials{}
	newLoginCredentials.ID = id

	return newLoginCredentials
}

func GetLoginCredentialsFrom(cache *goald.BObjCache, id goald.BObjID) *LoginCredentials {
	if cachedLoginCredentials := cache.Get("LoginCredentials", id); cachedLoginCredentials != nil {
		return cachedLoginCredentials.(*LoginCredentials)
	}

	return nil
}

func CachedOrNewLoginCredentials(cache *goald.BObjCache, id goald.BObjID) *LoginCredentials {
	if cachedLoginCredentials := GetLoginCredentialsFrom(cache, id); cachedLoginCredentials != nil {
		return cachedLoginCredentials
	}

	return cache.Set(NewLoginCredentials(id))
}

// ------------------------------------------------------------------------------------------------
// Cloning
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a LoginCredentials, without using reflection
func (bo *LoginCredentials) Clone(withFields, withRelationships bool) goald.IBusinessObject {
	clone := &LoginCredentials{}
	clone.ID = bo.ID

	if withFields {
		clone.Creation = bo.Creation
		clone.Modification = bo.Modification
		clone.Password = bo.Password
		clone.Username = bo.Username
	}

	if withRelationships {

	}

	return clone
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a LoginCredentials, without using reflection
func (bo *LoginCredentials) GetModelName() utils.ModelName {
	return "LoginCredentials"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *LoginCredentials) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Creation":
		return core.DateToString(bo.Creation)
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "Modification":
		return core.DateToString(bo.Modification)
	case "Password":
		return bo.Password
	case "Username":
		return bo.Username
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *LoginCredentials) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "Modification":
		bo.Modification = core.StringToDate(valueAsString, "Modification")
	case "Password":
		bo.Password = valueAsString
	case "Username":
		bo.Username = valueAsString
	}

	return goald.Error("[SetValueAsString] Unknown property: %T.%s", bo, propertyName)
}

// ------------------------------------------------------------------------------------------------
// Explicit relationship access
// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
// Generic relationship access
// ------------------------------------------------------------------------------------------------

// setting this LoginCredentials's parent
func (bo *LoginCredentials) SetParent(parent goald.IBusinessObject) {
	// no parent for this model
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *LoginCredentials) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("[SetRelationshipValue] Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *LoginCredentials) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("[AddRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *LoginCredentials) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("[ClearRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *LoginCredentials) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("[GetSingleRelationshipValue]Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *LoginCredentials) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("[GetMultipleRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *LoginCredentials) IsModelValid() error {
	if bo.Username == "" {
		return goald.Error("'Username' is mandatory and must have a non-zero value")
	}
	if bo.Password == "" {
		return goald.Error("'Password' is mandatory and must have a non-zero value")
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Diffing
// ------------------------------------------------------------------------------------------------

// Creates 2 synthetic instances gathering the added and removed relationships
func (bo *LoginCredentials) DiffWith(other goald.IBusinessObject, forLinks map[string]bool) (goald.IBusinessObject, goald.IBusinessObject) {
	added := NewLoginCredentials(bo.ID)
	removed := NewLoginCredentials(bo.ID)

	return added, removed
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *LoginCredentials) RemoveCycles() {
}

// setting the models names on all the business objects associated with this one
func (bo *LoginCredentials) SetModelNames() {
}
