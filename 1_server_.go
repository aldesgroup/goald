// ------------------------------------------------------------------------------------------------
// Here is the starting point of any Goald app: the initialisation of a server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// TODO add version to BO, isFieldValid (& links ?), patching methods, etc

// ------------------------------------------------------------------------------------------------
// Initialisation
// ------------------------------------------------------------------------------------------------

// This function should be called in each Goald-based app
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
	server := &server{config: serverConfig,
		instance: RandomString(3), // TODO remove ?
	}

	// init the logger
	slog.SetLogLoggerLevel(slog.LevelDebug) // TODO configure
	slog.Info(fmt.Sprintf("Instance: %s", server.instance))

	// init the router
	server.initRoutes()

	// running the app in code generation mode, i.e. no server started here - should only be used by devs
	if codegen > 0 {
		server.runCodeGen(srcdir, codeGenLevel(codegen))
	}

	// performing some checks on the code - but only in dev mode of course
	if server.IsDev() {
		server.runCodeChecks()
	}

	// initialising the DBs
	for _, dbConfig := range serverConfig.commonPart().Databases {
		initAndRegisterDB(dbConfig)
	}

	// migrating the DBs
	if migrate {
		autoMigrateDBs()
	}

	return server
}

// ------------------------------------------------------------------------------------------------
// Initialising the routes
// ------------------------------------------------------------------------------------------------

func (thisServer *server) initRoutes() {
	// no HTTP configured? Let's WARN about it
	if thisServer.config.commonPart().HTTP == nil {
		panicf("No \"HTTP\" section configured!")
	}

	// new router
	thisServer.router = httprouter.New()
	thisServer.router.RedirectTrailingSlash = false

	// configuring & adding the REST API endpoints - should we have to serve an API
	apiPath := thisServer.config.commonPart().HTTP.ApiPath
	if apiPath != "" {
		for _, endpoint := range restRegistry.endpoints {
			slog.Info(fmt.Sprintf("Serving: %s %s", endpoint.getMethod(), apiPath+endpoint.getFullPath()))
			thisServer.router.Handle(endpoint.getMethod(), apiPath+endpoint.getFullPath(), thisServer.handleFor(endpoint))
		}
	} else {
		panicf("No path provided for the API!")
	}

	// configuring the static routes
	for _, route := range thisServer.config.commonPart().HTTP.StaticRoutes {
		if fileToServe := route.ServeFile; fileToServe != "" {
			thisServer.router.HandlerFunc(http.MethodGet, route.For, func(w http.ResponseWriter, r *http.Request) { // e.g.: "/"
				slog.Debug(fmt.Sprintf("Serving file %s for %s", fileToServe, r.URL.Path))
				http.ServeFile(w, r, fileToServe) // e.g. serving "webapp/dist/index.html"
			})
		} else {
			path := route.For
			if strings.HasSuffix(path, "*") {
				path += "filepath"
			}

			thisServer.router.HandlerFunc(http.MethodGet, path, func(w http.ResponseWriter, r *http.Request) {
				slog.Debug(fmt.Sprintf("Serving file %s from %s", r.URL.Path, route.ServeDir))
				http.ServeFile(w, r, route.ServeDir+r.URL.Path) // e.g. serving index.html
			})
		}
	}
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

	// listening to HTTP requests (blocking process)
	port := thisServer.config.commonPart().HTTP.Port
	addr := fmt.Sprintf(":%d", port)
	slog.Info(fmt.Sprintf("Serving at: http://localhost:%d/", port))
	if errListen := http.ListenAndServe(addr, thisServer.router); errListen != nil && errListen != http.ErrServerClosed {
		panicErrf(errListen, "Could not start the server!")
	}
}

func (thisServer *server) handleFor(ep iEndpoint) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		thisServer.ServeEndpoint(ep, w, req, params)
	}
}
