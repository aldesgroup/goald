// ------------------------------------------------------------------------------------------------
// Here we define how Goald delegates authentication (login & bearer token validation) to a
// pluggable identity backend. See features/auth for the shared types, and features/auth/azuread
// for the Microsoft Entra ID / Entra External ID implementation - docs/authentication.md has the
// full picture, including the Azure setup steps, curl examples and a sequence diagram.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/auth"
)

// ------------------------------------------------------------------------------------------------
// The interface any identity backend must implement - e.g. Microsoft Entra ID, Auth0, a custom
// homegrown DB-backed login, a fake provider for tests, etc. A single implementation instance is
// shared by every realm configured with the same provider Type; realm-specific settings (tenant,
// audience...) are passed in through the *auth.ProviderConfig argument on each call.
//
// Real users are expected to sign in at the identity provider itself (e.g. authorization code +
// PKCE), and to hand Goald the resulting bearer token: that's all ValidateToken needs.
// ------------------------------------------------------------------------------------------------

type IAuthProvider interface {
	// ProviderType identifies this implementation, e.g. "azuread"; must match the "type" configured for a realm.
	ProviderType() auth.ProviderType

	// ValidateToken validates a bearer token (signature, issuer, audience, expiry...), returning
	// the caller's normalized claims.
	ValidateToken(ctx context.Context, cfg *auth.ProviderConfig, rawToken string) (*auth.Claims, error)
}

// IPasswordLoginProvider is optionally implemented by providers able to exchange a username and
// password for a token set, which only makes sense for local dev & tests (e.g. "devauth"), never
// for a real identity provider. When the backend rejects the credentials themselves, the returned
// error must wrap auth.ErrInvalidCredentials.
type IPasswordLoginProvider interface {
	Login(ctx context.Context, cfg *auth.ProviderConfig, creds auth.Credentials) (*auth.TokenSet, error)
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
// the given realm - which must support password login (see IPasswordLoginProvider), otherwise the
// error wraps auth.ErrPasswordLoginUnsupported. Meant to be called from a custom "login" endpoint's
// handler - see features/accessmgt/server for a working example.
func Login(webCtx WebContext, realm auth.Realm, username, password string) (*auth.TokenSet, error) {
	ctxImpl, ok := webCtx.(*webContextImpl)
	if !ok {
		return nil, Error("Login() can only be called from a real HTTP request context")
	}

	resolved := ctxImpl.server.authProviders[realm]
	if resolved == nil {
		return nil, Error("No auth provider configured for realm '%s'", realm)
	}

	passwordProvider, ok := resolved.impl.(IPasswordLoginProvider)
	if !ok {
		return nil, fmt.Errorf("realm '%s': %w", realm, auth.ErrPasswordLoginUnsupported)
	}

	return passwordProvider.Login(context.Background(), resolved.cfg, auth.Credentials{Username: username, Password: password})
}

// orderedAuthProviders returns the resolved providers by ascending LoginOrder, then realm name,
// so that anything trying several realms in turn behaves the same way on every run.
func (thisServer *server) orderedAuthProviders() []*resolvedAuthProvider {
	ordered := make([]*resolvedAuthProvider, 0, len(thisServer.authProviders))
	for _, resolved := range thisServer.authProviders {
		ordered = append(ordered, resolved)
	}

	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].cfg.LoginOrder != ordered[j].cfg.LoginOrder {
			return ordered[i].cfg.LoginOrder < ordered[j].cfg.LoginOrder
		}

		return ordered[i].cfg.Realm < ordered[j].cfg.Realm
	})

	return ordered
}

// LoginAny tries the given credentials against every configured realm supporting password login,
// in LoginOrder, and returns the first token set obtained along with the realm that issued it. When
// no realm accepts them, the error wraps auth.ErrInvalidCredentials if every realm rejected the
// credentials themselves; otherwise it's the first other failure met (e.g. a misconfigured provider).
// If no realm supports password login at all, the error wraps auth.ErrPasswordLoginUnsupported.
func LoginAny(webCtx WebContext, username, password string) (*auth.TokenSet, auth.Realm, error) {
	ctxImpl, ok := webCtx.(*webContextImpl)
	if !ok {
		return nil, "", Error("LoginAny() can only be called from a real HTTP request context")
	}

	if len(ctxImpl.server.authProviders) == 0 {
		return nil, "", Error("No authentication realm is configured")
	}

	creds := auth.Credentials{Username: username, Password: password}

	tried := 0
	var firstOtherErr error
	for _, resolved := range ctxImpl.server.orderedAuthProviders() {
		passwordProvider, ok := resolved.impl.(IPasswordLoginProvider)
		if !ok {
			continue
		}

		tried++

		tokenSet, err := passwordProvider.Login(context.Background(), resolved.cfg, creds)
		if err == nil {
			return tokenSet, resolved.cfg.Realm, nil
		}

		if !errors.Is(err, auth.ErrInvalidCredentials) {
			ctxImpl.Warn(fmt.Sprintf("Login against realm '%s' failed: %s", resolved.cfg.Realm, err))

			if firstOtherErr == nil {
				firstOtherErr = err
			}
		}
	}

	if tried == 0 {
		return nil, "", auth.ErrPasswordLoginUnsupported
	}

	if firstOtherErr != nil {
		return nil, "", firstOtherErr
	}

	return nil, "", auth.ErrInvalidCredentials
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
	for _, resolved := range thisServer.orderedAuthProviders() {
		claims, err := resolved.impl.ValidateToken(req.Context(), resolved.cfg, rawToken)
		if err == nil {
			return claims, nil
		}

		lastErr = err
	}

	return nil, ErrorC(lastErr, "Invalid or expired token")
}
