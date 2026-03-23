// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the Device specs
type deviceSpecs struct {
	g.IBusinessObjectSpecs
	status       *g.EnumField
	statusString *g.StringField
	model        *g.StringField
	serial       *g.StringField
}

// this is the main way to refer to the Device specs in the applicative code
func Device() *deviceSpecs {
	return device
}

// internal variables
var (
	device     *deviceSpecs
	deviceOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceSpecs() *deviceSpecs {
	newSpecs := &deviceSpecs{IBusinessObjectSpecs: g.NewBusinessObjectSpecs()}
	newSpecs.status = g.NewEnumField(newSpecs, "Status", false, "iot.DeviceStatus")
	newSpecs.statusString = g.NewStringField(newSpecs, "StatusString", false)
	newSpecs.model = g.NewStringField(newSpecs, "Model", false)
	newSpecs.serial = g.NewStringField(newSpecs, "Serial", false)

	return newSpecs
}

// making sure the Device specs exists at app startup
func init() {
	deviceOnce.Do(func() {
		device = newDeviceSpecs()
	})

	// this helps dynamically access to the Device specs
	g.RegisterSpecs("Device", device)
}

// accessing all the Device class' properties and relationships

func (d *deviceSpecs) Status() *g.EnumField {
	return d.status
}

func (d *deviceSpecs) StatusString() *g.StringField {
	return d.statusString
}

func (d *deviceSpecs) Model() *g.StringField {
	return d.model
}

func (d *deviceSpecs) Serial() *g.StringField {
	return d.serial
}
