package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
)

type DeviceBootstrapPayload struct {
	goald.BusinessObject
	Serial    string
	Model     string
	Timestamp int64  // to block replays of old payloads; UNIX Epoch time = seconds since 1970 (Go)
	Nonce     string // to block replays of recent payloads
}

func init() {
	model.DeviceBootstrapPayload().SetNotPersisted() // TDO change
}
