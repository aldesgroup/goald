// ------------------------------------------------------------------------------------------------
// Here is the starting point of any Goald app: the initialisation of a server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO add version to BO, isFieldValid (& links ?), patching methods, etc

// ------------------------------------------------------------------------------------------------
// Initialisation
// ------------------------------------------------------------------------------------------------

func InitServer(serverConfig IServerConfig) ServerContext {
	// reading the program's arguments
	var codegen int     // if > 0, the server cannot be started, but code is generated instead
	var srcdir string   // if codegen > 0, this is where to find the go source code
	var confPath string // the path to the config file
	var migrate bool    // if true, then the configured databases are auto-migrated to fit the BOs' persistency requirements

	flag.IntVar(&codegen, "codegen", 0, "if > 0, runs code generation and exits; 1 = objects, 2 = classes")
	flag.StringVar(&srcdir, "srcdir", "go", "where to find all the Go code, from the project's root")
	flag.StringVar(&confPath, "config", "", "the path to the config file")
	flag.BoolVar(&migrate, "migrate", false, "activates the auto-migration of the configured databases")
	flag.Parse()

	// reading the config file
	readAndCheckConfig(confPath, serverConfig)

	// new server
	server := &server{config: serverConfig}

	// init the router, configuring & adding the REST API endpoints
	server.router = httprouter.New()
	server.router.RedirectTrailingSlash = false
	apiPath := server.config.getCommonConfig().HTTP.ApiPath
	for _, endpoint := range restRegistry.endpoints {
		server.router.Handle(endpoint.getMethod(), apiPath+endpoint.getFullPath(), server.handleFor(endpoint))
	}

	// running the app in code generation mode, i.e. no server started here - should only be used by devs
	if codegen > 0 {
		server.runCodeGen(srcdir, codeGenLevel(codegen))
	}

	// performing some checks on the code - but only in dev mode of course
	if server.IsDev() {
		server.runCodeChecks()
	}

	// initialising the DBs
	for _, dbConfig := range serverConfig.getCommonConfig().Databases {
		initAndRegisterDB(dbConfig)
	}

	// migrating the DBs
	if migrate {
		autoMigrateDBs()
	}

	return server
}

// ------------------------------------------------------------------------------------------------
// Starting the server
// ------------------------------------------------------------------------------------------------

func (thisServer *server) Start() {
	// TODO check the configured host / port, etc

	// TODO fill the requestHandler pool

	if len(restRegistry.endpoints) == 0 {
		log.Printf("No endpoint configured, so no starting of the HTTP server!")
	}

	// TODO set router PanicHandler

	// listening to HTTP requests (blocking process) // TODO use config
	if errListen := http.ListenAndServe(":55555", thisServer.router); errListen != nil && errListen != http.ErrServerClosed {
		panicErrf(errListen, "Could not start the server!")
	}
}

func (thisServer *server) handleFor(ep iEndpoint) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		thisServer.ServeEndpoint(ep, w, req, params)
	}
}
