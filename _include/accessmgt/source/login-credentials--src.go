// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

type LoginCredentialsModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForLoginCredentials(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &LoginCredentialsModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "LoginCredentials", lastMod)}
}

func (this *LoginCredentialsModelSource) NewObject() any {
	return &accessmgt.LoginCredentials{}
}

func (this *LoginCredentialsModelSource) NewSlice() any {
	return &[]*accessmgt.LoginCredentials{}
}

func (this *LoginCredentialsModelSource) AppendToSlice(slicePtr any, bObj any) any {
	s := slicePtr.(*[]*accessmgt.LoginCredentials)
	*s = append(*s, bObj.(*accessmgt.LoginCredentials))
	return s
}
