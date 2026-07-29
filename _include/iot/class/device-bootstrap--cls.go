// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceBootstrapClass struct {
	goald.IClass
}

func ClassForDeviceBootstrap(srcPath, lastMod string) goald.IClass {
	return &DeviceBootstrapClass{IClass: goald.NewClass(srcPath, "DeviceBootstrap", lastMod)}
}

func (thisClass *DeviceBootstrapClass) NewObject() any {
	return &iot.DeviceBootstrap{}
}

func (thisClass *DeviceBootstrapClass) NewSlice() any {
	return []*iot.DeviceBootstrap{}
}
