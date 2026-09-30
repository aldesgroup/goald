package accessmgt

import "github.com/aldesgroup/goald"

// Type check
var _ goald.IUser = (*User)(nil)

// Implementing goald.IUser interface for the User business object
func (thisUser *User) GetUsername() string {
	return thisUser.Email
}

// GetMemberships implements [goald.IUser].
func (thisUser *User) GetMemberships() []goald.IUserGroup {
	panic("unimplemented")
}
