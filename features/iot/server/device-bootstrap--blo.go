package server

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log/slog"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// ----------------------------------------------------------------------------
// Init
// ----------------------------------------------------------------------------

const (
	intermediateFactoryCrtPath = "data/certs/intermediate-factory.crt.pem"
)

var (
	intermediateFactoryCert     []byte
	intermediateFactoryCertPool *x509.CertPool
)

func init() {
	intermediateFactoryCert = core.ReadFile(intermediateFactoryCrtPath, true)
	intermediateFactoryCertPool = x509.NewCertPool()
	if !intermediateFactoryCertPool.AppendCertsFromPEM(intermediateFactoryCert) {
		core.PanicMsg("Could not make a cert pool from the intermediate factory cert")
	}
}

// ----------------------------------------------------------------------------
// Main bootstrapping function
// ----------------------------------------------------------------------------

func doBootstrapDevice(bloContext g.BloContext, req *iot.DeviceBootstrapRequest) (*iot.DeviceBootstrap, error) {
	// loading the device's public factory cert
	factoryCert, errLoadCert := loadCert([]byte(req.FactoryCertPEM))
	if errLoadCert != nil {
		return nil, g.ErrorC(errLoadCert, "Error while loading the device's public factory cert")
	}

	// checking the factory cert is legit, i.e. derived from the common intermediate factory cert
	if _, errVerify := factoryCert.Verify(x509.VerifyOptions{
		Roots:     intermediateFactoryCertPool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny}}); errVerify != nil {
		return &iot.DeviceBootstrap{Status: iot.BootstrapStatusACCESSxDENIED}, nil
	}

	// decoding the payload and the signature
	payloadBytes, errDecode := base64.StdEncoding.DecodeString(req.PayloadB64)
	if errDecode != nil {
		return nil, g.ErrorC(errDecode, "Could not decode the payload")
	}
	signatureBytes, errDecode := base64.StdEncoding.DecodeString(req.SignatureB64)
	if errDecode != nil {
		return nil, g.ErrorC(errDecode, "Could not decode the signature")
	}

	// verifying the signature
	if errSignature := verifyPayloadSignature(factoryCert, payloadBytes, signatureBytes); errSignature != nil {
		return &iot.DeviceBootstrap{Status: iot.BootstrapStatusACCESSxDENIED}, nil
	}

	// unmarshalling the payload
	payload := &iot.DeviceBootstrapPayload{}
	if errUnmarsh := json.Unmarshal(payloadBytes, payload); errUnmarsh != nil {
		return nil, g.ErrorC(errUnmarsh, "Payload could not be unmarshaled!")
	}

	// anti-replay measure
	if !isFreshCall(payload) {
		return &iot.DeviceBootstrap{Status: iot.BootstrapStatusACCESSxDENIED}, nil
	}

	// Here, we can enforce having (CN/SAN = serial) in the certs
	// Not doing it for now.

	slog.Debug(fmt.Sprintf("Device %s n°%s is correctly authenticated!", payload.Model, payload.Serial))

	// AT THIS POINT, the device has made a legit bootstrap call,
	// and we can materialise it's will to bootstrap by saving its bootrapping state
	currentBootstrap := iot.GetDeviceBoostrap(payload.Serial)

	// IF the actual bootstrapping with the DPS is not pending nor finished,
	// AND the device has been linked by an authenticated user
	// THEN let's actually initiate the provisioning of the device
	// TODO better with workers etc
	if getDevice(payload.Serial) != nil && currentBootstrap.Status == iot.BootstrapStatusAUTHENTICATEDxONLY {
		// TODO better : ENQUEUE A JOB
		go startProvisioning(payload.Serial, currentBootstrap)
	}

	return currentBootstrap, nil
}

// ----------------------------------------------------------------------------
// Utilities
// ----------------------------------------------------------------------------

// Load a given certificate
func loadCert(certBytes []byte) (*x509.Certificate, error) {
	// Decoding the given certificate
	pemBlock, _ := pem.Decode(certBytes)
	if pemBlock == nil {
		return nil, g.Error("Invalid or empty certificate")
	}
	if pemBlock.Type != "CERTIFICATE" {
		return nil, g.Error("Bad certificate type ('%s' instead of 'CERTIFICATE')", pemBlock.Type)
	}

	// Parsing it to an object
	cert, errParse := x509.ParseCertificate(pemBlock.Bytes)
	if errParse != nil {
		return nil, g.ErrorC(errParse, "Could not parse the given certificate")
	}

	return cert, nil
}

// Verify the RSA-PKCS1v15 SHA-256 signature of the payload, using the factory cert
func verifyPayloadSignature(factoryCert *x509.Certificate, payload []byte, signature []byte) error {
	// retrieving the public key
	publicKey, isRSA := factoryCert.PublicKey.(*rsa.PublicKey)
	if !isRSA {
		return goald.Error("factory public key is not RSA")
	}

	// computing the hash of the payload
	payloadHash := sha256.Sum256(payload)

	// verifyoing the signature
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, payloadHash[:], signature)
}

// isFreshCall checks :
//  1. the timestamp is fresh enough (|now - ts| <= MaxClockSkew)
//  2. the nonce has never been seen for this serial, which should be cluster-safe thanks to a (serial, nonce) PK.
func isFreshCall(payload *iot.DeviceBootstrapPayload) bool {
	now := time.Now()
	timestamp := time.Unix(payload.Timestamp, 0) // Convert Unix epoch (seconds) -> time.Time

	// 1) Strict freshness window (anti "slow replay")
	bootstrapAge := now.Sub(timestamp)
	if bootstrapAge < 0 || 2*time.Minute < bootstrapAge {
		return false
	}

	// TODO change it someday to do something like this:
	//    INSERT ... ON CONFLICT DO NOTHING retourne RowsAffected=0 si (serial,nonce) existe.
	// const insert = `
	// 	INSERT INTO device_nonce_log (serial, nonce, seen_at)
	// 	SELECT @serial, @nonce, SYSUTCDATETIME()
	// 	WHERE NOT EXISTS (
	// 		SELECT 1 FROM device_nonce_log WHERE serial = @serial AND nonce = @nonce
	// 	);

	// res, err := db.ExecContext(ctx, insert, serial, nonce)
	// if err != nil {
	// 	return fmt.Errorf("db insert nonce: %w", err)
	// }
	// aff, _ := res.RowsAffected()
	// if aff == 0 {
	// 	// Conflit → nonce déjà vu pour ce serial
	// 	return ErrReplayNonce
	// }

	// 2) Checking if the exact same request has already been submitted, by using the nonce
	if payload.ExistsInDb() { // (better: single DB UPSERT)
		return false // seen → replay
	}

	payload.InsertInDb() // first time → remember
	return true
}

// Contacting the device provisioning service on behalf of the device
func startProvisioning(serial string, bootstrap *iot.DeviceBootstrap) {
	// TODO better with DB
	bootstrap.Status = iot.BootstrapStatusPENDING
	// update in DB

	// TODO contacting the DPS

	slog.Info("----------------------------------------------------------------")
	slog.Info("CONTACTING THE DPS")
	slog.Info("----------------------------------------------------------------")

	// implement a retry logic
	// TODO what if this fails ?

	// TODO better with DB
	// bootstrap.Status = iot.BootstrapStatusREADY
}
