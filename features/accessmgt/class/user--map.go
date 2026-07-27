// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

// getting a property's value as a string, without using reflection
func (thisClass *UserClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "Email":
		return bo.(*accessmgt.User).Email
	case "FirstName":
		return bo.(*accessmgt.User).FirstName
	case "ID":
		return core.Int64ToString(int64(bo.(*accessmgt.User).ID))
	case "LastName":
		return bo.(*accessmgt.User).LastName
	case "Password":
		return bo.(*accessmgt.User).Password
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *UserClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "Email":
		bo.(*accessmgt.User).Email = valueAsString
	case "FirstName":
		bo.(*accessmgt.User).FirstName = valueAsString
	case "ID":
		bo.(*accessmgt.User).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*accessmgt.User).ID"))
	case "LastName":
		bo.(*accessmgt.User).LastName = valueAsString
	case "Password":
		bo.(*accessmgt.User).Password = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (thisClass *UserClass) SetRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (thisClass *UserClass) AddRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "MemberOf":
		targetValue, ok := value.(goald.IUserGroup)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUserGroup' for 'User.MemberOf', got %T", value)
		}
		bo.(*accessmgt.User).MemberOf = append(bo.(*accessmgt.User).MemberOf, targetValue)
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (thisClass *UserClass) ClearRelationshipValue(bo goald.IBusinessObject, relationshipName string) error {
	switch relationshipName {
	case "MemberOf":
		bo.(*accessmgt.User).MemberOf = []goald.IUserGroup{}
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}
