package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
	// "github.com/aldesgroup/goald/_include/auth/model"
)

// Provides a partial
type User struct {
	goald.BusinessObject
	Email     string             `json:"email"     io:"i*" desc:"the user's email address, which serves as her/his username"`
	Password  string             `json:"password"  io:"i*" desc:"the user's password"`
	FirstName string             `json:"firstName" io:"i*" desc:"the user's first name"`
	LastName  string             `json:"lastName"  io:"i*" desc:"the user's last name"`
	MemberOf  []goald.IUserGroup `json:"memberOf"  io:"in" desc:"the list of user groups the user is a member of"`
}

func init() {
	u := model.User()
	u.SetAbstract()
	u.FirstName().SetSize(24)
	u.FirstName().SetPersonal()
	u.LastName().SetSize(24)
	u.LastName().SetPersonal()
	u.Email().SetSize(64)
	u.Email().SetUnique()
	u.Email().SetPersonal()
	u.Password().SetSize(64)
	u.Password().SetSecret()
	u.Password().SetPersonal()
	u.MemberOf().SetSourceToTarget(model.UserGroup().Members())
}
