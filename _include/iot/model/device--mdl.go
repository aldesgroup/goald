// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the Device model
type DeviceModel struct {
	g.IBusinessObjectModel
	status          *g.EnumField
	statusString    *g.StringField
	model           *g.StringField
	serial          *g.StringField
	associatedUsers *g.Relationship
}

// this is the main way to refer to the Device model in the applicative code
func Device() *DeviceModel {
	return device
}

// internal variables
var (
	device     *DeviceModel
	deviceOnce sync.Once
)

// fully describing each of this class' properties & relationships
func NewDeviceModel() *DeviceModel {
	thisModel := &DeviceModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.status = g.AddEnumField(thisModel, "Device", "Status", false, "iot.DeviceStatus")
	thisModel.statusString = g.AddStringField(thisModel, "Device", "StatusString", false)
	thisModel.model = g.AddStringField(thisModel, "Device", "Model", false)
	thisModel.serial = g.AddStringField(thisModel, "Device", "Serial", false)
	thisModel.associatedUsers = g.AddPolyRelationship(thisModel, "Device", "AssociatedUsers", true)

	return thisModel
}

// making sure the Device model exists at app startup
func init() {
	deviceOnce.Do(func() {
		device = NewDeviceModel()
	})

	// this helps dynamically access to the Device model
	g.RegisterModel("Device", device)
}

// accessing all the Device class' properties and relationships

func (D *DeviceModel) Status() *g.EnumField {
	return D.status
}

func (D *DeviceModel) StatusString() *g.StringField {
	return D.statusString
}

func (D *DeviceModel) Model() *g.StringField {
	return D.model
}

func (D *DeviceModel) Serial() *g.StringField {
	return D.serial
}

func (D *DeviceModel) AssociatedUsers() *g.Relationship {
	return D.associatedUsers
}
