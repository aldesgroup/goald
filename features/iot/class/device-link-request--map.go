// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// getting a property's value as a string, without using reflection
func (thisClass *DeviceLinkRequestClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "ID":
		return core.Int64ToString(int64(bo.(*iot.DeviceLinkRequest).ID))
	case "Model":
		return bo.(*iot.DeviceLinkRequest).Model
	case "Serial":
		return bo.(*iot.DeviceLinkRequest).Serial
	case "UserFullName":
		return bo.(*iot.DeviceLinkRequest).UserFullName
	case "UserID":
		return core.IntToString(bo.(*iot.DeviceLinkRequest).UserID)
	case "VerificationCode":
		return bo.(*iot.DeviceLinkRequest).VerificationCode
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *DeviceLinkRequestClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.(*iot.DeviceLinkRequest).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*iot.DeviceLinkRequest).ID"))
	case "Model":
		bo.(*iot.DeviceLinkRequest).Model = valueAsString
	case "Serial":
		bo.(*iot.DeviceLinkRequest).Serial = valueAsString
	case "UserFullName":
		bo.(*iot.DeviceLinkRequest).UserFullName = valueAsString
	case "UserID":
		bo.(*iot.DeviceLinkRequest).UserID = core.StringToInt(valueAsString, "(*iot.DeviceLinkRequest).UserID")
	case "VerificationCode":
		bo.(*iot.DeviceLinkRequest).VerificationCode = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
