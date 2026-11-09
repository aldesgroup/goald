// ------------------------------------------------------------------------------------------------
// A goald.IAuthProvider implementation backed by Microsoft Entra ID, usable for a workforce tenant
// as well as for an Entra External ID / CIAM tenant; each realm points to 1 tenant - see
// docs/authentication.md for the full picture and the Azure setup steps.
//
// It only validates bearer tokens: ValidateToken verifies a token's signature (RS256) and standard
// claims (issuer, audience, expiry) against the tenant's published JSON Web Key Set. Users sign in
// at Microsoft itself (authorization code + PKCE, e.g. with MSAL), and the app hands this API the
// resulting access token - there's deliberately no password login here. App-only tokens, from
// external callers using the client credentials flow, are recognized too: see Claims.IsApp & ClientID.
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
	// print the raw and decoded token for debugging purposes
	// fmt.Println("Raw token:", rawToken)
	// if segments := strings.Split(rawToken, "."); len(segments) == 3 {
	// 	decoded, _ := base64.RawURLEncoding.DecodeString(segments[1])
	// 	fmt.Println("Decoded token:", string(decoded))
	// }

	// a realm can't validate anything without knowing its tenant and its own audience
	if cfg.TenantID == "" || cfg.Audience == "" {
		return nil, fmt.Errorf("realm '%s' is missing its tenantId or audience", cfg.Realm)
	}

	issuer := p.issuer(cfg)

	// the signing keys are fetched from the issuer's JWKS, and cached per issuer
	kf, errKeys := p.keyfuncs.get(ctx, issuer)
	if errKeys != nil {
		return nil, fmt.Errorf("could not get the signing keys: %w", errKeys)
	}

	// checks the signature first, then the issuer, audience and expiry; fills in claims
	claims := jwt.MapClaims{}
	if _, errParse := jwt.ParseWithClaims(rawToken, claims, kf.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(cfg.Audience),
		jwt.WithExpirationRequired(),
	); errParse != nil {
		return nil, fmt.Errorf("invalid token: %w", errParse)
	}

	// can't fail on a token that passed the checks above, short of a non-numeric "exp"
	expiry, expErr := claims.GetExpirationTime()
	if expErr != nil {
		return nil, fmt.Errorf("invalid expiration time: %w", expErr)
	}

	roles := stringSliceClaim(claims, "roles")

	// app-only tokens carry app roles but no delegated scopes; "idtyp" is only there if configured as an optional claim
	_, hasScopes := claims["scp"]
	isApp := stringClaim(claims, "idtyp") == "app" || (stringClaim(claims, "idtyp") == "" && !hasScopes && len(roles) > 0)

	// normalizing the claims, whichever token version or tenant type they come from
	return &auth.Claims{
		Subject:  firstNonEmpty(stringClaim(claims, "oid"), stringClaim(claims, "sub")),
		Realm:    cfg.Realm,
		IsApp:    isApp,
		ClientID: firstNonEmpty(stringClaim(claims, "azp"), stringClaim(claims, "appid")),
		Email:    firstNonEmpty(stringClaim(claims, "preferred_username"), stringClaim(claims, "email"), stringClaim(claims, "upn")),
		Name:     stringClaim(claims, "name"),
		TenantID: stringClaim(claims, "tid"),
		Roles:    roles,
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
