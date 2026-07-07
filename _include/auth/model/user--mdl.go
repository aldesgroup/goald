// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the User model
type UserModel struct {
	g.IBusinessObjectModel
	email     *g.StringField
	password  *g.StringField
	firstName *g.StringField
	lastName  *g.StringField
}

// this is the main way to refer to the User model in the applicative code
func User() *UserModel {
	return user
}

// internal variables
var (
	user     *UserModel
	userOnce sync.Once
)

// fully describing each of this class' properties & relationships
func NewUserModel() *UserModel {
	thisModel := &UserModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.email = g.NewStringField(thisModel, "Email", false)
	thisModel.password = g.NewStringField(thisModel, "Password", false)
	thisModel.firstName = g.NewStringField(thisModel, "FirstName", false)
	thisModel.lastName = g.NewStringField(thisModel, "LastName", false)

	return thisModel
}

// making sure the User model exists at app startup
func init() {
	userOnce.Do(func() {
		user = NewUserModel()
	})

	// this helps dynamically access to the User model
	g.RegisterModel("User", user)
}

// accessing all the User class' properties and relationships

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
