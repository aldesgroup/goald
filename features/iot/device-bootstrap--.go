package iot

import (
	"github.com/aldesgroup/goald"
	specs "github.com/aldesgroup/goald/_include/_specs"
)

type DeviceBootstrap struct {
	goald.BusinessObject
	ScopeID        string
	StatusToRemove string
	// {
	//   "status": "ready|pending|denied",
	//   "signedCert": "pem",
	//   "caChain": "pem",
	//   "scopeId": "string",
	//   "deviceId": "string",
	//   "retryAfter": 5      // secondes (optionnel si pending)
	// }
}

func init() {
	specs.DeviceBootstrap().SetNotPersisted()
}
