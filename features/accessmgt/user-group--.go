package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
	// "github.com/aldesgroup/goald/_include/auth/model"
)

type UserGroup struct {
	goald.BusinessObject
	Name        string        `json:"name"        io:"i*" desc:"the user group's name"`
	Description string        `json:"description" io:"in" desc:"the user group's description"`
	Members     []goald.IUser `json:"members"     io:"in" desc:"the user group's members"`
}

func init() {
	ug := model.UserGroup()
	ug.SetDescription("User groups are used to group users together, for example to assign them the same permissions")
	ug.SetInDbByName("accessmgt") // you have to have an 'accessmgt' database to use this feature, or at least alias your database as 'accessmgt' in the config file
	ug.Name().SetSize(24).SetUnique()
	ug.Description().SetSize(64)

	goald.SetAutoCRUD[*UserGroup](groupAUTH)
}
