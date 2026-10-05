// ------------------------------------------------------------------------------------------------
// A goald.IAuthProvider implementation backed by Microsoft Entra ID, usable both for a workforce
// tenant (internal users) and for an Entra External ID / CIAM tenant (external users) - see
// docs/authentication.md for the full picture and the Azure setup steps.
//
// It only validates bearer tokens: ValidateToken verifies a token's signature (RS256) and standard
// claims (issuer, audience, expiry) against the tenant's published JSON Web Key Set. Users sign in
// at Microsoft itself (authorization code + PKCE, e.g. with MSAL), and the app hands this API the
// resulting access token - there's deliberately no password login here.
// ------------------------------------------------------------------------------------------------
package azuread

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/auth"
)

// ProviderType is the value to use in a realm's ProviderConfig.Type to select this implementation.
const ProviderType auth.ProviderType = "azuread"

func init() {
	goald.RegisterAuthProvider(&provider{
		keyfuncs: newKeyfuncCache(&http.Client{Timeout: 10 * time.Second}),
	})
}

type provider struct {
	keyfuncs *keyfuncCache
}

// ProviderType implements [goald.IAuthProvider].
func (*provider) ProviderType() auth.ProviderType {
	return ProviderType
}

// ------------------------------------------------------------------------------------------------
// Deriving the tenant-specific URLs from the config, allowing overrides for non-standard setups
// (e.g. an Entra External ID tenant using a custom domain)
// ------------------------------------------------------------------------------------------------

// authority returns the OAuth2 authority (token endpoint base) for the given config, e.g.
// "https://login.microsoftonline.com/<tenant-id>".
func (p *provider) authority(cfg *auth.ProviderConfig) string {
	if cfg.Authority != "" {
		return strings.TrimSuffix(cfg.Authority, "/")
	}

	return "https://login.microsoftonline.com/" + cfg.TenantID
}

// issuer returns the expected "iss" claim for tokens from this config, e.g.
// "https://login.microsoftonline.com/<tenant-id>/v2.0", which also doubles as the OIDC discovery
// base URL ("<issuer>/.well-known/openid-configuration").
func (p *provider) issuer(cfg *auth.ProviderConfig) string {
	if cfg.Issuer != "" {
		return strings.TrimSuffix(cfg.Issuer, "/")
	}

	return p.authority(cfg) + "/v2.0"
}

// ------------------------------------------------------------------------------------------------
// Validating an incoming bearer token
// ------------------------------------------------------------------------------------------------

// ValidateToken implements [goald.IAuthProvider]. Signature verification, issuer/audience
// matching and expiry are all delegated to github.com/golang-jwt/jwt/v5, using a jwt.Keyfunc
// backed by the issuer's JWKS (github.com/MicahParks/keyfunc) - we just map the resulting claims.
func (p *provider) ValidateToken(ctx context.Context, cfg *auth.ProviderConfig, rawToken string) (*auth.Claims, error) {
	if cfg.TenantID == "" || cfg.Audience == "" {
		return nil, fmt.Errorf("realm '%s' is missing its tenantId or audience", cfg.Realm)
	}

	issuer := p.issuer(cfg)

	kf, err := p.keyfuncs.get(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("could not get the signing keys: %w", err)
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(rawToken, claims, kf.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(cfg.Audience),
		jwt.WithExpirationRequired(),
	); err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	expiry, _ := claims.GetExpirationTime()

	return &auth.Claims{
		Subject:  firstNonEmpty(stringClaim(claims, "oid"), stringClaim(claims, "sub")),
		Realm:    cfg.Realm,
		Email:    firstNonEmpty(stringClaim(claims, "preferred_username"), stringClaim(claims, "email"), stringClaim(claims, "upn")),
		Name:     stringClaim(claims, "name"),
		TenantID: stringClaim(claims, "tid"),
		Roles:    stringSliceClaim(claims, "roles"),
		Expiry:   expiry.Time,
		Raw:      claims,
	}, nil
}

// ------------------------------------------------------------------------------------------------
// Small helpers for reading values out of a jwt.MapClaims
// ------------------------------------------------------------------------------------------------

func stringClaim(claims jwt.MapClaims, name string) string {
	if s, ok := claims[name].(string); ok {
		return s
	}

	return ""
}

func stringSliceClaim(claims jwt.MapClaims, name string) []string {
	raw, ok := claims[name].([]any)
	if !ok {
		return nil
	}

	values := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			values = append(values, s)
		}
	}

	return values
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return ""
}
