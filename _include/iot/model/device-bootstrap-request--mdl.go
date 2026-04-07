// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrapRequest model
type deviceBootstrapRequestModel struct {
	g.IBusinessObjectModel
	factoryCertPEM *g.StringField
	payloadB64     *g.StringField
	signatureB64   *g.StringField
}

// this is the main way to refer to the DeviceBootstrapRequest model in the applicative code
func DeviceBootstrapRequest() *deviceBootstrapRequestModel {
	return deviceBootstrapRequest
}

// internal variables
var (
	deviceBootstrapRequest     *deviceBootstrapRequestModel
	deviceBootstrapRequestOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceBootstrapRequestModel() *deviceBootstrapRequestModel {
	newModel := &deviceBootstrapRequestModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.factoryCertPEM = g.NewStringField(newModel, "FactoryCertPEM", false)
	newModel.payloadB64 = g.NewStringField(newModel, "PayloadB64", false)
	newModel.signatureB64 = g.NewStringField(newModel, "SignatureB64", false)

	return newModel
}

// making sure the DeviceBootstrapRequest model exists at app startup
func init() {
	deviceBootstrapRequestOnce.Do(func() {
		deviceBootstrapRequest = newDeviceBootstrapRequestModel()
	})

	// this helps dynamically access to the DeviceBootstrapRequest model
	g.RegisterModel("DeviceBootstrapRequest", deviceBootstrapRequest)
}

// accessing all the DeviceBootstrapRequest class' properties and relationships

func (d *deviceBootstrapRequestModel) FactoryCertPEM() *g.StringField {
	return d.factoryCertPEM
}

func (d *deviceBootstrapRequestModel) PayloadB64() *g.StringField {
	return d.payloadB64
}

func (d *deviceBootstrapRequestModel) SignatureB64() *g.StringField {
	return d.signatureB64
}
