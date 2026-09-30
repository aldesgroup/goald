// ------------------------------------------------------------------------------------------------
// Here we define how Goald delegates authentication (login & bearer token validation) to a
// pluggable identity backend. See features/auth for the shared types, and features/auth/azuread
// for the Microsoft Entra ID / Entra External ID implementation - docs/authentication.md has the
// full picture, including the Azure setup steps, curl examples and a sequence diagram.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/auth"
)

// ------------------------------------------------------------------------------------------------
// The interface any identity backend must implement - e.g. Microsoft Entra ID, Auth0, a custom
// homegrown DB-backed login, a fake provider for tests, etc. A single implementation instance is
// shared by every realm configured with the same provider Type; realm-specific settings (tenant,
// client, audience...) are passed in through the *auth.ProviderConfig argument on each call.
// ------------------------------------------------------------------------------------------------

type IAuthProvider interface {
	// ProviderType identifies this implementation, e.g. "azuread"; must match the "type" configured for a realm.
	ProviderType() auth.ProviderType

	// Login exchanges a set of credentials for a token set. Mainly meant for first-party/trusted
	// clients and automated testing - see docs/authentication.md.
	Login(ctx context.Context, cfg *auth.ProviderConfig, creds auth.Credentials) (*auth.TokenSet, error)

	// ValidateToken validates a bearer token (signature, issuer, audience, expiry...), returning
	// the caller's normalized claims.
	ValidateToken(ctx context.Context, cfg *auth.ProviderConfig, rawToken string) (*auth.Claims, error)
}

// ------------------------------------------------------------------------------------------------
// Registering the available provider implementations, by provider type
// ------------------------------------------------------------------------------------------------

var authProviderRegistry = &struct {
	items map[auth.ProviderType]IAuthProvider
	mx    sync.Mutex
}{items: map[auth.ProviderType]IAuthProvider{}}

// RegisterAuthProvider makes an auth provider implementation available under its own ProviderType -
// typically called from such an implementation's package init(), e.g. features/auth/azuread's.
// Panics if a provider is already registered for that type.
func RegisterAuthProvider(provider IAuthProvider) IAuthProvider {
	authProviderRegistry.mx.Lock()
	defer authProviderRegistry.mx.Unlock()

	core.PanicMsgIf(authProviderRegistry.items[provider.ProviderType()] != nil,
		"There's already an auth provider registered for type '%s'", provider.ProviderType())

	authProviderRegistry.items[provider.ProviderType()] = provider

	return provider
}

func getAuthProviderImpl(providerType auth.ProviderType) IAuthProvider {
	authProviderRegistry.mx.Lock()
	defer authProviderRegistry.mx.Unlock()

	return authProviderRegistry.items[providerType]
}

// ------------------------------------------------------------------------------------------------
// Resolving the configured providers, once, at server startup
// ------------------------------------------------------------------------------------------------

// a fully resolved (implementation + config) auth provider, ready to use for 1 realm
type resolvedAuthProvider struct {
	impl IAuthProvider
	cfg  *auth.ProviderConfig
}

// resolveAuthProviders reads the server's configured "Auth" realms, and pairs each of them up with
// its registered provider implementation. Called once at startup, right after the DB schemas are
// resolved. Authentication stays entirely opt-in: a server with no "Auth" section configured at
// all behaves exactly as before, with no gate whatsoever on its endpoints.
func (thisServer *server) resolveAuthProviders() {
	for realmID, cfg := range thisServer.config.base().Auth {
		impl := getAuthProviderImpl(cfg.Type)
		core.PanicMsgIf(impl == nil, "No auth provider registered for type '%s' (auth realm '%s' in config). "+
			"Did you add a blank import for the right package, e.g. _ \"github.com/aldesgroup/goald/features/auth/azuread\"?",
			cfg.Type, realmID)

		if cfg.Realm == "" {
			cfg.Realm = auth.Realm(realmID)
		}

		thisServer.authProviders[cfg.Realm] = &resolvedAuthProvider{impl: impl, cfg: cfg}
		thisServer.Info(fmt.Sprintf("Authentication realm '%s' is served by provider '%s'", cfg.Realm, cfg.Type))
	}
}

// ------------------------------------------------------------------------------------------------
// Using the resolved providers: logging in, and validating incoming bearer tokens
// ------------------------------------------------------------------------------------------------

// Login exchanges a set of credentials for a token set, using whichever provider is configured for
// the given realm. Meant to be called from a custom "login" endpoint's handler - see
// features/accessmgt/server for a working example.
func Login(webCtx WebContext, realm auth.Realm, username, password string) (*auth.TokenSet, error) {
	ctxImpl, ok := webCtx.(*webContextImpl)
	if !ok {
		return nil, Error("Login() can only be called from a real HTTP request context")
	}

	resolved := ctxImpl.server.authProviders[realm]
	if resolved == nil {
		return nil, Error("No auth provider configured for realm '%s'", realm)
	}

	return resolved.impl.Login(context.Background(), resolved.cfg, auth.Credentials{Username: username, Password: password})
}

// authenticateRequest extracts the bearer token from the incoming request, and tries every
// configured realm's provider until one of them accepts it - the token itself carries enough
// information (issuer, audience...) for a provider to quickly reject a token that isn't its own.
// Returns (nil, nil) when this server has no auth provider configured at all, since authentication
// is an entirely opt-in feature.
func (thisServer *server) authenticateRequest(req *http.Request) (*auth.Claims, error) {
	if len(thisServer.authProviders) == 0 {
		return nil, nil
	}

	tokenHeader := thisServer.config.base().AuthTokenHeader
	if tokenHeader == "" {
		tokenHeader = "Authorization"
	}

	authHeader := req.Header.Get(tokenHeader)
	if authHeader == "" {
		return nil, Error("Missing '%s' header", tokenHeader)
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return nil, Error("The '%s' header must use the 'Bearer' scheme", tokenHeader)
	}

	rawToken := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if rawToken == "" {
		return nil, Error("Empty bearer token")
	}

	var lastErr error
	for _, resolved := range thisServer.authProviders {
		claims, err := resolved.impl.ValidateToken(req.Context(), resolved.cfg, rawToken)
		if err == nil {
			return claims, nil
		}

		lastErr = err
	}

	return nil, ErrorC(lastErr, "Invalid or expired token")
}
