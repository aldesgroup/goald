package server

import (
	g "github.com/aldesgroup/goald"
	specs "github.com/aldesgroup/goald/_include/_specs"
	"github.com/aldesgroup/goald/features/hstatus"
	"github.com/aldesgroup/goald/features/iot"
)

func init() {
	g.PostOneGetOne(handleLinkDevice, "").
		At("link").
		Label("Allows a user to link to a device, bootstrapping it if it's not enrolled yet")

	g.GetOne(handleGetDevice, "").TargetWith(specs.Device().Serial())
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
	device := getDevice(webCtx.GetTargetRefOrID())

	if device == nil {
		return nil, hstatus.NotFound, "No device has been found with this serial"
	}

	return device, hstatus.OK, ""
}
