// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// getting a property's value as a string, without using reflection
func (thisClass *DeviceBootstrapClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "DeviceID":
		return bo.(*iot.DeviceBootstrap).DeviceID
	case "ID":
		return core.Int64ToString(int64(bo.(*iot.DeviceBootstrap).ID))
	case "IotCertPEM":
		return bo.(*iot.DeviceBootstrap).IotCertPEM
	case "IotChainPEM":
		return bo.(*iot.DeviceBootstrap).IotChainPEM
	case "ScopeID":
		return bo.(*iot.DeviceBootstrap).ScopeID
	case "Status":
		return core.IntToString(bo.(*iot.DeviceBootstrap).Status.Val())
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *DeviceBootstrapClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "DeviceID":
		bo.(*iot.DeviceBootstrap).DeviceID = valueAsString
	case "ID":
		bo.(*iot.DeviceBootstrap).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*iot.DeviceBootstrap).ID"))
	case "IotCertPEM":
		bo.(*iot.DeviceBootstrap).IotCertPEM = valueAsString
	case "IotChainPEM":
		bo.(*iot.DeviceBootstrap).IotChainPEM = valueAsString
	case "ScopeID":
		bo.(*iot.DeviceBootstrap).ScopeID = valueAsString
	case "Status":
		bo.(*iot.DeviceBootstrap).Status = iot.BootstrapStatus(core.StringToInt(valueAsString, "(*iot.DeviceBootstrap).Status"))
		core.PanicMsgIf(bo.(*iot.DeviceBootstrap).Status.String() == "", "Could not set '(*iot.DeviceBootstrap).Status' to %s since it's not a listed value", valueAsString)
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapClass) SetRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapClass) AddRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapClass) ClearRelationshipValue(bo goald.IBusinessObject, relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}
