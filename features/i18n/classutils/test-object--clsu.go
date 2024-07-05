// Generated file, do not edit!
package classutils

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

type TestObjectClassUtils struct {
	goald.IClassUtilsCore
}

func ClassUtilsForTestObject(srcPath, lastMod string) goald.IClassUtils {
	return &TestObjectClassUtils{IClassUtilsCore: goald.NewClassUtilsCore(srcPath, lastMod)}
}

func (thisUtils *TestObjectClassUtils) NewObject() any {
	return &i18n.TestObject{}
}

func (thisUtils *TestObjectClassUtils) NewSlice() any {
	return []*i18n.TestObject{}
}
