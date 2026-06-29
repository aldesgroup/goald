// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceBootstrap model
type DeviceBootstrapModel struct {
	g.IBusinessObjectModel
	status      *g.EnumField
	iotCertPem  *g.StringField
	iotChainPem *g.StringField
	deviceId    *g.StringField
	scopeId     *g.StringField
}

// this is the main way to refer to the DeviceBootstrap model in the applicative code
func DeviceBootstrap() *DeviceBootstrapModel {
	return deviceBootstrap
}

// internal variables
var (
	deviceBootstrap     *DeviceBootstrapModel
	deviceBootstrapOnce sync.Once
)

// fully describing each of this class' properties & relationships
func NewDeviceBootstrapModel() *DeviceBootstrapModel {
	newModel := &DeviceBootstrapModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
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
		deviceBootstrap = NewDeviceBootstrapModel()
	})

	// this helps dynamically access to the DeviceBootstrap model
	g.RegisterModel("DeviceBootstrap", deviceBootstrap)
}

// accessing all the DeviceBootstrap class' properties and relationships

func (D *DeviceBootstrapModel) Status() *g.EnumField {
	return D.status
}

func (D *DeviceBootstrapModel) IotCertPEM() *g.StringField {
	return D.iotCertPem
}

func (D *DeviceBootstrapModel) IotChainPEM() *g.StringField {
	return D.iotChainPem
}

func (D *DeviceBootstrapModel) DeviceID() *g.StringField {
	return D.deviceId
}

func (D *DeviceBootstrapModel) ScopeID() *g.StringField {
	return D.scopeId
}
