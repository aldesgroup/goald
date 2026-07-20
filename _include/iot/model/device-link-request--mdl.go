// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceLinkRequest model
type DeviceLinkRequestModel struct {
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
func DeviceLinkRequest() *DeviceLinkRequestModel {
	return deviceLinkRequest
}

// internal variables
var (
	deviceLinkRequest     *DeviceLinkRequestModel
	deviceLinkRequestOnce sync.Once
)

// fully describing each of this class' properties & relationships
func NewDeviceLinkRequestModel() *DeviceLinkRequestModel {
	thisModel := &DeviceLinkRequestModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.model = g.AddStringField(thisModel, "DeviceLinkRequest", "Model", false)
	thisModel.serial = g.AddStringField(thisModel, "DeviceLinkRequest", "Serial", false)
	thisModel.verificationCode = g.AddStringField(thisModel, "DeviceLinkRequest", "VerificationCode", false)
	thisModel.userId = g.AddIntField(thisModel, "DeviceLinkRequest", "UserID", false)
	thisModel.userFullName = g.AddStringField(thisModel, "DeviceLinkRequest", "UserFullName", false)
	thisModel.users = g.AddPolyRelationship(thisModel, "DeviceLinkRequest", "Users", true)
	thisModel.mainContact = g.AddPolyRelationship(thisModel, "DeviceLinkRequest", "MainContact", false)
	thisModel.forWho = g.AddPolyRelationship(thisModel, "DeviceLinkRequest", "ForWho", false)
	thisModel.enTranslation = g.AddRelationship(thisModel, "DeviceLinkRequest", "ENTranslation", false, "Translation")

	return thisModel
}

// making sure the DeviceLinkRequest model exists at app startup
func init() {
	deviceLinkRequestOnce.Do(func() {
		deviceLinkRequest = NewDeviceLinkRequestModel()
	})

	// this helps dynamically access to the DeviceLinkRequest model
	g.RegisterModel("DeviceLinkRequest", deviceLinkRequest)
}

// accessing all the DeviceLinkRequest class' properties and relationships

func (D *DeviceLinkRequestModel) Model() *g.StringField {
	return D.model
}

func (D *DeviceLinkRequestModel) Serial() *g.StringField {
	return D.serial
}

func (D *DeviceLinkRequestModel) VerificationCode() *g.StringField {
	return D.verificationCode
}

func (D *DeviceLinkRequestModel) UserID() *g.IntField {
	return D.userId
}

func (D *DeviceLinkRequestModel) UserFullName() *g.StringField {
	return D.userFullName
}

func (D *DeviceLinkRequestModel) Users() *g.Relationship {
	return D.users
}

func (D *DeviceLinkRequestModel) MainContact() *g.Relationship {
	return D.mainContact
}

func (D *DeviceLinkRequestModel) ForWho() *g.Relationship {
	return D.forWho
}

func (D *DeviceLinkRequestModel) ENTranslation() *g.Relationship {
	return D.enTranslation
}
