// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

// getting a property's value as a string, without using reflection
func (thisClass *UserGroupClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "Description":
		return bo.(*accessmgt.UserGroup).Description
	case "ID":
		return core.Int64ToString(int64(bo.(*accessmgt.UserGroup).ID))
	case "Name":
		return bo.(*accessmgt.UserGroup).Name
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *UserGroupClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "Description":
		bo.(*accessmgt.UserGroup).Description = valueAsString
	case "ID":
		bo.(*accessmgt.UserGroup).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*accessmgt.UserGroup).ID"))
	case "Name":
		bo.(*accessmgt.UserGroup).Name = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (thisClass *UserGroupClass) SetRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (thisClass *UserGroupClass) AddRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "Members":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'UserGroup.Members', got %T", value)
		}
		bo.(*accessmgt.UserGroup).Members = append(bo.(*accessmgt.UserGroup).Members, targetValue)
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (thisClass *UserGroupClass) ClearRelationshipValue(bo goald.IBusinessObject, relationshipName string) error {
	switch relationshipName {
	case "Members":
		bo.(*accessmgt.UserGroup).Members = []goald.IUser{}
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}
