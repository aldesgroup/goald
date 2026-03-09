package iot

// import (
// 	"crypto/rsa"
// 	"crypto/sha256"
// 	"crypto/x509"
// 	"encoding/base64"
// 	"encoding/json"
// 	"encoding/pem"
// 	"errors"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"
// )

// const (
// 	addr                   = ":8080"
// 	intermediateFactoryPem = "pki/intermediate-usine.crt.pem"
// 	maxSkew                = 2 * time.Minute
// )

// var seenNonces = map[string]time.Time{}

// func loadCert(pemBytes []byte) *x509.Certificate {
// 	block, _ := pem.Decode(pemBytes)
// 	if block == nil || block.Type != "CERTIFICATE" {
// 		log.Fatal("bad cert")
// 	}
// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return cert
// }

// func verifyChain(cert *x509.Certificate, intermPEM []byte) error {
// 	inters := x509.NewCertPool()
// 	if !inters.AppendCertsFromPEM(intermPEM) {
// 		return errors.New("cannot append intermediate")
// 	}
// 	_, err := cert.Verify(x509.VerifyOptions{
// 		Roots:     inters, // ← l’intermediate = trust anchor
// 		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
// 	})
// 	return err
// }

// func verifySignature(cert *x509.Certificate, payload []byte, sigB64 string) error {
// 	pub, ok := cert.PublicKey.(*rsa.PublicKey)
// 	if !ok {
// 		return errors.New("not RSA")
// 	}
// 	sig, _ := base64.StdEncoding.DecodeString(sigB64)
// 	h := sha256.Sum256(payload)
// 	return rsa.VerifyPKCS1v15(pub, cryptoHashSHA256, h[:], sig)
// }

// // Avoid import cycles
// var cryptoHashSHA256 = x509.SHA256WithRSA.SignatureAlgorithm.HashFunc()

// func isFresh(ts int64, nonce string) bool {
// 	now := time.Now()
// 	tm := time.Unix(ts, 0)
// 	if now.Sub(tm) > maxSkew || tm.Sub(now) > maxSkew {
// 		return false
// 	}
// 	if _, exists := seenNonces[nonce]; exists {
// 		return false
// 	}
// 	seenNonces[nonce] = now
// 	return true
// }

// func bootstrap(w http.ResponseWriter, r *http.Request) {

// 	var req bootstrapReq
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "bad json", 400)
// 		return
// 	}

// 	factoryCert := loadCert([]byte(req.FactoryCertPEM))
// 	interm := mustRead(intermediateFactoryPem)

// 	if err := verifyChain(factoryCert, interm); err != nil {
// 		http.Error(w, "invalid cert chain: "+err.Error(), 403)
// 		return
// 	}

// 	plBytes, _ := base64.StdEncoding.DecodeString(req.PayloadB64)

// 	if err := verifySignature(factoryCert, plBytes, req.SignatureB64); err != nil {
// 		http.Error(w, "sig invalid: "+err.Error(), 403)
// 		return
// 	}

// 	var pl payload
// 	if err := json.Unmarshal(plBytes, &pl); err != nil {
// 		http.Error(w, "bad payload", 400)
// 		return
// 	}

// 	if !isFresh(pl.TS, pl.Nonce) {
// 		http.Error(w, "replay or stale", 403)
// 		return
// 	}

// 	resp := map[string]any{
// 		"status":   "recognized",
// 		"serial":   pl.Serial,
// 		"model":    pl.Model,
// 		"deviceCN": factoryCert.Subject.CommonName,
// 	}
// 	json.NewEncoder(w).Encode(resp)
// }

// func mustRead(path string) []byte {
// 	b, err := os.ReadFile(path)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return b
// }

// func main() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/device/bootstrap", bootstrap)
// 	log.Println("Backend listening on", addr, "(behind ACA TLS)")
// 	http.ListenAndServe(addr, mux)
// }
