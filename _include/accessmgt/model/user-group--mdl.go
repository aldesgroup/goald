// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the UserGroup model
type UserGroupModel struct {
	g.IBusinessObjectModel
	name        *g.StringField
	description *g.StringField
	members     *g.Relationship
}

// this is the main way to refer to the UserGroup model in the applicative code
func UserGroup() *UserGroupModel {
	return userGroup
}

// internal variables
var (
	userGroup     *UserGroupModel
	userGroupOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewUserGroupModel() *UserGroupModel {
	thisModel := &UserGroupModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.name = g.AddStringField(thisModel, "UserGroup", "Name", false)
	thisModel.description = g.AddStringField(thisModel, "UserGroup", "Description", false)
	thisModel.members = g.AddPolyRelationship(thisModel, "UserGroup", "Members", true)

	return thisModel
}

// making sure the UserGroup model exists at app startup
func init() {
	userGroupOnce.Do(func() {
		userGroup = NewUserGroupModel()
	})

	// this helps dynamically access to the UserGroup model
	g.RegisterModel("UserGroup", userGroup)
}

// accessing all the UserGroup model's properties and relationships

func (U *UserGroupModel) Name() *g.StringField {
	return U.name
}

func (U *UserGroupModel) Description() *g.StringField {
	return U.description
}

func (U *UserGroupModel) Members() *g.Relationship {
	return U.members
}
