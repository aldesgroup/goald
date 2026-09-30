// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the LoginCredentials model
type LoginCredentialsModel struct {
	g.IBusinessObjectModel
	username *g.StringField
	password *g.StringField
}

// this is the main way to refer to the LoginCredentials model in the applicative code
func LoginCredentials() *LoginCredentialsModel {
	return loginCredentials
}

// internal variables
var (
	loginCredentials     *LoginCredentialsModel
	loginCredentialsOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewLoginCredentialsModel() *LoginCredentialsModel {
	thisModel := &LoginCredentialsModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.username = g.AddStringField(thisModel, "LoginCredentials", "Username", false)
	thisModel.password = g.AddStringField(thisModel, "LoginCredentials", "Password", false)

	return thisModel
}

// making sure the LoginCredentials model exists at app startup
func init() {
	loginCredentialsOnce.Do(func() {
		loginCredentials = NewLoginCredentialsModel()
	})

	// this helps dynamically access to the LoginCredentials model
	g.RegisterModel("LoginCredentials", loginCredentials)
}

// accessing all the LoginCredentials model's properties and relationships

func (L *LoginCredentialsModel) Username() *g.StringField {
	return L.username
}

func (L *LoginCredentialsModel) Password() *g.StringField {
	return L.password
}
