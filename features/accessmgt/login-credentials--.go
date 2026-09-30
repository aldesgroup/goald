// Generated file, do not edit!
package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
)

type LoginCredentials struct {
	goald.BusinessObject
	Username string `json:"username,omitempty" io:"i*" desc:"the caller's username, e.g. their email address"`
	Password string `json:"password,omitempty" io:"i*" desc:"the caller's password"`
}

func init() {
	model.LoginCredentials().SetDescription("The username/password pair used to log in and obtain an access token")
	model.LoginCredentials().SetNotPersisted()
	model.LoginCredentials().Password().SetSecret()
}
