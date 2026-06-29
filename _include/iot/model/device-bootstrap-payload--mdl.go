// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrapPayload model
type DeviceBootstrapPayloadModel struct {
	g.IBusinessObjectModel
	serial    *g.StringField
	model     *g.StringField
	timestamp *g.BigIntField
	nonce     *g.StringField
}

// this is the main way to refer to the DeviceBootstrapPayload model in the applicative code
func DeviceBootstrapPayload() *DeviceBootstrapPayloadModel {
	return deviceBootstrapPayload
}

// internal variables
var (
	deviceBootstrapPayload     *DeviceBootstrapPayloadModel
	deviceBootstrapPayloadOnce sync.Once
)

// fully describing each of this class' properties & relationships
func NewDeviceBootstrapPayloadModel() *DeviceBootstrapPayloadModel {
	newModel := &DeviceBootstrapPayloadModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.serial = g.NewStringField(newModel, "Serial", false)
	newModel.model = g.NewStringField(newModel, "Model", false)
	newModel.timestamp = g.NewBigIntField(newModel, "Timestamp", false)
	newModel.nonce = g.NewStringField(newModel, "Nonce", false)

	return newModel
}

// making sure the DeviceBootstrapPayload model exists at app startup
func init() {
	deviceBootstrapPayloadOnce.Do(func() {
		deviceBootstrapPayload = NewDeviceBootstrapPayloadModel()
	})

	// this helps dynamically access to the DeviceBootstrapPayload model
	g.RegisterModel("DeviceBootstrapPayload", deviceBootstrapPayload)
}

// accessing all the DeviceBootstrapPayload class' properties and relationships

func (D *DeviceBootstrapPayloadModel) Serial() *g.StringField {
	return D.serial
}

func (D *DeviceBootstrapPayloadModel) Model() *g.StringField {
	return D.model
}

func (D *DeviceBootstrapPayloadModel) Timestamp() *g.BigIntField {
	return D.timestamp
}

func (D *DeviceBootstrapPayloadModel) Nonce() *g.StringField {
	return D.nonce
}
