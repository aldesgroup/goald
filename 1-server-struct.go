// ------------------------------------------------------------------------------------------------
// Defining here all the information
// ------------------------------------------------------------------------------------------------
package goald

import (
	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/logging"
	r "github.com/julienschmidt/httprouter"
)

// ------------------------------------------------------------------------------------------------
// Server & methods
// ------------------------------------------------------------------------------------------------

type server struct {
	logging.ILogger               // being able to log from the server
	instance        string        // identifying this particular server instance
	config          IServerConfig // keeping tracks of the server's configuration
	router          *r.Router     // the HTTP router used to handle the incoming requests
}

// Implementing the interface ServerContext
func (thisServer *server) CustomConfig() ICustomConfig {
	return thisServer.config.CustomConfig()
}

// Shortcut; true if the 'EnvType' config item is "LOCAL"
func (thisServer *server) IsLocal() bool {
	return thisServer.config.base().resolvedEnvType == core.EnvTypeLOCAL
}

// Shortcut; true if the 'EnvType' config item is "SANDBOX"
func (thisServer *server) IsSandbox() bool {
	return thisServer.config.base().resolvedEnvType == core.EnvTypeSANDBOX
}

// ------------------------------------------------------------------------------------------------
// HTTP Request contexts, that form a limited pool
// ------------------------------------------------------------------------------------------------

// an HTTP request context proxies the main server, but also contains the info
// specific to the currently handled HTTP request
type httpRequestContext struct {
	*server               // proxying the server
	targetRefOrID  string // the ID or ref, or whatever property value used to clearly identify a resource
	inputBodyBytes []byte // keeping track of the incoming request body
}

func (thisReqCtx *httpRequestContext) withTargetRefOrID(targetRefOrID string) *httpRequestContext {
	thisReqCtx.targetRefOrID = targetRefOrID
	return thisReqCtx
}
