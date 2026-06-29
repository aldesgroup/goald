// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
"github.com/aldesgroup/goald/features/auth"
)

type UserClass struct {
	goald.IClassCore
}

func ClassForUser(srcPath, lastMod string) goald.IClass {
	return &UserClass{IClassCore: goald.NewClassCore(srcPath, "User", lastMod)}
}

func (thisClass *UserClass) NewObject() any {
	return &auth.User{}
}

func (thisClass *UserClass) NewSlice() any {
	return []*auth.User{}
}
