// Generated file, do not edit!
package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
	"github.com/aldesgroup/goald/features/i18n"
)

type DeviceLinkRequest struct {
	goald.BusinessObject
	Model            string            `json:"model,omitempty"            io:"i*" desc:"the model of the device to link"`
	Serial           string            `json:"serial,omitempty"           io:"i*" desc:"the device serial number"`
	VerificationCode string            `json:"verificationCode,omitempty" io:"in" desc:"VerificationCode: NOT STABLE, will probably be removed"`
	UserID           int               `json:"userId,omitempty"           io:"in" desc:"UserID: NOT STABLE, will probably be removed"`
	UserFullName     string            `json:"userFullName,omitempty"     io:"in" desc:"UserFullName:NOT STABLE, will probably be removed"`
	Users            []goald.IUser     `json:"users,omitempty"            io:"i*" desc:"the associated user"`
	MainContact      goald.IUser       `json:"mainContact,omitempty"      io:"in" desc:"the main contact user"`
	ForWho           goald.IUser       `json:"forWho,omitempty"           io:"i*" desc:"the user for whom the device is being linked"`
	ENTranslation    *i18n.Translation `json:"enTranslation,omitempty"    io:"i*" desc:"the English translation of the device model, for better understanding by the provisioning service"`
}

func init() {
	model.DeviceLinkRequest().SetDescription("A request to link a device to a user account, which should be processed by the provisioning service.")
	model.DeviceLinkRequest().SetNotPersisted()
}
