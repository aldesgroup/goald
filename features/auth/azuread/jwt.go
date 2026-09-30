// ------------------------------------------------------------------------------------------------
// Minimal, dependency-free JWT (RS256) verification + JWKS fetching/caching, used to validate
// bearer tokens issued by Microsoft Entra ID / Entra External ID, via their standard OIDC
// discovery document and JSON Web Key Set - see https://openid.net/specs/openid-connect-discovery-1_0.html
// ------------------------------------------------------------------------------------------------
package azuread

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ------------------------------------------------------------------------------------------------
// JWT header & signature verification
// ------------------------------------------------------------------------------------------------

type jwtHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

// parseJWTHeader reads (without verifying) the header segment of a compact JWT.
func parseJWTHeader(rawToken string) (*jwtHeader, error) {
	headerPart, _, found := strings.Cut(rawToken, ".")
	if !found {
		return nil, fmt.Errorf("malformed token: expected at least a header and a payload segment")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(headerPart)
	if err != nil {
		return nil, fmt.Errorf("could not decode token header: %w", err)
	}

	var header jwtHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("could not parse token header: %w", err)
	}

	return &header, nil
}

// verifyRS256 checks the RS256 signature of a compact JWT against the given public key, and
// returns its (still unvalidated) claims, i.e. the payload segment decoded as a generic map.
func verifyRS256(rawToken string, publicKey *rsa.PublicKey) (map[string]any, error) {
	parts := strings.Split(rawToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed token: expected 3 segments, got %d", len(parts))
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("could not decode token signature: %w", err)
	}

	hashed := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); err != nil {
		return nil, fmt.Errorf("token signature verification failed: %w", err)
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("could not decode token payload: %w", err)
	}

	var claims map[string]any
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("could not parse token payload: %w", err)
	}

	return claims, nil
}

// ------------------------------------------------------------------------------------------------
// OIDC discovery & JWKS fetching
// ------------------------------------------------------------------------------------------------

type oidcDiscoveryDoc struct {
	Issuer        string `json:"issuer"`
	JWKSURI       string `json:"jwks_uri"`
	TokenEndpoint string `json:"token_endpoint"`
}

type jsonWebKeySet struct {
	Keys []struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

// parseRSAPublicKey builds an *rsa.PublicKey from a JWK's base64url-encoded modulus (n) and
// exponent (e), as found in a standard JSON Web Key Set.
func parseRSAPublicKey(nB64, eB64 string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, fmt.Errorf("invalid modulus: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, fmt.Errorf("invalid exponent: %w", err)
	}

	var exponent int
	for _, b := range eBytes {
		exponent = exponent<<8 | int(b)
	}

	return &rsa.PublicKey{N: new(big.Int).SetBytes(nBytes), E: exponent}, nil
}

// ------------------------------------------------------------------------------------------------
// Caching the signing keys, per issuer (1 issuer = 1 tenant), so we don't hit the network on
// every single request
// ------------------------------------------------------------------------------------------------

const jwksCacheTTL = 1 * time.Hour

type cachedKeySet struct {
	keysByKid map[string]*rsa.PublicKey
	fetchedAt time.Time
}

type jwksCache struct {
	mx       sync.Mutex
	byIssuer map[string]*cachedKeySet
	client   *http.Client
}

func newJWKSCache(client *http.Client) *jwksCache {
	return &jwksCache{
		byIssuer: map[string]*cachedKeySet{},
		client:   client,
	}
}

// getKey returns the RSA public key to use for the given issuer & key ID, fetching (or
// refreshing) the issuer's JWKS as needed - in particular when the key ID isn't found in a
// still-fresh cache entry, which handles the identity provider's periodic key rotation.
func (c *jwksCache) getKey(ctx context.Context, issuer, kid string) (*rsa.PublicKey, error) {
	c.mx.Lock()
	set := c.byIssuer[issuer]
	c.mx.Unlock()

	if set != nil && time.Since(set.fetchedAt) < jwksCacheTTL {
		if key, found := set.keysByKid[kid]; found {
			return key, nil
		}
	}

	freshSet, err := c.fetchKeySet(ctx, issuer)
	if err != nil {
		// falling back to a stale cache entry rather than failing outright, if we have one
		if set != nil {
			if key, found := set.keysByKid[kid]; found {
				return key, nil
			}
		}

		return nil, err
	}

	c.mx.Lock()
	c.byIssuer[issuer] = freshSet
	c.mx.Unlock()

	key, found := freshSet.keysByKid[kid]
	if !found {
		return nil, fmt.Errorf("no signing key found for kid '%s' at issuer '%s'", kid, issuer)
	}

	return key, nil
}

func (c *jwksCache) fetchKeySet(ctx context.Context, issuer string) (*cachedKeySet, error) {
	discoveryURL := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"

	var doc oidcDiscoveryDoc
	if err := c.getJSON(ctx, discoveryURL, &doc); err != nil {
		return nil, fmt.Errorf("could not fetch OIDC discovery document: %w", err)
	}

	if doc.JWKSURI == "" {
		return nil, fmt.Errorf("OIDC discovery document at '%s' has no 'jwks_uri'", discoveryURL)
	}

	var jwks jsonWebKeySet
	if err := c.getJSON(ctx, doc.JWKSURI, &jwks); err != nil {
		return nil, fmt.Errorf("could not fetch JWKS: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey, len(jwks.Keys))
	for _, key := range jwks.Keys {
		if key.Kty != "RSA" || key.N == "" || key.E == "" {
			continue // only RSA signing keys are relevant to us
		}

		publicKey, err := parseRSAPublicKey(key.N, key.E)
		if err != nil {
			continue // ignoring any key we can't parse, rather than failing the whole set
		}

		keys[key.Kid] = publicKey
	}

	return &cachedKeySet{keysByKid: keys, fetchedAt: time.Now()}, nil
}

func (c *jwksCache) getJSON(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("'%s' returned status %d", url, resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
