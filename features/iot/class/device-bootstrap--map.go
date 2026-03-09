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
	case "ID":
		return core.Int64ToString(int64(bo.(*iot.DeviceBootstrap).ID))
	case "ScopeID":
		return bo.(*iot.DeviceBootstrap).ScopeID
	case "StatusToRemove":
		return bo.(*iot.DeviceBootstrap).StatusToRemove
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *DeviceBootstrapClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.(*iot.DeviceBootstrap).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*iot.DeviceBootstrap).ID"))
	case "ScopeID":
		bo.(*iot.DeviceBootstrap).ScopeID = valueAsString
	case "StatusToRemove":
		bo.(*iot.DeviceBootstrap).StatusToRemove = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
