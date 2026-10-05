// ------------------------------------------------------------------------------------------------
// A goald.IAuthProvider implementation for local development and automated testing - no network
// call, no real identity backend. Any non-empty username/password logs in successfully; the
// password field doubles as a comma-separated list of roles, e.g. password "admin,support" grants
// both roles, so you can exercise authorization logic without a real IdP - see docs/authentication.md.
//
// NEVER select "devauth" as a realm's provider type outside local dev / tests. It's the only shipped
// provider with a password login (goald.IPasswordLoginProvider).
// ------------------------------------------------------------------------------------------------
package devauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/auth"
)

// ProviderType is the value to use in a realm's ProviderConfig.Type to select this implementation.
const ProviderType auth.ProviderType = "devauth"

// devIssuer is the fixed "iss" claim minted & expected for every devauth token.
const devIssuer = "goald-devauth"

// tokenTTL is how long a devauth token stays valid for.
const tokenTTL = 8 * time.Hour

// signingKey is a fixed, process-local HMAC secret - fine since these tokens never leave your
// machine/CI and are only ever meant to be verified by this very same provider.
var signingKey = []byte("goald-devauth-local-only-do-not-use-outside-dev-or-tests")

func init() {
	goald.RegisterAuthProvider(&provider{})
}

type provider struct{}

// ProviderType implements [goald.IAuthProvider].
func (*provider) ProviderType() auth.ProviderType {
	return ProviderType
}

// Login implements [goald.IPasswordLoginProvider] - mints a locally-signed token for any non-empty
// username/password, no real credential check involved.
func (*provider) Login(ctx context.Context, cfg *auth.ProviderConfig, creds auth.Credentials) (*auth.TokenSet, error) {
	if creds.Username == "" || creds.Password == "" {
		return nil, fmt.Errorf("%w: username and password are both required", auth.ErrInvalidCredentials)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   devIssuer,
		"sub":   "dev-" + creds.Username,
		"aud":   cfg.Audience,
		"realm": string(cfg.Realm),
		"email": creds.Username,
		"name":  displayNameFor(creds.Username),
		"roles": rolesFromPassword(creds.Password),
		"iat":   now.Unix(),
		"exp":   now.Add(tokenTTL).Unix(),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
	if err != nil {
		return nil, fmt.Errorf("could not sign the dev token: %w", err)
	}

	return &auth.TokenSet{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int64(tokenTTL.Seconds()),
	}, nil
}

// ValidateToken implements [goald.IAuthProvider].
func (*provider) ValidateToken(ctx context.Context, cfg *auth.ProviderConfig, rawToken string) (*auth.Claims, error) {
	opts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(devIssuer),
		jwt.WithExpirationRequired(),
	}
	if cfg.Audience != "" {
		opts = append(opts, jwt.WithAudience(cfg.Audience))
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(rawToken, claims, func(*jwt.Token) (any, error) {
		return signingKey, nil
	}, opts...); err != nil {
		return nil, fmt.Errorf("invalid dev token: %w", err)
	}

	// each realm only accepts the tokens it minted itself
	if realm, _ := claims["realm"].(string); realm != string(cfg.Realm) {
		return nil, fmt.Errorf("dev token was issued for realm '%s', not '%s'", realm, cfg.Realm)
	}

	expiry, _ := claims.GetExpirationTime()
	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)

	var roles []string
	if raw, ok := claims["roles"].([]any); ok {
		for _, r := range raw {
			if s, ok := r.(string); ok {
				roles = append(roles, s)
			}
		}
	}

	return &auth.Claims{
		Subject: sub,
		Realm:   cfg.Realm,
		Email:   email,
		Name:    name,
		Roles:   roles,
		Expiry:  expiry.Time,
		Raw:     claims,
	}, nil
}

func displayNameFor(username string) string {
	local, _, _ := strings.Cut(username, "@")
	if local == "" {
		return username
	}

	return strings.ToUpper(local[:1]) + local[1:]
}

func rolesFromPassword(password string) []string {
	if password == "" {
		return nil
	}

	return strings.Split(password, ",")
}
