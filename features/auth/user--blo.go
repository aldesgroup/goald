package auth

import (
	"github.com/aldesgroup/goald"
)

func (thisUser *User) IsValid(bloCtx goald.BloContext) error {
	if errSuper := thisUser.BusinessObject.IsValid(bloCtx); errSuper != nil {
		return errSuper
	}

	if thisUser.FirstName == "" {
		return goald.Error("First name is required")
	}
	if thisUser.LastName == "" {
		return goald.Error("Last name is required")
	}
	if thisUser.Email == "" {
		return goald.Error("Email is required")
	}
	if thisUser.Password == "" {
		return goald.Error("Password is required")
	}

	return nil
}
