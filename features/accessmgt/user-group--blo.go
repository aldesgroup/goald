package accessmgt

import "github.com/aldesgroup/goald"

// type check
var _ goald.IUserGroup = (*UserGroup)(nil)

// GetDescription implements [goald.IUserGroup].
func (u *UserGroup) GetDescription() string {
	return u.Description
}

// GetMembers implements [goald.IUserGroup].
func (u *UserGroup) GetMembers() []goald.IUser {
	return u.Members
}

// GetName implements [goald.IUserGroup].
func (u *UserGroup) GetName() string {
	return u.Name
}
