// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForDevice(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &DeviceModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "Device", lastMod)}
}

func (this *DeviceModelSource) NewObject() any {
	return &iot.Device{}
}

func (this *DeviceModelSource) NewSlice() any {
	return []*iot.Device{}
}
