// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceBootstrapRequestClass struct {
	goald.IClass
}

func ClassForDeviceBootstrapRequest(srcPath, lastMod string) goald.IClass {
	return &DeviceBootstrapRequestClass{IClass: goald.NewClass(srcPath, "DeviceBootstrapRequest", lastMod)}
}

func (thisClass *DeviceBootstrapRequestClass) NewObject() any {
	return &iot.DeviceBootstrapRequest{}
}

func (thisClass *DeviceBootstrapRequestClass) NewSlice() any {
	return []*iot.DeviceBootstrapRequest{}
}
