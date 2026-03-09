// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrapRequest specs
type deviceBootstrapRequestSpecs struct {
	g.IBusinessObjectSpecs
	factoryCertPEM *g.StringField
	payloadB64     *g.StringField
	signatureB64   *g.StringField
}

// this is the main way to refer to the DeviceBootstrapRequest specs in the applicative code
func DeviceBootstrapRequest() *deviceBootstrapRequestSpecs {
	return deviceBootstrapRequest
}

// internal variables
var (
	deviceBootstrapRequest     *deviceBootstrapRequestSpecs
	deviceBootstrapRequestOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceBootstrapRequestSpecs() *deviceBootstrapRequestSpecs {
	newSpecs := &deviceBootstrapRequestSpecs{IBusinessObjectSpecs: g.NewBusinessObjectSpecs()}
	newSpecs.factoryCertPEM = g.NewStringField(newSpecs, "FactoryCertPEM", false)
	newSpecs.payloadB64 = g.NewStringField(newSpecs, "PayloadB64", false)
	newSpecs.signatureB64 = g.NewStringField(newSpecs, "SignatureB64", false)

	return newSpecs
}

// making sure the DeviceBootstrapRequest specs exists at app startup
func init() {
	deviceBootstrapRequestOnce.Do(func() {
		deviceBootstrapRequest = newDeviceBootstrapRequestSpecs()
	})

	// this helps dynamically access to the DeviceBootstrapRequest specs
	g.RegisterSpecs("DeviceBootstrapRequest", deviceBootstrapRequest)
}

// accessing all the DeviceBootstrapRequest class' properties and relationships

func (d *deviceBootstrapRequestSpecs) FactoryCertPEM() *g.StringField {
	return d.factoryCertPEM
}

func (d *deviceBootstrapRequestSpecs) PayloadB64() *g.StringField {
	return d.payloadB64
}

func (d *deviceBootstrapRequestSpecs) SignatureB64() *g.StringField {
	return d.signatureB64
}
