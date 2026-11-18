package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/accessmgt/model"
)

type PermissionOwner struct {
	goald.BusinessObject
	Permissions []goald.IPermission `json:"permissions,omitempty" io:"in" desc:"The permissions owned by this entity"`
}

func init() {
	model.PermissionOwner().SetAbstract()
}
