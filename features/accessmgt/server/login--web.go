// ------------------------------------------------------------------------------------------------
// The login endpoints: 1 per realm, each delegating to that realm's configured auth provider (see
// docs/authentication.md for the full picture, curl examples, and a sequence diagram).
// ------------------------------------------------------------------------------------------------
package server

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
	"github.com/aldesgroup/goald/features/auth"
	"github.com/aldesgroup/goald/features/hstatus"
)

var Auth = g.NewEndpointGroup("Authentication", "Logging in as a colleague or a customer, to obtain an access token")

func init() {
	g.PostOneGetOne(handleLoginColleague, nil).
		InGroup(Auth).
		At("login/colleague").
		Public().
		Label("Colleague login").
		Description("Logs a colleague in against the Entra ID (workforce) tenant, returning an access token").
		TrimBodyLogging(1)

	g.PostOneGetOne(handleLoginCustomer, nil).
		InGroup(Auth).
		At("login/customer").
		Public().
		Label("Customer login").
		Description("Logs a customer in against the Entra External ID (CIAM) tenant, returning an access token").
		TrimBodyLogging(1)
}

func handleLoginColleague(webCtx g.WebContext, creds *accessmgt.LoginCredentials) (*accessmgt.AuthToken, hstatus.Code, string) {
	return doLogin(webCtx, auth.RealmColleague, creds)
}

func handleLoginCustomer(webCtx g.WebContext, creds *accessmgt.LoginCredentials) (*accessmgt.AuthToken, hstatus.Code, string) {
	return doLogin(webCtx, auth.RealmCustomer, creds)
}

func doLogin(webCtx g.WebContext, realm auth.Realm, creds *accessmgt.LoginCredentials) (*accessmgt.AuthToken, hstatus.Code, string) {
	tokenSet, err := g.Login(webCtx, realm, creds.Username, creds.Password)
	if err != nil {
		return nil, hstatus.Unauthorized, g.ErrorC(err, "Login failed").Error()
	}

	return &accessmgt.AuthToken{
		AccessToken:  tokenSet.AccessToken,
		TokenType:    tokenSet.TokenType,
		ExpiresIn:    tokenSet.ExpiresIn,
		RefreshToken: tokenSet.RefreshToken,
		IDToken:      tokenSet.IDToken,
	}, hstatus.OK, ""
}
