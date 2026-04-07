// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrapPayload model
type deviceBootstrapPayloadModel struct {
	g.IBusinessObjectModel
	serial    *g.StringField
	model     *g.StringField
	timestamp *g.BigIntField
	nonce     *g.StringField
}

// this is the main way to refer to the DeviceBootstrapPayload model in the applicative code
func DeviceBootstrapPayload() *deviceBootstrapPayloadModel {
	return deviceBootstrapPayload
}

// internal variables
var (
	deviceBootstrapPayload     *deviceBootstrapPayloadModel
	deviceBootstrapPayloadOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceBootstrapPayloadModel() *deviceBootstrapPayloadModel {
	newModel := &deviceBootstrapPayloadModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.serial = g.NewStringField(newModel, "Serial", false)
	newModel.model = g.NewStringField(newModel, "Model", false)
	newModel.timestamp = g.NewBigIntField(newModel, "Timestamp", false)
	newModel.nonce = g.NewStringField(newModel, "Nonce", false)

	return newModel
}

// making sure the DeviceBootstrapPayload model exists at app startup
func init() {
	deviceBootstrapPayloadOnce.Do(func() {
		deviceBootstrapPayload = newDeviceBootstrapPayloadModel()
	})

	// this helps dynamically access to the DeviceBootstrapPayload model
	g.RegisterModel("DeviceBootstrapPayload", deviceBootstrapPayload)
}

// accessing all the DeviceBootstrapPayload class' properties and relationships

func (d *deviceBootstrapPayloadModel) Serial() *g.StringField {
	return d.serial
}

func (d *deviceBootstrapPayloadModel) Model() *g.StringField {
	return d.model
}

func (d *deviceBootstrapPayloadModel) Timestamp() *g.BigIntField {
	return d.timestamp
}

func (d *deviceBootstrapPayloadModel) Nonce() *g.StringField {
	return d.nonce
}
