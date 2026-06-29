package auth

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/auth/model"
)

// Provides a partial
type User struct {
	goald.BusinessObject
	Email     string `json:"email"     io:"i*" desc:"the user's email address, which serves as her/his username"`
	Password  string `json:"password"  io:"i*" desc:"the user's password"`
	FirstName string `json:"firstName" io:"i*" desc:"the user's first name"`
	LastName  string `json:"lastName"  io:"i*" desc:"the user's last name"`
}

func init() {
	model.User().SetAbstract()
}
