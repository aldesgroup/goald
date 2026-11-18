// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

type PermissionOwnerModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForPermissionOwner(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &PermissionOwnerModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "PermissionOwner", lastMod)}
}

func (this *PermissionOwnerModelSource) NewObject() any {
	return &accessmgt.PermissionOwner{}
}

func (this *PermissionOwnerModelSource) NewSlice() any {
	return &[]*accessmgt.PermissionOwner{}
}

func (this *PermissionOwnerModelSource) AppendToSlice(slicePtr any, bObj any) any {
	s := slicePtr.(*[]*accessmgt.PermissionOwner)
	*s = append(*s, bObj.(*accessmgt.PermissionOwner))
	return s
}
