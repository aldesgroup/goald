// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrap specs
type deviceBootstrapSpecs struct {
	g.IBusinessObjectSpecs
	status      *g.EnumField
	retryAfter  *g.IntField
	iotCertPEM  *g.StringField
	iotChainPEM *g.StringField
	deviceID    *g.StringField
	scopeID     *g.StringField
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
	newSpecs.status = g.NewEnumField(newSpecs, "Status", false, "iot.BootstrapStatus")
	newSpecs.retryAfter = g.NewIntField(newSpecs, "RetryAfter", false)
	newSpecs.iotCertPEM = g.NewStringField(newSpecs, "IotCertPEM", false)
	newSpecs.iotChainPEM = g.NewStringField(newSpecs, "IotChainPEM", false)
	newSpecs.deviceID = g.NewStringField(newSpecs, "DeviceID", false)
	newSpecs.scopeID = g.NewStringField(newSpecs, "ScopeID", false)

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

func (d *deviceBootstrapSpecs) Status() *g.EnumField {
	return d.status
}

func (d *deviceBootstrapSpecs) RetryAfter() *g.IntField {
	return d.retryAfter
}

func (d *deviceBootstrapSpecs) IotCertPEM() *g.StringField {
	return d.iotCertPEM
}

func (d *deviceBootstrapSpecs) IotChainPEM() *g.StringField {
	return d.iotChainPEM
}

func (d *deviceBootstrapSpecs) DeviceID() *g.StringField {
	return d.deviceID
}

func (d *deviceBootstrapSpecs) ScopeID() *g.StringField {
	return d.scopeID
}
