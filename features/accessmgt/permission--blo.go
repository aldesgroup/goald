package accessmgt

import "github.com/aldesgroup/goald"

// ----------------------------------------------------------------------------
// Implementing IPermission
// ----------------------------------------------------------------------------

// type check
var _ goald.IPermission = (*Permission)(nil)

// GetOwner implements [goald.IPermission].
func (p *Permission) GetOwner() goald.IPermissionOwner {
	panic("unimplemented")
}

// GetType implements [goald.IPermission].
func (p *Permission) GetType() goald.PermissionType {
	panic("unimplemented")
}

// GetUriPath implements [goald.IPermission].
func (p *Permission) GetUriPath() string {
	panic("unimplemented")
}
