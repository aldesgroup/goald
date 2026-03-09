package iot

import (
	"github.com/aldesgroup/goald"
	specs "github.com/aldesgroup/goald/_include/_specs"
)

type DeviceBootstrapRequest struct {
	goald.BusinessObject
	FactoryCertPEM string
	PayloadB64     string
	SignatureB64   string
}

func init() {
	specs.DeviceBootstrapRequest().SetNotPersisted()
}
