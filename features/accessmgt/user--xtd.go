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

func NewUser(id goald.BObjID) *User {
	// TODO use sync.Pool?
	newUser := &User{}
	newUser.ID = id

	return newUser
}

func GetUserFrom(cache *goald.BObjCache, id goald.BObjID) *User {
	if cachedUser := cache.Get("User", id); cachedUser != nil {
		return cachedUser.(*User)
	}

	return nil
}

func CachedOrNewUser(cache *goald.BObjCache, id goald.BObjID) *User {
	if cachedUser := GetUserFrom(cache, id); cachedUser != nil {
		return cachedUser
	}

	return cache.Set(NewUser(id))
}

// ------------------------------------------------------------------------------------------------
// Cloning
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a User, without using reflection
func (bo *User) Clone(withFields, withRelationships bool) goald.IBusinessObject {
	clone := &User{}
	clone.ID = bo.ID

	if withFields {
		clone.Creation = bo.Creation
		clone.Email = bo.Email
		clone.ExternalID = bo.ExternalID
		clone.FirstName = bo.FirstName
		clone.LastName = bo.LastName
		clone.Modification = bo.Modification
		clone.UserStatus = bo.UserStatus
	}

	if withRelationships {
		clone.MemberOf = bo.MemberOf
		clone.Permissions = bo.Permissions
		clone.Profiles = bo.Profiles
	}

	return clone
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a User, without using reflection
func (bo *User) GetModelName() utils.ModelName {
	return "User"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *User) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Creation":
		return core.DateToString(bo.Creation)
	case "Email":
		return bo.Email
	case "ExternalID":
		return bo.ExternalID
	case "FirstName":
		return bo.FirstName
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "LastName":
		return bo.LastName
	case "Modification":
		return core.DateToString(bo.Modification)
	case "UserStatus":
		return core.IntToString(bo.UserStatus.Val())
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *User) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
	case "Email":
		bo.Email = valueAsString
	case "ExternalID":
		bo.ExternalID = valueAsString
	case "FirstName":
		bo.FirstName = valueAsString
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "LastName":
		bo.LastName = valueAsString
	case "Modification":
		bo.Modification = core.StringToDate(valueAsString, "Modification")
	case "UserStatus":
		bo.UserStatus = goald.UserStatus(core.StringToInt(valueAsString, "UserStatus"))
		core.PanicMsgIf(bo.UserStatus.String() == "", "Could not set 'UserStatus' to %s since it's not a listed value", valueAsString)
	}

	return goald.Error("[SetValueAsString] Unknown property: %T.%s", bo, propertyName)
}

// ------------------------------------------------------------------------------------------------
// Explicit relationship access
// ------------------------------------------------------------------------------------------------

func (bo *User) AddAndReturnMemberOf(added goald.IUserGroup) goald.IUserGroup {
	bo.MemberOf = append(bo.MemberOf, added)
	return added
}

func (bo *User) AddAndReturnPermissions(added goald.IPermission) goald.IPermission {
	bo.Permissions = append(bo.Permissions, added)
	return added
}

func (bo *User) AddAndReturnProfiles(added goald.IUserProfile) goald.IUserProfile {
	bo.Profiles = append(bo.Profiles, added)
	return added
}

// ------------------------------------------------------------------------------------------------
// Generic relationship access
// ------------------------------------------------------------------------------------------------

// setting this User's parent
func (bo *User) SetParent(parent goald.IBusinessObject) {
	// no parent for this model
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *User) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("[SetRelationshipValue] Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *User) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "MemberOf":
		targetValue, ok := value.(goald.IUserGroup)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUserGroup' for 'User.MemberOf', got %T", value)
		}
		bo.MemberOf = append(bo.MemberOf, targetValue)
		return nil
	case "Permissions":
		targetValue, ok := value.(goald.IPermission)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IPermission' for 'User.Permissions', got %T", value)
		}
		bo.Permissions = append(bo.Permissions, targetValue)
		return nil
	case "Profiles":
		targetValue, ok := value.(goald.IUserProfile)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUserProfile' for 'User.Profiles', got %T", value)
		}
		bo.Profiles = append(bo.Profiles, targetValue)
		return nil
	}

	return goald.Error("[AddRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *User) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {
	case "MemberOf":
		bo.MemberOf = []goald.IUserGroup{}
		return nil
	case "Permissions":
		bo.Permissions = []goald.IPermission{}
		return nil
	case "Profiles":
		bo.Profiles = []goald.IUserProfile{}
		return nil
	}

	return goald.Error("[ClearRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *User) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("[GetSingleRelationshipValue]Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *User) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {
	case "MemberOf":
		memberOf := make([]goald.IBusinessObject, len(bo.MemberOf))
		for i, target := range bo.MemberOf {
			memberOf[i] = target
		}
		return memberOf, nil
	case "Permissions":
		permissions := make([]goald.IBusinessObject, len(bo.Permissions))
		for i, target := range bo.Permissions {
			if err := target.SetRelationshipValue("Owner", bo); err != nil { // ensuring the unique backref is set
				return nil, err
			}
			permissions[i] = target
		}
		return permissions, nil
	case "Profiles":
		profiles := make([]goald.IBusinessObject, len(bo.Profiles))
		for i, target := range bo.Profiles {
			if err := target.SetRelationshipValue("User", bo); err != nil { // ensuring the unique backref is set
				return nil, err
			}
			profiles[i] = target
		}
		return profiles, nil
	}

	return nil, goald.Error("[GetMultipleRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *User) IsModelValid() error {
	if bo.Email == "" {
		return goald.Error("'Email' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.Email, 64, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Email'")
	}
	if err := goald.CheckStringSize(bo.ExternalID, 64, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'ExternalID'")
	}
	if bo.FirstName == "" {
		return goald.Error("'FirstName' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.FirstName, 24, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'FirstName'")
	}
	if bo.LastName == "" {
		return goald.Error("'LastName' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.LastName, 26, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'LastName'")
	}
	if _, isLegitValue := bo.UserStatus.Values()[bo.UserStatus.Val()]; !isLegitValue {
		return goald.Error("Invalid value '%d' for 'User.UserStatus'", bo.UserStatus.Val())
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Diffing
// ------------------------------------------------------------------------------------------------

// Creates 2 synthetic instances gathering the added and removed relationships
func (bo *User) DiffWith(other goald.IBusinessObject, forLinks map[string]bool) (goald.IBusinessObject, goald.IBusinessObject) {
	added := NewUser(bo.ID)
	removed := NewUser(bo.ID)

	otherBo := other.(*User)

	if forLinks["MemberOf"] {
		added.MemberOf, removed.MemberOf = goald.DiffBusinessObjectSlices(otherBo.MemberOf, bo.MemberOf, true)
	}

	return added, removed
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *User) RemoveCycles() {
	for _, target := range bo.Permissions {
		core.PanicIfErr(target.SetRelationshipValue("Owner", nil))
	}
	for _, target := range bo.Profiles {
		core.PanicIfErr(target.SetRelationshipValue("User", nil))
	}
}

// setting the models names on all the business objects associated with this one
func (bo *User) SetModelNames() {
	for _, target := range bo.MemberOf {
		target.SetMdl(target.GetModelName())
		target.SetModelNames()
	}
	for _, target := range bo.Permissions {
		target.SetMdl(target.GetModelName())
		target.SetModelNames()
	}
	for _, target := range bo.Profiles {
		target.SetMdl(target.GetModelName())
		target.SetModelNames()
	}
}