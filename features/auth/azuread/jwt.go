// ------------------------------------------------------------------------------------------------
// Discovering & caching, per issuer (1 issuer = 1 tenant), a keyfunc.Keyfunc - a jwt.Keyfunc
// (github.com/golang-jwt/jwt/v5) backed by that issuer's JSON Web Key Set, fetched & refreshed
// automatically by github.com/MicahParks/keyfunc. We only ever want 1 per issuer, since each one
// launches its own background refresh goroutine - hence the cache, kept for the process's life.
// ------------------------------------------------------------------------------------------------
package azuread

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/MicahParks/keyfunc/v3"
)

type keyfuncCache struct {
	mx       sync.Mutex
	byIssuer map[string]keyfunc.Keyfunc
	client   *http.Client
}

func newKeyfuncCache(client *http.Client) *keyfuncCache {
	return &keyfuncCache{byIssuer: map[string]keyfunc.Keyfunc{}, client: client}
}

func (c *keyfuncCache) get(ctx context.Context, issuer string) (keyfunc.Keyfunc, error) {
	c.mx.Lock()
	kf, found := c.byIssuer[issuer]
	c.mx.Unlock()

	if found {
		return kf, nil
	}

	jwksURI, err := c.discoverJWKSURI(ctx, issuer)
	if err != nil {
		return nil, err
	}

	kf, err = keyfunc.NewDefaultCtx(ctx, []string{jwksURI})
	if err != nil {
		return nil, fmt.Errorf("could not build a key function from '%s': %w", jwksURI, err)
	}

	c.mx.Lock()
	c.byIssuer[issuer] = kf
	c.mx.Unlock()

	return kf, nil
}

type oidcDiscoveryDoc struct {
	JWKSURI string `json:"jwks_uri"`
}

// discoverJWKSURI reads an issuer's standard OIDC discovery document to find its JWKS endpoint.
func (c *keyfuncCache) discoverJWKSURI(ctx context.Context, issuer string) (string, error) {
	discoveryURL := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not fetch OIDC discovery document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("'%s' returned status %d", discoveryURL, resp.StatusCode)
	}

	var doc oidcDiscoveryDoc
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return "", fmt.Errorf("could not parse OIDC discovery document: %w", err)
	}

	if doc.JWKSURI == "" {
		return "", fmt.Errorf("OIDC discovery document at '%s' has no 'jwks_uri'", discoveryURL)
	}

	return doc.JWKSURI, nil
}
