// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceBootstrapRequestModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForDeviceBootstrapRequest(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &DeviceBootstrapRequestModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "DeviceBootstrapRequest", lastMod)}
}

func (this *DeviceBootstrapRequestModelSource) NewObject() any {
	return &iot.DeviceBootstrapRequest{}
}

func (this *DeviceBootstrapRequestModelSource) NewSlice() any {
	return []*iot.DeviceBootstrapRequest{}
}
