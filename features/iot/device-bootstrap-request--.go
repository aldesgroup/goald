package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
)

type DeviceBootstrapRequest struct {
	goald.BusinessObject
	FactoryCertPEM string
	PayloadB64     string
	SignatureB64   string
}

func init() {
	model.DeviceBootstrapRequest().SetNotPersisted()
}
