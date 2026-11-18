package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
)

type BaseUserProfile struct {
	goald.BusinessObject
	User          *User            `json:"user,omitempty"          io:"i*" desc:"The user associated with the profile"`
	ProfileStatus goald.UserStatus `json:"profileStatus,omitempty" io:"in" desc:"The status of the user profile"`
}

func init() {
	model.BaseUserProfile().SetAbstract()
	model.BaseUserProfile().User().SetChildToParent(model.User().Profiles())
	// model.BaseUserProfile().User()
}
