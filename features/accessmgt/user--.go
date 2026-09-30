// Generated file, do not edit!
package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
)

// Provides a partial implementation of the User business object
type User struct {
	goald.BusinessObject
	Email      string             `json:"email,omitempty"      io:"i*" desc:"the %s's email address, which serves as her/his username"`
	Password   string             `json:"password,omitempty"   io:"in" desc:"the %s's password, for local accounts not backed by an external identity provider"`
	FirstName  string             `json:"firstName,omitempty"  io:"i*" desc:"the %s's first name"`
	LastName   string             `json:"lastName,omitempty"   io:"i*" desc:"the %s's last name"`
	ExternalID string             `json:"externalId,omitempty" io:"in" desc:"the unique identifier ('oid' claim) of this %s in its external identity provider, if any"`
	Realm      string             `json:"realm,omitempty"      io:"in" desc:"which population this %s belongs to: 'colleague', 'customer', or 'local'"`
	MemberOf   []goald.IUserGroup `json:"memberOf,omitempty"   io:"in" desc:"the list of user groups the %s is a member of"`
}

func init() {
	u := model.User()
	u.SetAbstract()
	u.FirstName().SetSize(24)
	u.FirstName().SetPersonal()
	u.LastName().SetSize(26)
	u.LastName().SetPersonal()
	u.Email().SetSize(64)
	u.Email().SetUnique()
	u.Email().SetPersonal()
	u.Password().SetSize(64)
	u.Password().SetSecret()
	u.Password().SetPersonal()
	u.ExternalID().SetSize(64)
	u.ExternalID().SetPersonal()
	u.Realm().SetSize(16)
	u.MemberOf().SetSourceToTarget(model.UserGroup().Members())
}
