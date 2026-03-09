// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
"github.com/aldesgroup/goald/features/iot"
)

type DeviceBootstrapClass struct {
	goald.IClassCore
}

func ClassForDeviceBootstrap(srcPath, lastMod string) goald.IClass {
	return &DeviceBootstrapClass{IClassCore: goald.NewClassCore(srcPath, "DeviceBootstrap", lastMod)}
}

func (thisClass *DeviceBootstrapClass) NewObject() any {
	return &iot.DeviceBootstrap{}
}

func (thisClass *DeviceBootstrapClass) NewSlice() any {
	return []*iot.DeviceBootstrap{}
}
