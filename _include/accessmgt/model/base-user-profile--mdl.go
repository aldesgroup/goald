// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the BaseUserProfile model
type BaseUserProfileModel struct {
	g.IBusinessObjectModel
	user          *g.Relationship
	profileStatus *g.EnumField
}

// this is the main way to refer to the BaseUserProfile model in the applicative code
func BaseUserProfile() *BaseUserProfileModel {
	baseUserProfileOnce.Do(func() {
		baseUserProfile = NewBaseUserProfileModel()

		// this helps dynamically access to the BaseUserProfile model
		g.RegisterModel("BaseUserProfile", baseUserProfile)
	})

	return baseUserProfile
}

// internal variables
var (
	baseUserProfile     *BaseUserProfileModel
	baseUserProfileOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewBaseUserProfileModel() *BaseUserProfileModel {
	thisModel := &BaseUserProfileModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.user = g.AddRelationship(thisModel, "BaseUserProfile", "User", false, "User")
	thisModel.profileStatus = g.AddEnumField(thisModel, "BaseUserProfile", "ProfileStatus", false, "goald.UserStatus")

	return thisModel
}

// making sure the BaseUserProfile model exists at app startup
func init() {
	BaseUserProfile()
}

// accessing all the BaseUserProfile model's properties and relation	ships

func (B *BaseUserProfileModel) User() *g.Relationship {
	return B.user
}

func (B *BaseUserProfileModel) ProfileStatus() *g.EnumField {
	return B.profileStatus
}
