// Generated file, do not edit!
package accessmgt

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// getting the name of the model for a User, without using reflection
func (bo *User) GetModelName() utils.ModelName {
	return "User"
}

// getting a property's value as a string, without using reflection
func (bo *User) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Email":
		return bo.Email
	case "FirstName":
		return bo.FirstName
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "LastName":
		return bo.LastName
	case "Password":
		return bo.Password
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *User) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Email":
		bo.Email = valueAsString
	case "FirstName":
		bo.FirstName = valueAsString
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "LastName":
		bo.LastName = valueAsString
	case "Password":
		bo.Password = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *User) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
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
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *User) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {
	case "MemberOf":
		bo.MemberOf = []goald.IUserGroup{}
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
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
	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// checking a business object's general validity, without using reflection
func (bo *User) IsModelValid() error {
	if bo.Email == "" {
		return goald.Error("'Email' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.Email, 64, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Email'")
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
	if bo.Password == "" {
		return goald.Error("'Password' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.Password, 64, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Password'")
	}

	return nil
}

// removing any cycles from the business object, without using reflection
func (bo *User) RemoveCycles() {
}
