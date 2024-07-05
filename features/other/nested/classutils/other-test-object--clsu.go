// Generated file, do not edit!
package classutils

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/other/nested"
)

type OtherTestObjectClassUtils struct {
	goald.IClassUtilsCore
}

func ClassUtilsForOtherTestObject(srcPath, lastMod string) goald.IClassUtils {
	return &OtherTestObjectClassUtils{IClassUtilsCore: goald.NewClassUtilsCore(srcPath, lastMod)}
}

func (thisUtils *OtherTestObjectClassUtils) NewObject() any {
	return &nested.OtherTestObject{}
}

func (thisUtils *OtherTestObjectClassUtils) NewSlice() any {
	return []*nested.OtherTestObject{}
}
