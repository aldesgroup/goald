package iot

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/hstatus"
)

func init() {
	g.PostOneGetOne(handleCheckBootstrap, "").Label("Allows a device to retrieve its bootstrapping info")
}

func handleCheckBootstrap(webCtx g.WebContext, req *DeviceBootstrapRequest) (*DeviceBootstrap, hstatus.Code, string) {

	// check if there's a boostrap pending for the current user / device SN

	// else

	// check if the user exists

	// check if the device can be enrolled

	// if all ok, contact the DPS, wait until the enrollment is done, then update the device bootstrap infos

	return &DeviceBootstrap{ScopeID: "0000", StatusToRemove: "pending"}, hstatus.Accepted, "Enrollment is pending. Retry later to get the result."
}
