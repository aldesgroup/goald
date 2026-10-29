package accessmgt

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/hstatus"
)

func init() {
	goald.GetOne(handleGetCurrentUser, nil).
		Label("Get the current user").
		Description("Retrieve the user information from the currently connected user")
}

func handleGetCurrentUser(webCtx goald.WebContext) (*User, hstatus.Code, string) {
	return nil, hstatus.NotFound, "current user not found"
}
