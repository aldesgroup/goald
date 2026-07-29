// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceBootstrapPayloadClass struct {
	goald.IClass
}

func ClassForDeviceBootstrapPayload(srcPath, lastMod string) goald.IClass {
	return &DeviceBootstrapPayloadClass{IClass: goald.NewClass(srcPath, "DeviceBootstrapPayload", lastMod)}
}

func (thisClass *DeviceBootstrapPayloadClass) NewObject() any {
	return &iot.DeviceBootstrapPayload{}
}

func (thisClass *DeviceBootstrapPayloadClass) NewSlice() any {
	return []*iot.DeviceBootstrapPayload{}
}
