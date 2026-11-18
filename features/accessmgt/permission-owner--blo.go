package accessmgt

import "github.com/aldesgroup/goald"

// ----------------------------------------------------------------------------
// Implementing IPermissionOwner
// ----------------------------------------------------------------------------

// type check
var _ goald.IPermissionOwner = (*PermissionOwner)(nil)

// GetPermissions implements [goald.IPermissionOwner].
func (po *PermissionOwner) GetPermissions() []goald.IPermission {
	return po.Permissions
}
