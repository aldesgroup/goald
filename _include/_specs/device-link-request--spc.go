// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the DeviceLinkRequest specs
type deviceLinkRequestSpecs struct {
	g.IBusinessObjectSpecs
	model            *g.StringField
	serial           *g.StringField
	verificationCode *g.StringField
	userID           *g.IntField
	userFullName     *g.StringField
}

// this is the main way to refer to the DeviceLinkRequest specs in the applicative code
func DeviceLinkRequest() *deviceLinkRequestSpecs {
	return deviceLinkRequest
}

// internal variables
var (
	deviceLinkRequest     *deviceLinkRequestSpecs
	deviceLinkRequestOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newDeviceLinkRequestSpecs() *deviceLinkRequestSpecs {
	newSpecs := &deviceLinkRequestSpecs{IBusinessObjectSpecs: g.NewBusinessObjectSpecs()}
	newSpecs.model = g.NewStringField(newSpecs, "Model", false)
	newSpecs.serial = g.NewStringField(newSpecs, "Serial", false)
	newSpecs.verificationCode = g.NewStringField(newSpecs, "VerificationCode", false)
	newSpecs.userID = g.NewIntField(newSpecs, "UserID", false)
	newSpecs.userFullName = g.NewStringField(newSpecs, "UserFullName", false)

	return newSpecs
}

// making sure the DeviceLinkRequest specs exists at app startup
func init() {
	deviceLinkRequestOnce.Do(func() {
		deviceLinkRequest = newDeviceLinkRequestSpecs()
	})

	// this helps dynamically access to the DeviceLinkRequest specs
	g.RegisterSpecs("DeviceLinkRequest", deviceLinkRequest)
}

// accessing all the DeviceLinkRequest class' properties and relationships

func (d *deviceLinkRequestSpecs) Model() *g.StringField {
	return d.model
}

func (d *deviceLinkRequestSpecs) Serial() *g.StringField {
	return d.serial
}

func (d *deviceLinkRequestSpecs) VerificationCode() *g.StringField {
	return d.verificationCode
}

func (d *deviceLinkRequestSpecs) UserID() *g.IntField {
	return d.userID
}

func (d *deviceLinkRequestSpecs) UserFullName() *g.StringField {
	return d.userFullName
}
