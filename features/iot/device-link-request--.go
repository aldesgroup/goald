package iot

import (
	"github.com/aldesgroup/goald"
	specs "github.com/aldesgroup/goald/_include/_specs"
)

type DeviceLinkRequest struct {
	goald.BusinessObject
	Model            string
	Serial           string
	VerificationCode string
	UserID           int    // TODO remove
	UserFullName     string // TODO remove
}

func init() {
	specs.DeviceLinkRequest().SetNotPersisted()
}
