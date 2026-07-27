// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// getting a property's value as a string, without using reflection
func (thisClass *DeviceBootstrapPayloadClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "ID":
		return core.Int64ToString(int64(bo.(*iot.DeviceBootstrapPayload).ID))
	case "Model":
		return bo.(*iot.DeviceBootstrapPayload).Model
	case "Nonce":
		return bo.(*iot.DeviceBootstrapPayload).Nonce
	case "Serial":
		return bo.(*iot.DeviceBootstrapPayload).Serial
	case "Timestamp":
		return core.Int64ToString(bo.(*iot.DeviceBootstrapPayload).Timestamp)
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *DeviceBootstrapPayloadClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.(*iot.DeviceBootstrapPayload).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*iot.DeviceBootstrapPayload).ID"))
	case "Model":
		bo.(*iot.DeviceBootstrapPayload).Model = valueAsString
	case "Nonce":
		bo.(*iot.DeviceBootstrapPayload).Nonce = valueAsString
	case "Serial":
		bo.(*iot.DeviceBootstrapPayload).Serial = valueAsString
	case "Timestamp":
		bo.(*iot.DeviceBootstrapPayload).Timestamp = core.StringToInt64(valueAsString, "(*iot.DeviceBootstrapPayload).Timestamp")
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapPayloadClass) SetRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapPayloadClass) AddRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapPayloadClass) ClearRelationshipValue(bo goald.IBusinessObject, relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}
