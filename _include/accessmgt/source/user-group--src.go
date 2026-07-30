// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

type UserGroupModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForUserGroup(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &UserGroupModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "UserGroup", lastMod)}
}

func (this *UserGroupModelSource) NewObject() any {
	return &accessmgt.UserGroup{}
}

func (this *UserGroupModelSource) NewSlice() any {
	return []*accessmgt.UserGroup{}
}
