// ------------------------------------------------------------------------------------------------
// Here are defined the contexts passed in the different layers of our multi-tier architecture:
// - (no context) for the business object layer (*--.go files)
// - BloContext for the Business LOgic code (used in *--blo.go files)
// - DaoContext for the Data Access Objects (used in *--dao.go files)
// - WebContext for the web endpoints code (used in *--web.go files)
// ------------------------------------------------------------------------------------------------
package goald

import (
	"github.com/aldesgroup/goald/features/logging"
)

// ------------------------------------------------------------------------------------------------
// AppContext contains the minimal info set that should be accessible in all the layers of the app
// ------------------------------------------------------------------------------------------------
type AppContext interface {
	logging.ILogger              // we should be able to log from anywhere
	CustomConfig() ICustomConfig // returns the app's custom part of the config
}

// ------------------------------------------------------------------------------------------------
// restContext is used in the context of handling with a REST resource (single or plural)
// ------------------------------------------------------------------------------------------------

type restContext interface { // TODO keep ?
	AppContext
	getResourceModel() IBusinessObjectModel // the business object model of the resource being requested
}

// ------------------------------------------------------------------------------------------------
// ServerContext is a particular Business Logic Context used at app startup
// Implemented by the `server` struct
// ------------------------------------------------------------------------------------------------
type ServerContext interface {
	BloContext
	Start()
}
