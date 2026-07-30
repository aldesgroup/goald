// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceBootstrapPayloadModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForDeviceBootstrapPayload(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &DeviceBootstrapPayloadModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "DeviceBootstrapPayload", lastMod)}
}

func (this *DeviceBootstrapPayloadModelSource) NewObject() any {
	return &iot.DeviceBootstrapPayload{}
}

func (this *DeviceBootstrapPayloadModelSource) NewSlice() any {
	return []*iot.DeviceBootstrapPayload{}
}
