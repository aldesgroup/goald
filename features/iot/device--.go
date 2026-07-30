// Generated file, do not edit!
package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
)

type Device struct {
	goald.BusinessObject
	Status          DeviceStatus  `json:"status,omitempty"          io:"o*" desc:"the status of the device"`
	StatusString    string        `json:"statusString,omitempty"    io:"o*" desc:"the status of the device"`
	Model           string        `json:"model,omitempty"           io:"o*" desc:"the model of the device"`
	Serial          string        `json:"serial,omitempty"          io:"o*" desc:"the serial number of the device"`
	AssociatedUsers []goald.IUser `json:"associatedUsers,omitempty" io:"o*" desc:"the users associated with the device"` // removed for now
}

func init() {
	model.Device().SetDescription("A device able to connect to the IoT platform and exchange data with it.")
	model.Device().SetNotPersisted()
}
