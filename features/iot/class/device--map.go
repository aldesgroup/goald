// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// getting a property's value as a string, without using reflection
func (thisClass *DeviceClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "ID":
		return core.Int64ToString(int64(bo.(*iot.Device).ID))
	case "Model":
		return bo.(*iot.Device).Model
	case "Serial":
		return bo.(*iot.Device).Serial
	case "Status":
		return core.IntToString(bo.(*iot.Device).Status.Val())
	case "StatusString":
		return bo.(*iot.Device).StatusString
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *DeviceClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.(*iot.Device).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*iot.Device).ID"))
	case "Model":
		bo.(*iot.Device).Model = valueAsString
	case "Serial":
		bo.(*iot.Device).Serial = valueAsString
	case "Status":
		bo.(*iot.Device).Status = iot.DeviceStatus(core.StringToInt(valueAsString, "(*iot.Device).Status"))
		core.PanicMsgIf(bo.(*iot.Device).Status.String() == "", "Could not set '(*iot.Device).Status' to %s since it's not a listed value", valueAsString)
	case "StatusString":
		bo.(*iot.Device).StatusString = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
