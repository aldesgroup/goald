// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrapPayload specs
type deviceBootstrapPayloadSpecs struct {
	g.IBusinessObjectSpecs
	serial    *g.StringField
	model     *g.StringField
	timestamp *g.BigIntField
	nonce     *g.StringField
}

// this is the main way to refer to the DeviceBootstrapPayload specs in the applicative code
func DeviceBootstrapPayload() *deviceBootstrapPayloadSpecs {
	return deviceBootstrapPayload
}

// internal variables
var (
	deviceBootstrapPayload     *deviceBootstrapPayloadSpecs
	deviceBootstrapPayloadOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceBootstrapPayloadSpecs() *deviceBootstrapPayloadSpecs {
	newSpecs := &deviceBootstrapPayloadSpecs{IBusinessObjectSpecs: g.NewBusinessObjectSpecs()}
	newSpecs.serial = g.NewStringField(newSpecs, "Serial", false)
	newSpecs.model = g.NewStringField(newSpecs, "Model", false)
	newSpecs.timestamp = g.NewBigIntField(newSpecs, "Timestamp", false)
	newSpecs.nonce = g.NewStringField(newSpecs, "Nonce", false)

	return newSpecs
}

// making sure the DeviceBootstrapPayload specs exists at app startup
func init() {
	deviceBootstrapPayloadOnce.Do(func() {
		deviceBootstrapPayload = newDeviceBootstrapPayloadSpecs()
	})

	// this helps dynamically access to the DeviceBootstrapPayload specs
	g.RegisterSpecs("DeviceBootstrapPayload", deviceBootstrapPayload)
}

// accessing all the DeviceBootstrapPayload class' properties and relationships

func (d *deviceBootstrapPayloadSpecs) Serial() *g.StringField {
	return d.serial
}

func (d *deviceBootstrapPayloadSpecs) Model() *g.StringField {
	return d.model
}

func (d *deviceBootstrapPayloadSpecs) Timestamp() *g.BigIntField {
	return d.timestamp
}

func (d *deviceBootstrapPayloadSpecs) Nonce() *g.StringField {
	return d.nonce
}
