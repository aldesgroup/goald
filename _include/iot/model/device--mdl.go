// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the Device model
type deviceModel struct {
	g.IBusinessObjectModel
	status          *g.EnumField
	statusString    *g.StringField
	model           *g.StringField
	serial          *g.StringField
	associatedUsers *g.Relationship
}

// this is the main way to refer to the Device model in the applicative code
func Device() *deviceModel {
	return device
}

// internal variables
var (
	device     *deviceModel
	deviceOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceModel() *deviceModel {
	newModel := &deviceModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.status = g.NewEnumField(newModel, "Status", false, "iot.DeviceStatus")
	newModel.statusString = g.NewStringField(newModel, "StatusString", false)
	newModel.model = g.NewStringField(newModel, "Model", false)
	newModel.serial = g.NewStringField(newModel, "Serial", false)
	newModel.associatedUsers = g.NewPolyRelationship(newModel, "AssociatedUsers", true)

	return newModel
}

// making sure the Device model exists at app startup
func init() {
	deviceOnce.Do(func() {
		device = newDeviceModel()
	})

	// this helps dynamically access to the Device model
	g.RegisterModel("Device", device)
}

// accessing all the Device class' properties and relationships

func (d *deviceModel) Status() *g.EnumField {
	return d.status
}

func (d *deviceModel) StatusString() *g.StringField {
	return d.statusString
}

func (d *deviceModel) Model() *g.StringField {
	return d.model
}

func (d *deviceModel) Serial() *g.StringField {
	return d.serial
}

func (d *deviceModel) AssociatedUsers() *g.Relationship {
	return d.associatedUsers
}
