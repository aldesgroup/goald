package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
)

type Device struct {
	goald.BusinessObject
	Status          DeviceStatus
	StatusString    string
	Model           string
	Serial          string
	AssociatedUsers []goald.IUser
}

func init() {
	model.Device().SetNotPersisted()
}
