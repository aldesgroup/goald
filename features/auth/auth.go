// ------------------------------------------------------------------------------------------------
// Shared types for Goald's pluggable authentication.
//
// A Realm identifies a population of callers (e.g. internal vs. external users). Each realm is served
// by its own configured provider (see goald.IAuthProvider), so a single Goald server can
// authenticate against several identity backends at once - typically 1 Microsoft Entra ID
// (workforce) tenant for internal users, and 1 Microsoft Entra External ID (CIAM) tenant for external ones.
// ------------------------------------------------------------------------------------------------
package auth

import (
	"errors"
	"time"
)

// ErrInvalidCredentials is what a provider's Login() must wrap when the identity backend rejects the
// credentials themselves (as opposed to being unreachable or misconfigured), so a login spanning
// several realms knows it can safely try the next one.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Realm identifies a population of callers, each of which can be served by a different,
// independently configured identity provider.
type Realm string

const (
	RealmInternal Realm = "internal" // internal users, typically backed by a Microsoft Entra ID (workforce) tenant
	RealmExternal Realm = "external" // external users, typically backed by a Microsoft Entra External ID (CIAM) tenant
)

// ProviderType identifies which goald.IAuthProvider implementation should serve a given realm,
// e.g. "azuread". This lets an application swap in a different implementation (a different IDP, a
// homegrown DB-backed login, a fake provider for tests, etc.) just by changing the configured type.
type ProviderType string

// ProviderConfig gathers everything a goald.IAuthProvider implementation needs to serve 1 realm.
// A single provider implementation is typically instantiated once, but used for several realms at
// once, each with its own ProviderConfig (e.g. 2 different Entra tenants).
type ProviderConfig struct {
	Type         ProviderType // which IAuthProvider implementation to use, e.g. "azuread"
	Realm        Realm        // which realm this config is for; defaults to the config's map key when empty
	TenantID     string       // the Entra tenant ID (a GUID), for both workforce and External ID tenants
	ClientID     string       // the application (client) ID of the app registration representing this API
	ClientSecret string       // only needed for confidential-client flows; keep this out of source control
	Audience     string       // the expected "aud" claim on incoming access tokens, e.g. "api://<client-id>"
	Scope        string       // the scope(s) to request when logging in, e.g. "api://<client-id>/.default"
	Issuer       string       // overrides the expected "iss" claim & OIDC discovery base; auto-derived from TenantID when empty
	Authority    string       // overrides the OAuth2 authority (token endpoint base); auto-derived from TenantID when empty
	LoginOrder   int          // when logging in without naming a realm, realms are tried by ascending LoginOrder, then by name
}

// Credentials is a generic username/password pair, used for the direct (resource-owner) login flow.
// This flow is mainly meant for first-party/trusted clients and automated testing (see docs/authentication.md);
// interactive, browser-based apps should instead authenticate directly against the identity
// provider and only ever hand this API a bearer token.
type Credentials struct {
	Username string
	Password string
}

// TokenSet mirrors a standard OAuth2/OIDC token response, so it can be serialized as-is to API clients.
type TokenSet struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

// Claims is the normalized result of validating a bearer token, whatever the underlying provider.
type Claims struct {
	Subject  string         // the stable, unique identifier of the caller within its identity provider (e.g. Entra's "oid")
	Realm    Realm          // which realm this caller was authenticated against
	Email    string         // the caller's email / preferred username, if any
	Name     string         // the caller's display name, if any
	TenantID string         // the identity provider's tenant ID, if any
	Roles    []string       // app roles / permissions carried by the token, if any
	Expiry   time.Time      // when the token expires
	Raw      map[string]any // the full set of claims, for anything not normalized above
}
