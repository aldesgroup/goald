// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

type PermissionModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForPermission(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &PermissionModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "Permission", lastMod)}
}

func (this *PermissionModelSource) NewObject() any {
	return &accessmgt.Permission{}
}

func (this *PermissionModelSource) NewSlice() any {
	return &[]*accessmgt.Permission{}
}

func (this *PermissionModelSource) AppendToSlice(slicePtr any, bObj any) any {
	s := slicePtr.(*[]*accessmgt.Permission)
	*s = append(*s, bObj.(*accessmgt.Permission))
	return s
}
