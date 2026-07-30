// Generated file, do not edit!
package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
)

type DeviceBootstrapPayload struct {
	goald.BusinessObject
	Serial    string `json:"serial,omitempty"    io:"i*" desc:"the serial number of the device, which should be unique for each device model and factory certificate"`
	Model     string `json:"model,omitempty"     io:"i*" desc:"the model of the device"`
	Timestamp int64  `json:"timestamp,omitempty" io:"i*" desc:"the timestamp of the payload, to block replays of old payloads; UNIX Epoch time = seconds since 1970 (Go)"`
	Nonce     string `json:"nonce,omitempty"     io:"i*" desc:"the nonce of the payload, to block replays of recent payloads"`
}

func init() {
	model.DeviceBootstrapPayload().SetDescription("The payload sent by a device to the bootstrap endpoint")
	model.DeviceBootstrapPayload().SetNotPersisted() // TDO change
}
