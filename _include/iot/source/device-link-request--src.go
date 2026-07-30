// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceLinkRequestModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForDeviceLinkRequest(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &DeviceLinkRequestModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "DeviceLinkRequest", lastMod)}
}

func (this *DeviceLinkRequestModelSource) NewObject() any {
	return &iot.DeviceLinkRequest{}
}

func (this *DeviceLinkRequestModelSource) NewSlice() any {
	return []*iot.DeviceLinkRequest{}
}
