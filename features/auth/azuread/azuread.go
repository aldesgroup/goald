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

	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/auth"
)

// ProviderType is the value to use in a realm's ProviderConfig.Type to select this implementation.
const ProviderType auth.ProviderType = "azuread"

func init() {
	goald.RegisterAuthProvider(&provider{
		jwks:       newJWKSCache(&http.Client{Timeout: 10 * time.Second}),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	})
}

type provider struct {
	jwks       *jwksCache
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

// ValidateToken implements [goald.IAuthProvider].
func (p *provider) ValidateToken(ctx context.Context, cfg *auth.ProviderConfig, rawToken string) (*auth.Claims, error) {
	header, err := parseJWTHeader(rawToken)
	if err != nil {
		return nil, err
	}

	if header.Alg != "RS256" {
		return nil, fmt.Errorf("unsupported signing algorithm '%s'", header.Alg)
	}

	issuer := p.issuer(cfg)

	signingKey, err := p.jwks.getKey(ctx, issuer, header.Kid)
	if err != nil {
		return nil, fmt.Errorf("could not get the signing key: %w", err)
	}

	rawClaims, err := verifyRS256(rawToken, signingKey)
	if err != nil {
		return nil, err
	}

	if iss, _ := rawClaims["iss"].(string); iss != issuer {
		return nil, fmt.Errorf("unexpected issuer '%s'", iss)
	}

	if !audienceMatches(rawClaims["aud"], cfg.Audience) {
		return nil, errors.New("unexpected audience")
	}

	now := time.Now()

	exp, hasExp := numberClaim(rawClaims, "exp")
	if !hasExp {
		return nil, errors.New("token has no 'exp' claim")
	}
	expiry := time.Unix(exp, 0)
	if now.After(expiry) {
		return nil, errors.New("token has expired")
	}

	if nbf, hasNbf := numberClaim(rawClaims, "nbf"); hasNbf && now.Before(time.Unix(nbf, 0)) {
		return nil, errors.New("token is not yet valid")
	}

	return &auth.Claims{
		Subject:  firstNonEmpty(stringClaim(rawClaims, "oid"), stringClaim(rawClaims, "sub")),
		Realm:    cfg.Realm,
		Email:    firstNonEmpty(stringClaim(rawClaims, "preferred_username"), stringClaim(rawClaims, "email"), stringClaim(rawClaims, "upn")),
		Name:     stringClaim(rawClaims, "name"),
		TenantID: stringClaim(rawClaims, "tid"),
		Roles:    stringSliceClaim(rawClaims, "roles"),
		Expiry:   expiry,
		Raw:      rawClaims,
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
// Small helpers for reading values out of a raw, untyped claims map
// ------------------------------------------------------------------------------------------------

func stringClaim(claims map[string]any, name string) string {
	if s, ok := claims[name].(string); ok {
		return s
	}

	return ""
}

func numberClaim(claims map[string]any, name string) (int64, bool) {
	if n, ok := claims[name].(float64); ok {
		return int64(n), true
	}

	return 0, false
}

func stringSliceClaim(claims map[string]any, name string) []string {
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

// audienceMatches checks a JWT "aud" claim against the expected audience - the spec allows "aud"
// to be either a single string or an array of strings.
func audienceMatches(aud any, expected string) bool {
	if expected == "" {
		return false
	}

	switch v := aud.(type) {
	case string:
		return v == expected
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok && s == expected {
				return true
			}
		}
	}

	return false
}

// firstLine trims Microsoft's verbose, multi-line error descriptions down to their first line.
func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\r\n")
	line, _, _ = strings.Cut(line, "\n")

	return line
}
