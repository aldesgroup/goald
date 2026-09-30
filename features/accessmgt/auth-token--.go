// Generated file, do not edit!
package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
)

type AuthToken struct {
	goald.BusinessObject
	AccessToken  string `json:"accessToken,omitempty"  io:"o*" desc:"the bearer token to use in the 'Authorization' header of subsequent API calls"`
	TokenType    string `json:"tokenType,omitempty"    io:"o*" desc:"the type of the returned token, typically 'Bearer'"`
	ExpiresIn    int64  `json:"expiresIn,omitempty"    io:"o*" desc:"the number of seconds until the access token expires"`
	RefreshToken string `json:"refreshToken,omitempty" io:"o*" desc:"a token that can be used to obtain a new access token without logging in again, if any"`
	IDToken      string `json:"idToken,omitempty"      io:"o*" desc:"an OpenID Connect ID token identifying the caller, if any"`
}

func init() {
	model.AuthToken().SetDescription("The token set returned after a successful login")
	model.AuthToken().SetNotPersisted()
	model.AuthToken().AccessToken().SetSecret()
	model.AuthToken().RefreshToken().SetSecret()
	model.AuthToken().IDToken().SetSecret()
}
