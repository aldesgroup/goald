// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceClass struct {
	goald.IClass
}

func ClassForDevice(srcPath, lastMod string) goald.IClass {
	return &DeviceClass{IClass: goald.NewClass(srcPath, "Device", lastMod)}
}

func (thisClass *DeviceClass) NewObject() any {
	return &iot.Device{}
}

func (thisClass *DeviceClass) NewSlice() any {
	return []*iot.Device{}
}
