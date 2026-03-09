// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrap specs
type deviceBootstrapSpecs struct {
	g.IBusinessObjectSpecs
	scopeID        *g.StringField
	statusToRemove *g.StringField
}

// this is the main way to refer to the DeviceBootstrap specs in the applicative code
func DeviceBootstrap() *deviceBootstrapSpecs {
	return deviceBootstrap
}

// internal variables
var (
	deviceBootstrap     *deviceBootstrapSpecs
	deviceBootstrapOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceBootstrapSpecs() *deviceBootstrapSpecs {
	newSpecs := &deviceBootstrapSpecs{IBusinessObjectSpecs: g.NewBusinessObjectSpecs()}
	newSpecs.scopeID = g.NewStringField(newSpecs, "ScopeID", false)
	newSpecs.statusToRemove = g.NewStringField(newSpecs, "StatusToRemove", false)

	return newSpecs
}

// making sure the DeviceBootstrap specs exists at app startup
func init() {
	deviceBootstrapOnce.Do(func() {
		deviceBootstrap = newDeviceBootstrapSpecs()
	})

	// this helps dynamically access to the DeviceBootstrap specs
	g.RegisterSpecs("DeviceBootstrap", deviceBootstrap)
}

// accessing all the DeviceBootstrap class' properties and relationships

func (d *deviceBootstrapSpecs) ScopeID() *g.StringField {
	return d.scopeID
}

func (d *deviceBootstrapSpecs) StatusToRemove() *g.StringField {
	return d.statusToRemove
}
