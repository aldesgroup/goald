package goald

type IUser interface {
	IBusinessObject
	GetUsername() string          // returns the username of the user, which is used for authentication
	GetMemberships() []IUserGroup // returns the list of user groups the user is a member of
}

type IUserGroup interface {
	IBusinessObject
	GetName() string        // returns the name of the user group
	GetDescription() string // returns the description of the user group
	GetMembers() []IUser    // returns the list of users that are members of this user group
}
