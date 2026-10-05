// ------------------------------------------------------------------------------------------------
// The password login endpoints, only served by realms whose provider supports it (i.e. "devauth",
// for local dev & tests): 1 that tries every such realm in turn, and 1 per realm. Real identity
// providers answer 501, since their users sign in at the provider itself (authorization code + PKCE)
// and hand over the resulting bearer token - see docs/authentication.md.
// ------------------------------------------------------------------------------------------------
package server

import (
	"errors"

	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
	"github.com/aldesgroup/goald/features/auth"
	"github.com/aldesgroup/goald/features/hstatus"
)

var Auth = g.NewEndpointGroup("Authentication", "Password login for local dev & tests; real users sign in at their identity provider")

func init() {
	g.PostOneGetOne(handleLogin, nil).
		InGroup(Auth).
		At("login").
		Public().
		Label("Login").
		Description("Dev & tests only: logs a user in against the first realm accepting their credentials, returning an access token and that realm").
		TrimBodyLogging(1)

	g.PostOneGetOne(handleLoginInternal, nil).
		InGroup(Auth).
		At("login/internal").
		Public().
		Label("Internal login").
		Description("Dev & tests only: logs an internal user in, returning an access token").
		TrimBodyLogging(1)

	g.PostOneGetOne(handleLoginExternal, nil).
		InGroup(Auth).
		At("login/external").
		Public().
		Label("External login").
		Description("Dev & tests only: logs an external user in, returning an access token").
		TrimBodyLogging(1)
}

func handleLogin(webCtx g.WebContext, creds *accessmgt.LoginCredentials) (*accessmgt.AuthToken, hstatus.Code, string) {
	tokenSet, realm, err := g.LoginAny(webCtx, creds.Username, creds.Password)
	if err != nil {
		// same answer whichever realm(s) rejected the credentials, so it can't reveal where an account lives
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, hstatus.Unauthorized, "Login failed"
		}

		if errors.Is(err, auth.ErrPasswordLoginUnsupported) {
			return nil, hstatus.NotImplemented, auth.ErrPasswordLoginUnsupported.Error()
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
		if errors.Is(err, auth.ErrPasswordLoginUnsupported) {
			return nil, hstatus.NotImplemented, auth.ErrPasswordLoginUnsupported.Error()
		}

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
