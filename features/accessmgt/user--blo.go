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

// GetUserStatus implements [goald.IUser].
func (thisUser *User) GetUserStatus() goald.UserStatus {
	return thisUser.UserStatus
}

// GetProfiles implements [goald.IUser].
func (thisUser *User) GetProfiles() []goald.IUserProfile {
	return thisUser.Profiles
}
