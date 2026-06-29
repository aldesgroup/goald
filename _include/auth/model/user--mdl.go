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
	newModel := &UserModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.email = g.NewStringField(newModel, "Email", false)
	newModel.password = g.NewStringField(newModel, "Password", false)
	newModel.firstName = g.NewStringField(newModel, "FirstName", false)
	newModel.lastName = g.NewStringField(newModel, "LastName", false)

	return newModel
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
