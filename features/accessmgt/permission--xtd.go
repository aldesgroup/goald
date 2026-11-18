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

func NewPermission(id goald.BObjID) *Permission {
	// TODO use sync.Pool?
	newPermission := &Permission{}
	newPermission.ID = id

	return newPermission
}

func GetPermissionFrom(cache *goald.BObjCache, id goald.BObjID) *Permission {
	if cachedPermission := cache.Get("Permission", id); cachedPermission != nil {
		return cachedPermission.(*Permission)
	}

	return nil
}

func CachedOrNewPermission(cache *goald.BObjCache, id goald.BObjID) *Permission {
	if cachedPermission := GetPermissionFrom(cache, id); cachedPermission != nil {
		return cachedPermission
	}

	return cache.Set(NewPermission(id))
}

// ------------------------------------------------------------------------------------------------
// Cloning
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a Permission, without using reflection
func (bo *Permission) Clone(withFields, withRelationships bool) goald.IBusinessObject {
	clone := &Permission{}
	clone.ID = bo.ID

	if withFields {
		clone.Creation = bo.Creation
		clone.Modification = bo.Modification
	}

	if withRelationships {
		clone.Owner = bo.Owner
	}

	return clone
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a Permission, without using reflection
func (bo *Permission) GetModelName() utils.ModelName {
	return "Permission"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *Permission) GetValueAsString(propertyName string) string {
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
func (bo *Permission) SetValueAsString(propertyName string, valueAsString string) error {
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

// setting this Permission's parent
func (bo *Permission) SetParent(parent goald.IBusinessObject) {
	bo.Owner = parent.(goald.IPermissionOwner)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *Permission) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "Owner":
		if value == nil {
			bo.Owner = nil
		} else {
			targetValue, ok := value.(goald.IPermissionOwner)
			if !ok {
				return goald.Error("Expected a value of type 'goald.IPermissionOwner' for 'Permission.Owner', got %T", value)
			}
			bo.Owner = targetValue
		}

		return nil
	}

	return goald.Error("[SetRelationshipValue] Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *Permission) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("[AddRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *Permission) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("[ClearRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *Permission) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {
	case "Owner":
		return bo.Owner, nil
	}

	return nil, goald.Error("[GetSingleRelationshipValue]Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *Permission) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("[GetMultipleRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *Permission) IsModelValid() error {
	if bo.Owner == nil {
		return goald.Error("'Owner' is required on 'Permission'")
	}
	if bo.Owner.GetID() <= 0 {
		return goald.Error("'Owner' must reference an existing, persisted business object")
	}
	if !core.InSlice([]string{"User", "UserGroup"}, string(bo.Owner.GetModelName())) {
		return goald.Error("Invalid target model for 'Owner'")
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Diffing
// ------------------------------------------------------------------------------------------------

// Creates 2 synthetic instances gathering the added and removed relationships
func (bo *Permission) DiffWith(other goald.IBusinessObject, forLinks map[string]bool) (goald.IBusinessObject, goald.IBusinessObject) {
	added := NewPermission(bo.ID)
	removed := NewPermission(bo.ID)

	return added, removed
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *Permission) RemoveCycles() {
}

// setting the models names on all the business objects associated with this one
func (bo *Permission) SetModelNames() {
	if bo.Owner != nil {
		bo.Owner.SetMdl(bo.Owner.GetModelName())
		bo.Owner.SetModelNames()
	}
}