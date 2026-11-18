// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

type BaseUserProfileModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForBaseUserProfile(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &BaseUserProfileModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "BaseUserProfile", lastMod)}
}

func (this *BaseUserProfileModelSource) NewObject() any {
	return &accessmgt.BaseUserProfile{}
}

func (this *BaseUserProfileModelSource) NewSlice() any {
	return &[]*accessmgt.BaseUserProfile{}
}

func (this *BaseUserProfileModelSource) AppendToSlice(slicePtr any, bObj any) any {
	s := slicePtr.(*[]*accessmgt.BaseUserProfile)
	*s = append(*s, bObj.(*accessmgt.BaseUserProfile))
	return s
}
