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
