package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
)

type DeviceBootstrap struct {
	goald.BusinessObject
	Status      BootstrapStatus `json:"status"      io:"o*" desc:"the status of the bootstrap process for this device"`
	IotCertPEM  string          `json:"iotCertPem"  io:"o*" desc:"the IoT certificate in PEM format"`       // once the bootstrap is done in the provisioning service
	IotChainPEM string          `json:"iotChainPem" io:"o*" desc:"the IoT certificate chain in PEM format"` // once the bootstrap is done in the provisioning service
	DeviceID    string          `json:"deviceId"    io:"o*" desc:"the attributed device ID"`                // once the bootstrap is done in the provisioning service
	ScopeID     string          `json:"scopeId"     io:"o*" desc:"the attributed scope ID"`                 // once the bootstrap is done in the provisioning service
	// RetryAfter  int    // in seconds
}

func init() {
	model.DeviceBootstrap().SetDescription("The device bootstrap response")
	model.DeviceBootstrap().SetNotPersisted()
}
