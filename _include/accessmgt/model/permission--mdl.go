// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the Permission model
type PermissionModel struct {
	g.IBusinessObjectModel
	owner *g.Relationship
}

// this is the main way to refer to the Permission model in the applicative code
func Permission() *PermissionModel {
	permissionOnce.Do(func() {
		permission = NewPermissionModel()

		// this helps dynamically access to the Permission model
		g.RegisterModel("Permission", permission)
	})

	return permission
}

// internal variables
var (
	permission     *PermissionModel
	permissionOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewPermissionModel() *PermissionModel {
	thisModel := &PermissionModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.owner = g.AddPolyRelationship(thisModel, "Permission", "Owner", false)

	return thisModel
}

// making sure the Permission model exists at app startup
func init() {
	Permission()
}

// accessing all the Permission model's properties and relation	ships

func (P *PermissionModel) Owner() *g.Relationship {
	return P.owner
}
