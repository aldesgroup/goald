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

// fully describing each of this model's properties & relationships
func NewDeviceBootstrapPayloadModel() *DeviceBootstrapPayloadModel {
	thisModel := &DeviceBootstrapPayloadModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.serial = g.AddStringField(thisModel, "DeviceBootstrapPayload", "Serial", false)
	thisModel.model = g.AddStringField(thisModel, "DeviceBootstrapPayload", "Model", false)
	thisModel.timestamp = g.AddBigIntField(thisModel, "DeviceBootstrapPayload", "Timestamp", false)
	thisModel.nonce = g.AddStringField(thisModel, "DeviceBootstrapPayload", "Nonce", false)

	return thisModel
}

// making sure the DeviceBootstrapPayload model exists at app startup
func init() {
	deviceBootstrapPayloadOnce.Do(func() {
		deviceBootstrapPayload = NewDeviceBootstrapPayloadModel()
	})

	// this helps dynamically access to the DeviceBootstrapPayload model
	g.RegisterModel("DeviceBootstrapPayload", deviceBootstrapPayload)
}

// accessing all the DeviceBootstrapPayload model's properties and relationships

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
