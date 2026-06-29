// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrapRequest model
type DeviceBootstrapRequestModel struct {
	g.IBusinessObjectModel
	factoryCertPem *g.StringField
	payloadB64     *g.StringField
	signatureB64   *g.StringField
}

// this is the main way to refer to the DeviceBootstrapRequest model in the applicative code
func DeviceBootstrapRequest() *DeviceBootstrapRequestModel {
	return deviceBootstrapRequest
}

// internal variables
var (
	deviceBootstrapRequest     *DeviceBootstrapRequestModel
	deviceBootstrapRequestOnce sync.Once
)

// fully describing each of this class' properties & relationships
func NewDeviceBootstrapRequestModel() *DeviceBootstrapRequestModel {
	newModel := &DeviceBootstrapRequestModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.factoryCertPem = g.NewStringField(newModel, "FactoryCertPEM", false)
	newModel.payloadB64 = g.NewStringField(newModel, "PayloadB64", false)
	newModel.signatureB64 = g.NewStringField(newModel, "SignatureB64", false)

	return newModel
}

// making sure the DeviceBootstrapRequest model exists at app startup
func init() {
	deviceBootstrapRequestOnce.Do(func() {
		deviceBootstrapRequest = NewDeviceBootstrapRequestModel()
	})

	// this helps dynamically access to the DeviceBootstrapRequest model
	g.RegisterModel("DeviceBootstrapRequest", deviceBootstrapRequest)
}

// accessing all the DeviceBootstrapRequest class' properties and relationships

func (D *DeviceBootstrapRequestModel) FactoryCertPEM() *g.StringField {
	return D.factoryCertPem
}

func (D *DeviceBootstrapRequestModel) PayloadB64() *g.StringField {
	return D.payloadB64
}

func (D *DeviceBootstrapRequestModel) SignatureB64() *g.StringField {
	return D.signatureB64
}
