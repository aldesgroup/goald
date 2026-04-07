package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
)

type DeviceBootstrap struct {
	goald.BusinessObject
	Status      BootstrapStatus
	IotCertPEM  string // once the bootstrap is done in the provisioning service
	IotChainPEM string // once the bootstrap is done in the provisioning service
	DeviceID    string // once the bootstrap is done in the provisioning service
	ScopeID     string // once the bootstrap is done in the provisioning service
	// RetryAfter  int    // in seconds
}

func init() {
	model.DeviceBootstrap().SetNotPersisted()
}
