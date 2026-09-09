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

func NewUserGroup(id goald.BObjID) *UserGroup {
	// TODO use sync.Pool?
	newUserGroup := &UserGroup{}
	newUserGroup.ID = id

	return newUserGroup
}

func GetUserGroupFrom(cache *goald.BObjCache, id goald.BObjID) *UserGroup {
	if cachedUserGroup := cache.Get("UserGroup", id); cachedUserGroup != nil {
		return cachedUserGroup.(*UserGroup)
	}

	return nil
}

func CachedOrNewUserGroup(cache *goald.BObjCache, id goald.BObjID) *UserGroup {
	if cachedUserGroup := GetUserGroupFrom(cache, id); cachedUserGroup != nil {
		return cachedUserGroup
	}

	return cache.Set(NewUserGroup(id))
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a UserGroup, without using reflection
func (bo *UserGroup) GetModelName() utils.ModelName {
	return "UserGroup"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *UserGroup) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Creation":
		return core.DateToString(bo.Creation)
	case "Description":
		return bo.Description
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "Name":
		return bo.Name
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *UserGroup) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
	case "Description":
		bo.Description = valueAsString
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "Name":
		bo.Name = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// ------------------------------------------------------------------------------------------------
// Explicit relationship access
// ------------------------------------------------------------------------------------------------

func (bo *UserGroup) WithAddedMembers(added goald.IUser) goald.IUser {
	bo.Members = append(bo.Members, added)
	return added
}

// ------------------------------------------------------------------------------------------------
// Generic relationship access
// ------------------------------------------------------------------------------------------------

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *UserGroup) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *UserGroup) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "Members":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'UserGroup.Members', got %T", value)
		}
		bo.Members = append(bo.Members, targetValue)
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *UserGroup) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {
	case "Members":
		bo.Members = []goald.IUser{}
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *UserGroup) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *UserGroup) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {
	case "Members":
		members := make([]goald.IBusinessObject, len(bo.Members))
		for i, target := range bo.Members {
			members[i] = target
		}
		return members, nil
	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *UserGroup) IsModelValid() error {
	if err := goald.CheckStringSize(bo.Description, 64, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Description'")
	}
	if bo.Name == "" {
		return goald.Error("'Name' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.Name, 24, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Name'")
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *UserGroup) RemoveCycles() {
}
