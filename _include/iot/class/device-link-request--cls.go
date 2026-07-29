// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

type DeviceLinkRequestClass struct {
	goald.IClass
}

func ClassForDeviceLinkRequest(srcPath, lastMod string) goald.IClass {
	return &DeviceLinkRequestClass{IClass: goald.NewClass(srcPath, "DeviceLinkRequest", lastMod)}
}

func (thisClass *DeviceLinkRequestClass) NewObject() any {
	return &iot.DeviceLinkRequest{}
}

func (thisClass *DeviceLinkRequestClass) NewSlice() any {
	return []*iot.DeviceLinkRequest{}
}
