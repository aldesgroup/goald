package goald

// Allowing a user or user group to operate on a URI Path
type IPermission interface {
	IBusinessObject
	GetOwner() IPermissionOwner // returns the owner of the authorization
	GetUriPath() string         // returns the allowed or restricted URI path; "/rest/resource/resid" means exact path, "/rest/resource/*" means all subpaths
	GetType() PermissionType    // returns the type of the permission (read, write, delete)
}

// IPermissionOwner represents an entity that can own permissions
type IPermissionOwner interface {
	IBusinessObject
	GetPermissions() []IPermission // returns the list of permissions associated with this entity
}

// IUser represents a user in the system
type IUser interface {
	IPermissionOwner
	GetUniqueAuthID() string      // returns the unique identifier used for authentication
	GetMemberships() []IUserGroup // returns the list of user groups the user is a member of
	GetProfiles() []IUserProfile  // returns the list of profiles associated with this user
	GetUserStatus() UserStatus    // returns the status of the user associated with this user
}

// IUserProfile represents a profile associated with a user in a given application
type IUserProfile interface {
	IBusinessObject
	GetUser() IUser               // returns the user associated with this profile
	GetProfileStatus() UserStatus // returns the status of the user associated with this profile
	// IUser                         // within an application, a user profile should behave like a user
}

// IUserGroup represents a group of users in the system
type IUserGroup interface {
	IPermissionOwner
	GetName() string        // returns the name of the user group
	GetDescription() string // returns the description of the user group
	GetMembers() []IUser    // returns the list of users that are members of this user group
}
