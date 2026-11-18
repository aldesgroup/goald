// Generated file, do not edit!
package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
)

// Provides a partial implementation of the User business object
type User struct {
	PermissionOwner
	Email      string               `json:"email,omitempty"      io:"i*" desc:"the %s's email address, which serves as her/his username"`
	FirstName  string               `json:"firstName,omitempty"  io:"i*" desc:"the %s's first name"`
	LastName   string               `json:"lastName,omitempty"   io:"i*" desc:"the %s's last name"`
	ExternalID string               `json:"externalId,omitempty" io:"in" desc:"the unique identifier ('oid' claim) of this %s in its external identity provider, if any"`
	UserStatus goald.UserStatus     `json:"userStatus,omitempty" io:"o*" desc:"the status of the user"`
	MemberOf   []goald.IUserGroup   `json:"memberOf,omitempty"   io:"in" desc:"the list of user groups the %s is a member of"`
	Profiles   []goald.IUserProfile `json:"profiles,omitempty"   io:"in" desc:"the list of profiles associated with the %s"`
}

func init() {
	u := model.User()
	u.SetInDbByName("accessmgt")
	u.SetDescription("Represents a physical user")
	u.FirstName().SetSize(24)
	u.FirstName().SetPersonal()
	u.LastName().SetSize(26)
	u.LastName().SetPersonal()
	u.Email().SetSize(64)
	u.Email().SetUnique()
	u.Email().SetPersonal()
	u.ExternalID().SetSize(64)
	u.ExternalID().SetUnique()
	u.MemberOf().SetSourceToTarget(model.UserGroup().Members())
}
