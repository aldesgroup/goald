package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
)

type DeviceBootstrapRequest struct {
	goald.BusinessObject
	FactoryCertPEM string `json:"factoryCertPem" io:"i*" desc:"the PEM-encoded public factory certificate of the device, which should be signed by the intermediate factory cert"`
	PayloadB64     string `json:"payloadB64"     io:"i*" desc:"the base64-encoded payload for the bootstrap request"`
	SignatureB64   string `json:"signatureB64"   io:"i*" desc:"the base64-encoded signature of the payload"`
}

func init() {
	model.DeviceBootstrapRequest().SetDescription("A request to bootstrap a device")
	model.DeviceBootstrapRequest().SetNotPersisted()
}
