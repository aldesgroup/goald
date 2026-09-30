// ------------------------------------------------------------------------------------------------
// A goald.IAuthProvider implementation backed by Microsoft Entra ID, usable both for a workforce
// tenant (colleagues) and for an Entra External ID / CIAM tenant (customers) - see
// docs/authentication.md for the full picture, the Azure setup steps, and curl examples.
//
//   - ValidateToken verifies an incoming bearer token's signature (RS256) and standard claims
//     (issuer, audience, expiry) against the tenant's published JSON Web Key Set.
//   - Login exchanges a username/password for a token set, using the OAuth2 Resource Owner
//     Password Credentials grant. This requires the app registration to allow public client flows,
//     and is mainly intended for first-party/trusted clients and automated testing: interactive
//     apps should instead redirect users to Microsoft's own login page (authorization code + PKCE).
//
// ------------------------------------------------------------------------------------------------
package azuread

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
		keyfuncs:   newKeyfuncCache(&http.Client{Timeout: 10 * time.Second}),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	})
}

type provider struct {
	keyfuncs   *keyfuncCache
	httpClient *http.Client
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
// Logging in with a username & password (OAuth2 Resource Owner Password Credentials grant)
// ------------------------------------------------------------------------------------------------

// Login implements [goald.IAuthProvider].
func (p *provider) Login(ctx context.Context, cfg *auth.ProviderConfig, creds auth.Credentials) (*auth.TokenSet, error) {
	if creds.Username == "" || creds.Password == "" {
		return nil, errors.New("username and password are both required")
	}

	if cfg.TenantID == "" || cfg.ClientID == "" {
		return nil, fmt.Errorf("realm '%s' is missing its tenantId or clientId", cfg.Realm)
	}

	scope := cfg.Scope
	if scope == "" {
		scope = "openid profile offline_access"
	}

	form := url.Values{
		"grant_type": {"password"},
		"client_id":  {cfg.ClientID},
		"username":   {creds.Username},
		"password":   {creds.Password},
		"scope":      {scope},
	}
	if cfg.ClientSecret != "" {
		form.Set("client_secret", cfg.ClientSecret)
	}

	tokenURL := p.authority(cfg) + "/oauth2/v2.0/token"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not reach the identity provider: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read the identity provider's response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errBody struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		_ = json.Unmarshal(bodyBytes, &errBody)

		if errBody.Error != "" {
			return nil, fmt.Errorf("login failed: %s (%s)", errBody.Error, firstLine(errBody.Description))
		}

		return nil, fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	var tokenSet auth.TokenSet
	if err := json.Unmarshal(bodyBytes, &tokenSet); err != nil {
		return nil, fmt.Errorf("could not parse the identity provider's token response: %w", err)
	}

	return &tokenSet, nil
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

// firstLine trims Microsoft's verbose, multi-line error descriptions down to their first line.
func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\r\n")
	line, _, _ = strings.Cut(line, "\n")

	return line
}
