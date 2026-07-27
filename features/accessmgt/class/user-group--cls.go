// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
"github.com/aldesgroup/goald/features/accessmgt"
)

type UserGroupClass struct {
	goald.IClassCore
}

func ClassForUserGroup(srcPath, lastMod string) goald.IClass {
	return &UserGroupClass{IClassCore: goald.NewClassCore(srcPath, "UserGroup", lastMod)}
}

func (thisClass *UserGroupClass) NewObject() any {
	return &accessmgt.UserGroup{}
}

func (thisClass *UserGroupClass) NewSlice() any {
	return []*accessmgt.UserGroup{}
}
