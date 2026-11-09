package accessmgt

import "github.com/aldesgroup/goald"

// Type check
var _ goald.IUser = (*User)(nil)

// Implementing goald.IUser interface for the User business object
func (thisUser *User) GetUniqueAuthID() string {
	return thisUser.ExternalID
}

// GetMemberships implements [goald.IUser].
func (thisUser *User) GetMemberships() []goald.IUserGroup {
	return thisUser.MemberOf
}
