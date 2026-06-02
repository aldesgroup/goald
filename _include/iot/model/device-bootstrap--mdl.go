// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrap model
type deviceBootstrapModel struct {
	g.IBusinessObjectModel
	status      *g.EnumField
	iotCertPem  *g.StringField
	iotChainPem *g.StringField
	deviceId    *g.StringField
	scopeId     *g.StringField
}

// this is the main way to refer to the DeviceBootstrap model in the applicative code
func DeviceBootstrap() *deviceBootstrapModel {
	return deviceBootstrap
}

// internal variables
var (
	deviceBootstrap     *deviceBootstrapModel
	deviceBootstrapOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceBootstrapModel() *deviceBootstrapModel {
	newModel := &deviceBootstrapModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.status = g.NewEnumField(newModel, "Status", false, "iot.BootstrapStatus")
	newModel.iotCertPem = g.NewStringField(newModel, "IotCertPEM", false)
	newModel.iotChainPem = g.NewStringField(newModel, "IotChainPEM", false)
	newModel.deviceId = g.NewStringField(newModel, "DeviceID", false)
	newModel.scopeId = g.NewStringField(newModel, "ScopeID", false)

	return newModel
}

// making sure the DeviceBootstrap model exists at app startup
func init() {
	deviceBootstrapOnce.Do(func() {
		deviceBootstrap = newDeviceBootstrapModel()
	})

	// this helps dynamically access to the DeviceBootstrap model
	g.RegisterModel("DeviceBootstrap", deviceBootstrap)
}

// accessing all the DeviceBootstrap class' properties and relationships

func (d *deviceBootstrapModel) Status() *g.EnumField {
	return d.status
}

func (d *deviceBootstrapModel) IotCertPEM() *g.StringField {
	return d.iotCertPem
}

func (d *deviceBootstrapModel) IotChainPEM() *g.StringField {
	return d.iotChainPem
}

func (d *deviceBootstrapModel) DeviceID() *g.StringField {
	return d.deviceId
}

func (d *deviceBootstrapModel) ScopeID() *g.StringField {
	return d.scopeId
}
