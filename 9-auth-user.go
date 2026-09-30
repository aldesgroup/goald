package goald

import "github.com/aldesgroup/goald/features/auth"

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

// ------------------------------------------------------------------------------------------------
// Turning a validated external identity into 1 of the application's own IUser business objects
// ------------------------------------------------------------------------------------------------

// UserResolverFunc maps the claims obtained from a successfully validated bearer token onto 1 of
// the application's own IUser business objects - e.g. looking an existing user up by external ID
// or email, and/or just-in-time provisioning a new one on first login.
type UserResolverFunc func(bloCtx BloContext, claims *auth.Claims) (IUser, error)

var userResolver UserResolverFunc

// RegisterUserResolver plugs in the application-specific logic used to resolve the IUser
// corresponding to an authenticated request's claims. Without it, GetCurrentUser() always returns nil.
func RegisterUserResolver(fn UserResolverFunc) {
	userResolver = fn
}

func getUserResolver() UserResolverFunc {
	return userResolver
}
