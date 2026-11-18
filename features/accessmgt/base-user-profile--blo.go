package accessmgt

import "github.com/aldesgroup/goald"

// ----------------------------------------------------------------------------
// Implementing IUserProfile
// ----------------------------------------------------------------------------

// type check
var _ goald.IUserProfile = (*BaseUserProfile)(nil)

// GetUser implements [IUserProfile].
func (b *BaseUserProfile) GetUser() goald.IUser {
	return b.User
}

// GetProfileStatus implements [IUserProfile].
func (b *BaseUserProfile) GetProfileStatus() goald.UserStatus {
	return b.ProfileStatus
}

// // ----------------------------------------------------------------------------
// // Implementing IUser
// // ----------------------------------------------------------------------------

// // type check
// var _ goald.IUser = (*BaseUserProfile)(nil)

// // GetUniqueAuthID implements [IUser].
// func (b *BaseUserProfile) GetUniqueAuthID() string {
// 	return b.GetUser().GetUniqueAuthID()
// }

// // GetUserStatus implements [IUserProfile].
// func (b *BaseUserProfile) GetUserStatus() goald.UserStatus {
// 	return b.GetUser().GetUserStatus()
// }

// // GetMemberships implements [IUser].
// func (b *BaseUserProfile) GetMemberships() []goald.IUserGroup {
// 	return b.GetUser().GetMemberships()
// }

// // GetProfiles implements [IUser].
// func (b *BaseUserProfile) GetProfiles() []goald.IUserProfile {
// 	return b.GetUser().GetProfiles()
// }

// // GetUserStatus implements [IUser].
// func (b *BaseUserProfile) GetPermissions() []goald.IPermission {
// 	return b.GetUser().GetPermissions()
// }
