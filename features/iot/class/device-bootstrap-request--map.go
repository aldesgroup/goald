// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// getting a property's value as a string, without using reflection
func (thisClass *DeviceBootstrapRequestClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "FactoryCertPEM":
		return bo.(*iot.DeviceBootstrapRequest).FactoryCertPEM
	case "ID":
		return core.Int64ToString(int64(bo.(*iot.DeviceBootstrapRequest).ID))
	case "PayloadB64":
		return bo.(*iot.DeviceBootstrapRequest).PayloadB64
	case "SignatureB64":
		return bo.(*iot.DeviceBootstrapRequest).SignatureB64
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *DeviceBootstrapRequestClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "FactoryCertPEM":
		bo.(*iot.DeviceBootstrapRequest).FactoryCertPEM = valueAsString
	case "ID":
		bo.(*iot.DeviceBootstrapRequest).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*iot.DeviceBootstrapRequest).ID"))
	case "PayloadB64":
		bo.(*iot.DeviceBootstrapRequest).PayloadB64 = valueAsString
	case "SignatureB64":
		bo.(*iot.DeviceBootstrapRequest).SignatureB64 = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapRequestClass) SetRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapRequestClass) AddRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (thisClass *DeviceBootstrapRequestClass) ClearRelationshipValue(bo goald.IBusinessObject, relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}
