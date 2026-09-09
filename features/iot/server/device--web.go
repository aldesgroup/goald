package server

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
	"github.com/aldesgroup/goald/features/hstatus"
	"github.com/aldesgroup/goald/features/iot"
)

var IoTBackend = g.NewEndpointGroup("IoT Backend", "Endpoints to manage IoT devices and their data")

func init() {
	g.PostOneGetOne(handleLinkDevice, nil).
		InGroup(IoTBackend).
		At("link").
		Label("Link a device").
		Description("Allows a user to link to a device, bootstrapping it if it's not enrolled yet")

	g.GetOne(handleGetDevice, nil).
		InGroup(IoTBackend).
		TargetWith(model.Device().Serial()).
		Label("Read a device").
		Description("Allows to read a device with its serial number")
}

func handleLinkDevice(webCtx g.WebContext, req *iot.DeviceLinkRequest) (*iot.Device, hstatus.Code, string) {

	// using the creating / linking function
	device, errCreate := doLinkDevice(webCtx.GetBloContext(), req.Model, req.Serial, req.UserID)
	if errCreate != nil {
		return nil, hstatus.InternalServerError, g.ErrorC(errCreate, "Could not link this device").Error()
	}

	// adjusting the response
	device.StatusString = device.Status.String()

	switch device.Status {
	case iot.DeviceStatusBADxDEVICExSERIAL:
		return nil, hstatus.BadRequest, "The given serial number is not valid"
	case iot.DeviceStatusFORBIDDEN:
		return nil, hstatus.BadRequest, "This model is not handled"
	}

	return device, hstatus.OK, ""
}

func handleGetDevice(webCtx g.WebContext) (*iot.Device, hstatus.Code, string) {
	// TODO better
	device := getDevice(webCtx.GetResourceRefOrID())

	if device == nil {
		return nil, hstatus.NotFound, "No device has been found with this serial: " + webCtx.GetResourceRefOrID()
	}

	return device, hstatus.OK, ""
}
