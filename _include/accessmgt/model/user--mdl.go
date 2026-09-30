// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the User model
type UserModel struct {
	g.IBusinessObjectModel
	email      *g.StringField
	password   *g.StringField
	firstName  *g.StringField
	lastName   *g.StringField
	externalId *g.StringField
	realm      *g.StringField
	memberOf   *g.Relationship
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
	thisModel := &UserModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.email = g.AddStringField(thisModel, "User", "Email", false)
	thisModel.password = g.AddStringField(thisModel, "User", "Password", false)
	thisModel.firstName = g.AddStringField(thisModel, "User", "FirstName", false)
	thisModel.lastName = g.AddStringField(thisModel, "User", "LastName", false)
	thisModel.externalId = g.AddStringField(thisModel, "User", "ExternalID", false)
	thisModel.realm = g.AddStringField(thisModel, "User", "Realm", false)
	thisModel.memberOf = g.AddPolyRelationship(thisModel, "User", "MemberOf", true)

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

func (U *UserModel) Password() *g.StringField {
	return U.password
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

func (U *UserModel) Realm() *g.StringField {
	return U.realm
}

func (U *UserModel) MemberOf() *g.Relationship {
	return U.memberOf
}
