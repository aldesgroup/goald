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

func NewAuthToken(id goald.BObjID) *AuthToken {
	// TODO use sync.Pool?
	newAuthToken := &AuthToken{}
	newAuthToken.ID = id

	return newAuthToken
}

func GetAuthTokenFrom(cache *goald.BObjCache, id goald.BObjID) *AuthToken {
	if cachedAuthToken := cache.Get("AuthToken", id); cachedAuthToken != nil {
		return cachedAuthToken.(*AuthToken)
	}

	return nil
}

func CachedOrNewAuthToken(cache *goald.BObjCache, id goald.BObjID) *AuthToken {
	if cachedAuthToken := GetAuthTokenFrom(cache, id); cachedAuthToken != nil {
		return cachedAuthToken
	}

	return cache.Set(NewAuthToken(id))
}

// ------------------------------------------------------------------------------------------------
// Cloning
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a AuthToken, without using reflection
func (bo *AuthToken) Clone(withFields, withRelationships bool) goald.IBusinessObject {
	clone := &AuthToken{}
	clone.ID = bo.ID

	if withFields {
		clone.AccessToken = bo.AccessToken
		clone.Creation = bo.Creation
		clone.ExpiresIn = bo.ExpiresIn
		clone.IDToken = bo.IDToken
		clone.Modification = bo.Modification
		clone.RefreshToken = bo.RefreshToken
		clone.TokenType = bo.TokenType
	}

	if withRelationships {

	}

	return clone
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a AuthToken, without using reflection
func (bo *AuthToken) GetModelName() utils.ModelName {
	return "AuthToken"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *AuthToken) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "AccessToken":
		return bo.AccessToken
	case "Creation":
		return core.DateToString(bo.Creation)
	case "ExpiresIn":
		return core.Int64ToString(bo.ExpiresIn)
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "IDToken":
		return bo.IDToken
	case "Modification":
		return core.DateToString(bo.Modification)
	case "RefreshToken":
		return bo.RefreshToken
	case "TokenType":
		return bo.TokenType
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *AuthToken) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "AccessToken":
		bo.AccessToken = valueAsString
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
	case "ExpiresIn":
		bo.ExpiresIn = core.StringToInt64(valueAsString, "ExpiresIn")
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "IDToken":
		bo.IDToken = valueAsString
	case "Modification":
		bo.Modification = core.StringToDate(valueAsString, "Modification")
	case "RefreshToken":
		bo.RefreshToken = valueAsString
	case "TokenType":
		bo.TokenType = valueAsString
	}

	return goald.Error("[SetValueAsString] Unknown property: %T.%s", bo, propertyName)
}

// ------------------------------------------------------------------------------------------------
// Explicit relationship access
// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
// Generic relationship access
// ------------------------------------------------------------------------------------------------

// setting this AuthToken's parent
func (bo *AuthToken) SetParent(parent goald.IBusinessObject) {
	// no parent for this model
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *AuthToken) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("[SetRelationshipValue] Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *AuthToken) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("[AddRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *AuthToken) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("[ClearRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *AuthToken) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("[GetSingleRelationshipValue]Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *AuthToken) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("[GetMultipleRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *AuthToken) IsModelValid() error {
	return nil
}

// ------------------------------------------------------------------------------------------------
// Diffing
// ------------------------------------------------------------------------------------------------

// Creates 2 synthetic instances gathering the added and removed relationships
func (bo *AuthToken) DiffWith(other goald.IBusinessObject, forLinks map[string]bool) (goald.IBusinessObject, goald.IBusinessObject) {
	added := NewAuthToken(bo.ID)
	removed := NewAuthToken(bo.ID)

	return added, removed
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *AuthToken) RemoveCycles() {
}

// setting the models names on all the business objects associated with this one
func (bo *AuthToken) SetModelNames() {
}
