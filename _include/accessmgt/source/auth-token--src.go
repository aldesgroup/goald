// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

type AuthTokenModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForAuthToken(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &AuthTokenModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "AuthToken", lastMod)}
}

func (this *AuthTokenModelSource) NewObject() any {
	return &accessmgt.AuthToken{}
}

func (this *AuthTokenModelSource) NewSlice() any {
	return &[]*accessmgt.AuthToken{}
}

func (this *AuthTokenModelSource) AppendToSlice(slicePtr any, bObj any) any {
	s := slicePtr.(*[]*accessmgt.AuthToken)
	*s = append(*s, bObj.(*accessmgt.AuthToken))
	return s
}
