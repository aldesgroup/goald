package iot

import (
	"github.com/aldesgroup/goald"
	specs "github.com/aldesgroup/goald/_include/_specs"
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
	specs.Device().SetNotPersisted()
}
