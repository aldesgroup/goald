// ------------------------------------------------------------------------------------------------
// Defining here all the information
// ------------------------------------------------------------------------------------------------
package goald

import r "github.com/julienschmidt/httprouter"

// ------------------------------------------------------------------------------------------------
// Server & methods
// ------------------------------------------------------------------------------------------------

type server struct {
	config IServerConfig
	router *r.Router
}

// Implementing the interface ServerContext
func (thisServer *server) CustomConfig() ICustomConfig {
	return thisServer.config.GetCustomConfig()
}

// Shortcut; true if the 'Env' config item is "dev"
func (thisServer *server) IsDev() bool {
	return thisServer.config.getCommonConfig().envAsType == envTypeDEV
}

// Shortcut; true if the 'Env' config item is "prod"
func (thisServer *server) IsProd() bool {
	return thisServer.config.getCommonConfig().envAsType == envTypePROD
}

// ------------------------------------------------------------------------------------------------
// HTTP Request contexts, that form a limited pool
// ------------------------------------------------------------------------------------------------

// an HTTP request context proxies the main server, but also contains the info
// specific to the currently handled HTTP request
type httpRequestContext struct {
	*server // proxying the server
}
