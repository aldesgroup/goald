// ------------------------------------------------------------------------------------------------
// Defining here all the information
// ------------------------------------------------------------------------------------------------
package goald

import (
	"sync/atomic"
	"time"

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
	reqCount        atomic.Int64  // the number of requests handled by this server instance since its startup
}

// Implementing the interface AppContext
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
	*server                   // proxying the server...
	logging.ILogger           // ... but providing this context with it's own logger
	inputBodyBytes  []byte    // keeping track of the incoming request body
	reqNum          int64     // the number of the request being handled by this context
	start           time.Time // the time when the request handling started
}
