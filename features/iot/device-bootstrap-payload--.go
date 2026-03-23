package iot

import (
	"github.com/aldesgroup/goald"
	specs "github.com/aldesgroup/goald/_include/_specs"
)

type DeviceBootstrapPayload struct {
	goald.BusinessObject
	Serial    string
	Model     string
	Timestamp int64  // to block replays of old payloads; UNIX Epoch time = seconds since 1970 (Go)
	Nonce     string // to block replays of recent payloads
}

func init() {
	specs.DeviceBootstrapPayload().SetNotPersisted() // TDO change
}
