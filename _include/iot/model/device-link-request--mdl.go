// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
	i18n_model "github.com/aldesgroup/goald/_include/i18n/model"
)

// static, reflect-free access to the definition of the DeviceLinkRequest model
type deviceLinkRequestModel struct {
	g.IBusinessObjectModel
	model            *g.StringField
	serial           *g.StringField
	verificationCode *g.StringField
	userId           *g.IntField
	userFullName     *g.StringField
	users            *g.Relationship
	mainContact      *g.Relationship
	forWho           *g.Relationship
	enTranslation    *g.Relationship
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
	newModel.userId = g.NewIntField(newModel, "UserID", false)
	newModel.userFullName = g.NewStringField(newModel, "UserFullName", false)
	newModel.users = g.NewPolyRelationship(newModel, "Users", true)
	newModel.mainContact = g.NewPolyRelationship(newModel, "MainContact", false)
	newModel.forWho = g.NewPolyRelationship(newModel, "ForWho", false)
	newModel.enTranslation = g.NewRelationship(newModel, "ENTranslation", false, i18n_model.Translation())

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
	return d.userId
}

func (d *deviceLinkRequestModel) UserFullName() *g.StringField {
	return d.userFullName
}

func (d *deviceLinkRequestModel) Users() *g.Relationship {
	return d.users
}

func (d *deviceLinkRequestModel) MainContact() *g.Relationship {
	return d.mainContact
}

func (d *deviceLinkRequestModel) ForWho() *g.Relationship {
	return d.forWho
}

func (d *deviceLinkRequestModel) ENTranslation() *g.Relationship {
	return d.enTranslation
}
