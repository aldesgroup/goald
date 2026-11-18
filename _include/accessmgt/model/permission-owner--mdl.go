// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the PermissionOwner model
type PermissionOwnerModel struct {
	g.IBusinessObjectModel
	permissions *g.Relationship
}

// this is the main way to refer to the PermissionOwner model in the applicative code
func PermissionOwner() *PermissionOwnerModel {
	permissionOwnerOnce.Do(func() {
		permissionOwner = NewPermissionOwnerModel()

		// this helps dynamically access to the PermissionOwner model
		g.RegisterModel("PermissionOwner", permissionOwner)
	})

	return permissionOwner
}

// internal variables
var (
	permissionOwner     *PermissionOwnerModel
	permissionOwnerOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewPermissionOwnerModel() *PermissionOwnerModel {
	thisModel := &PermissionOwnerModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.permissions = g.AddPolyRelationship(thisModel, "PermissionOwner", "Permissions", true)

	return thisModel
}

// making sure the PermissionOwner model exists at app startup
func init() {
	PermissionOwner()
}

// accessing all the PermissionOwner model's properties and relation	ships

func (P *PermissionOwnerModel) Permissions() *g.Relationship {
	return P.permissions
}
