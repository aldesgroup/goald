// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the User model
type UserModel struct {
	PermissionOwnerModel
	email      *g.StringField
	firstName  *g.StringField
	lastName   *g.StringField
	externalId *g.StringField
	userStatus *g.EnumField
	memberOf   *g.Relationship
	profiles   *g.Relationship
}

// this is the main way to refer to the User model in the applicative code
func User() *UserModel {
	userOnce.Do(func() {
		user = NewUserModel()

		// this helps dynamically access to the User model
		g.RegisterModel("User", user)
	})

	return user
}

// internal variables
var (
	user     *UserModel
	userOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewUserModel() *UserModel {
	thisModel := &UserModel{PermissionOwnerModel: *NewPermissionOwnerModel()}
	thisModel.email = g.AddStringField(thisModel, "User", "Email", false)
	thisModel.firstName = g.AddStringField(thisModel, "User", "FirstName", false)
	thisModel.lastName = g.AddStringField(thisModel, "User", "LastName", false)
	thisModel.externalId = g.AddStringField(thisModel, "User", "ExternalID", false)
	thisModel.userStatus = g.AddEnumField(thisModel, "User", "UserStatus", false, "goald.UserStatus")
	thisModel.memberOf = g.AddPolyRelationship(thisModel, "User", "MemberOf", true)
	thisModel.profiles = g.AddPolyRelationship(thisModel, "User", "Profiles", true)

	return thisModel
}

// making sure the User model exists at app startup
func init() {
	User()
}

// accessing all the User model's properties and relation	ships

func (U *UserModel) Email() *g.StringField {
	return U.email
}

func (U *UserModel) FirstName() *g.StringField {
	return U.firstName
}

func (U *UserModel) LastName() *g.StringField {
	return U.lastName
}

func (U *UserModel) ExternalID() *g.StringField {
	return U.externalId
}

func (U *UserModel) UserStatus() *g.EnumField {
	return U.userStatus
}

func (U *UserModel) MemberOf() *g.Relationship {
	return U.memberOf
}

func (U *UserModel) Profiles() *g.Relationship {
	return U.profiles
}
