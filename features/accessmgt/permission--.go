package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
)

type Permission struct {
	goald.BusinessObject
	Owner goald.IPermissionOwner `json:"owner,omitempty" io:"i*" desc:"The owner of the permission, i.e. a user or user group"`
}

func init() {
	model.Permission().SetDescription("Permissions define what actions can be performed and by whom")
	model.Permission().SetInDbByName("accessmgt")
	model.Permission().Owner().SetChildToParent(model.PermissionOwner().Permissions())
}
