// ------------------------------------------------------------------------------------------------
// The login endpoints: 1 that tries every realm in turn, and 1 per realm, each delegating to that
// realm's configured auth provider (see docs/authentication.md for the full picture, curl
// examples, and a sequence diagram).
// ------------------------------------------------------------------------------------------------
package server

import (
	"errors"

	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
	"github.com/aldesgroup/goald/features/auth"
	"github.com/aldesgroup/goald/features/hstatus"
)

var Auth = g.NewEndpointGroup("Authentication", "Logging in as an internal or an external user, to obtain an access token")

func init() {
	g.PostOneGetOne(handleLogin, nil).
		InGroup(Auth).
		At("login").
		Public().
		Label("Login").
		Description("Logs a user in against the first realm accepting their credentials, returning an access token and that realm").
		TrimBodyLogging(1)

	g.PostOneGetOne(handleLoginInternal, nil).
		InGroup(Auth).
		At("login/internal").
		Public().
		Label("Internal login").
		Description("Logs an internal user in against the Entra ID (workforce) tenant, returning an access token").
		TrimBodyLogging(1)

	g.PostOneGetOne(handleLoginExternal, nil).
		InGroup(Auth).
		At("login/external").
		Public().
		Label("External login").
		Description("Logs an external user in against the Entra External ID (CIAM) tenant, returning an access token").
		TrimBodyLogging(1)
}

func handleLogin(webCtx g.WebContext, creds *accessmgt.LoginCredentials) (*accessmgt.AuthToken, hstatus.Code, string) {
	tokenSet, realm, err := g.LoginAny(webCtx, creds.Username, creds.Password)
	if err != nil {
		// same answer whichever realm(s) rejected the credentials, so it can't reveal where an account lives
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, hstatus.Unauthorized, "Login failed"
		}

		return nil, hstatus.ServiceUnavailable, "Login is currently unavailable"
	}

	return newAuthToken(tokenSet, realm), hstatus.OK, ""
}

func handleLoginInternal(webCtx g.WebContext, creds *accessmgt.LoginCredentials) (*accessmgt.AuthToken, hstatus.Code, string) {
	return doLogin(webCtx, auth.RealmInternal, creds)
}

func handleLoginExternal(webCtx g.WebContext, creds *accessmgt.LoginCredentials) (*accessmgt.AuthToken, hstatus.Code, string) {
	return doLogin(webCtx, auth.RealmExternal, creds)
}

func doLogin(webCtx g.WebContext, realm auth.Realm, creds *accessmgt.LoginCredentials) (*accessmgt.AuthToken, hstatus.Code, string) {
	tokenSet, err := g.Login(webCtx, realm, creds.Username, creds.Password)
	if err != nil {
		return nil, hstatus.Unauthorized, g.ErrorC(err, "Login failed").Error()
	}

	return newAuthToken(tokenSet, realm), hstatus.OK, ""
}

func newAuthToken(tokenSet *auth.TokenSet, realm auth.Realm) *accessmgt.AuthToken {
	return &accessmgt.AuthToken{
		AccessToken:  tokenSet.AccessToken,
		TokenType:    tokenSet.TokenType,
		ExpiresIn:    tokenSet.ExpiresIn,
		RefreshToken: tokenSet.RefreshToken,
		IDToken:      tokenSet.IDToken,
		Realm:        string(realm),
	}
}
