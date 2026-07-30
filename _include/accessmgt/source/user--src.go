// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

type UserModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForUser(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &UserModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "User", lastMod)}
}

func (this *UserModelSource) NewObject() any {
	return &accessmgt.User{}
}

func (this *UserModelSource) NewSlice() any {
	return []*accessmgt.User{}
}
