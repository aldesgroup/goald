// Generated file, do not edit!
package iot

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// getting the name of the model for a DeviceBootstrap, without using reflection
func (bo *DeviceBootstrap) GetModelName() utils.ModelName {
	return "DeviceBootstrap"
}

// getting a property's value as a string, without using reflection
func (bo *DeviceBootstrap) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "DeviceID":
		return bo.DeviceID
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "IotCertPEM":
		return bo.IotCertPEM
	case "IotChainPEM":
		return bo.IotChainPEM
	case "ScopeID":
		return bo.ScopeID
	case "Status":
		return core.IntToString(bo.Status.Val())
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *DeviceBootstrap) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "DeviceID":
		bo.DeviceID = valueAsString
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "IotCertPEM":
		bo.IotCertPEM = valueAsString
	case "IotChainPEM":
		bo.IotChainPEM = valueAsString
	case "ScopeID":
		bo.ScopeID = valueAsString
	case "Status":
		bo.Status = BootstrapStatus(core.StringToInt(valueAsString, "Status"))
		core.PanicMsgIf(bo.Status.String() == "", "Could not set 'Status' to %s since it's not a listed value", valueAsString)
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *DeviceBootstrap) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *DeviceBootstrap) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *DeviceBootstrap) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *DeviceBootstrap) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// checking a business object's general validity, without using reflection
func (bo *DeviceBootstrap) IsModelValid() error {

	return nil
}

// removing any cycles from the business object, without using reflection
func (bo *DeviceBootstrap) RemoveCycles() {
}
