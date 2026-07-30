// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceBootstrapModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForDeviceBootstrap(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &DeviceBootstrapModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "DeviceBootstrap", lastMod)}
}

func (this *DeviceBootstrapModelSource) NewObject() any {
	return &iot.DeviceBootstrap{}
}

func (this *DeviceBootstrapModelSource) NewSlice() any {
	return []*iot.DeviceBootstrap{}
}
