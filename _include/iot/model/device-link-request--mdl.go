// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceLinkRequest model
type deviceLinkRequestModel struct {
	g.IBusinessObjectModel
	model            *g.StringField
	serial           *g.StringField
	verificationCode *g.StringField
	userID           *g.IntField
	userFullName     *g.StringField
}

// this is the main way to refer to the DeviceLinkRequest model in the applicative code
func DeviceLinkRequest() *deviceLinkRequestModel {
	return deviceLinkRequest
}

// internal variables
var (
	deviceLinkRequest     *deviceLinkRequestModel
	deviceLinkRequestOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceLinkRequestModel() *deviceLinkRequestModel {
	newModel := &deviceLinkRequestModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.model = g.NewStringField(newModel, "Model", false)
	newModel.serial = g.NewStringField(newModel, "Serial", false)
	newModel.verificationCode = g.NewStringField(newModel, "VerificationCode", false)
	newModel.userID = g.NewIntField(newModel, "UserID", false)
	newModel.userFullName = g.NewStringField(newModel, "UserFullName", false)

	return newModel
}

// making sure the DeviceLinkRequest model exists at app startup
func init() {
	deviceLinkRequestOnce.Do(func() {
		deviceLinkRequest = newDeviceLinkRequestModel()
	})

	// this helps dynamically access to the DeviceLinkRequest model
	g.RegisterModel("DeviceLinkRequest", deviceLinkRequest)
}

// accessing all the DeviceLinkRequest class' properties and relationships

func (d *deviceLinkRequestModel) Model() *g.StringField {
	return d.model
}

func (d *deviceLinkRequestModel) Serial() *g.StringField {
	return d.serial
}

func (d *deviceLinkRequestModel) VerificationCode() *g.StringField {
	return d.verificationCode
}

func (d *deviceLinkRequestModel) UserID() *g.IntField {
	return d.userID
}

func (d *deviceLinkRequestModel) UserFullName() *g.StringField {
	return d.userFullName
}
