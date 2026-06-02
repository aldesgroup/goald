package azure

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// BuildDPSSASToken builds a SharedAccessSignature for DPS service REST.
//   - dpsHost:  "<dpsName>.azure-devices-provisioning.net"
//   - policy:   "provisioningserviceowner" (or your custom service policy)
//   - keyB64:   base64-encoded primary/secondary key from DPS access policies
func BuildDPSSASToken(dpsHost, policy, keyB64 string, ttl time.Duration) (string, error) {
	resourceURI := url.QueryEscape(dpsHost)
	exp := time.Now().Add(ttl).Unix()
	stringToSign := fmt.Sprintf("%s\n%d", resourceURI, exp)

	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(stringToSign))
	sig := url.QueryEscape(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	// SAS: SharedAccessSignature sr=<resource>&sig=<sig>&se=<exp>&skn=<policy>
	token := fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%s&skn=%s",
		resourceURI, sig, strconv.FormatInt(exp, 10), url.QueryEscape(policy))
	return token, nil
}
