package iot

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/iot/model"
	"github.com/aldesgroup/goald/features/i18n"
)

type DeviceLinkRequest struct {
	goald.BusinessObject
	Model            string            `json:"model"            io:"i*" desc:"the model of the device to link"`
	Serial           string            `json:"serial"           io:"i*" desc:"the device serial number"`
	VerificationCode string            `json:"verificationCode" io:"in" desc:"VerificationCode: NOT STABLE, will probably be removed"`
	UserID           int               `json:"userId"           io:"in" desc:"UserID: NOT STABLE, will probably be removed"`
	UserFullName     string            `json:"userFullName"     io:"in" desc:"UserFullName:NOT STABLE, will probably be removed"`
	Users            []goald.IUser     `json:"users"            io:"i*" desc:"the associated user"`
	MainContact      goald.IUser       `json:"mainContact"      io:"in" desc:"the main contact user"`
	ForWho           goald.IUser       `json:"forWho"           io:"i*" desc:"the user for whom the device is being linked"`
	ENTranslation    *i18n.Translation `json:"enTranslation"    io:"i*" desc:"the English translation of the device model, for better understanding by the provisioning service"`
}

func init() {
	model.DeviceLinkRequest().SetDescription("A request to link a device to a user account, which should be processed by the provisioning service.")
	model.DeviceLinkRequest().SetNotPersisted()
}

////
