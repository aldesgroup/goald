// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the AuthToken model
type AuthTokenModel struct {
	g.IBusinessObjectModel
	accessToken  *g.StringField
	tokenType    *g.StringField
	expiresIn    *g.BigIntField
	refreshToken *g.StringField
	idToken      *g.StringField
	realm        *g.StringField
}

// this is the main way to refer to the AuthToken model in the applicative code
func AuthToken() *AuthTokenModel {
	return authToken
}

// internal variables
var (
	authToken     *AuthTokenModel
	authTokenOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewAuthTokenModel() *AuthTokenModel {
	thisModel := &AuthTokenModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.accessToken = g.AddStringField(thisModel, "AuthToken", "AccessToken", false)
	thisModel.tokenType = g.AddStringField(thisModel, "AuthToken", "TokenType", false)
	thisModel.expiresIn = g.AddBigIntField(thisModel, "AuthToken", "ExpiresIn", false)
	thisModel.refreshToken = g.AddStringField(thisModel, "AuthToken", "RefreshToken", false)
	thisModel.idToken = g.AddStringField(thisModel, "AuthToken", "IDToken", false)
	thisModel.realm = g.AddStringField(thisModel, "AuthToken", "Realm", false)

	return thisModel
}

// making sure the AuthToken model exists at app startup
func init() {
	authTokenOnce.Do(func() {
		authToken = NewAuthTokenModel()
	})

	// this helps dynamically access to the AuthToken model
	g.RegisterModel("AuthToken", authToken)
}

// accessing all the AuthToken model's properties and relationships

func (A *AuthTokenModel) AccessToken() *g.StringField {
	return A.accessToken
}

func (A *AuthTokenModel) TokenType() *g.StringField {
	return A.tokenType
}

func (A *AuthTokenModel) ExpiresIn() *g.BigIntField {
	return A.expiresIn
}

func (A *AuthTokenModel) RefreshToken() *g.StringField {
	return A.refreshToken
}

func (A *AuthTokenModel) IDToken() *g.StringField {
	return A.idToken
}

func (A *AuthTokenModel) Realm() *g.StringField {
	return A.realm
}
