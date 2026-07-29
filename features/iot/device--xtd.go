// Generated file, do not edit!
package iot

import (
	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// getting the name of the class for a Device, without using reflection
func (bo *Device) ClassName() utils.ClassName {
	return "Device"
}

// getting a property's value as a string, without using reflection
func (bo *Device) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "Model":
		return bo.Model
	case "Serial":
		return bo.Serial
	case "Status":
		return core.IntToString(bo.Status.Val())
	case "StatusString":
		return bo.StatusString
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *Device) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "Model":
		bo.Model = valueAsString
	case "Serial":
		bo.Serial = valueAsString
	case "Status":
		bo.Status = DeviceStatus(core.StringToInt(valueAsString, "Status"))
		core.PanicMsgIf(bo.Status.String() == "", "Could not set 'Status' to %s since it's not a listed value", valueAsString)
	case "StatusString":
		bo.StatusString = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *Device) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *Device) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "AssociatedUsers":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'Device.AssociatedUsers', got %T", value)
		}
		bo.AssociatedUsers = append(bo.AssociatedUsers, targetValue)
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *Device) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {
	case "AssociatedUsers":
		bo.AssociatedUsers = []goald.IUser{}
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *Device) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {
	case "AssociatedUsers":
		associatedUsers := make([]goald.IBusinessObject, len(bo.AssociatedUsers))
		for i, target := range bo.AssociatedUsers {
			associatedUsers[i] = target
		}
		return associatedUsers, nil
	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// checking a business object's general validity, without using reflection
func (bo *Device) IsModelValid() error {

	return nil
}
