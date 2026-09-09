package server

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/hstatus"
	"github.com/aldesgroup/goald/features/iot"
)

func init() {
	g.PostOneGetOne(handleDeviceBootstrap, nil).
		InGroup(IoTBackend).
		Label("Bootstrap a device").
		Description("Allows a device to retrieve its bootstrapping info").TrimBodyLogging(45)
}

// this endpoint should be contacted by the device, as many times as necessary, but at least 2 times in reality;
// 1) once to check the device linking has been made by the mobile = rendez-vous made, which should trigger an enrollment with the DPS by the backend
// 2) a final time to get the result of the DPS enrollment, i.e. the signature of the factory cert, and the IoT certs
func handleDeviceBootstrap(webCtx g.WebContext, req *iot.DeviceBootstrapRequest) (*iot.DeviceBootstrap, hstatus.Code, string) {
	result, errBoot := doBootstrapDevice(nil, webCtx.GetBloContext(), req)
	if errBoot != nil {
		errHandle := g.ErrorC(errBoot, "Error while bootstrapping").Error()
		return nil, hstatus.InternalServerError, errHandle
	}

	if result.Status == iot.BootstrapStatusACCESSxDENIED {
		return result, hstatus.Unauthorized, ""
	}

	return result, hstatus.OK, ""
}
