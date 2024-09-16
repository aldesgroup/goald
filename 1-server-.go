// ------------------------------------------------------------------------------------------------
// Here is the starting point of any Goald app: the initialisation of a server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aldesgroup/goald/features/utils"
	"github.com/julienschmidt/httprouter"
)

// TODO add version to BO, isFieldValid (& links ?), patching methods, etc

// ------------------------------------------------------------------------------------------------
// Initialisation
// ------------------------------------------------------------------------------------------------

// This function should be called in each Goald-based app
func NewServer() ServerContext {
	// reading the program's arguments
	var confPath string // the path to the config file
	var srcdir string   // if codegen > 0, this is where to find the go source code
	var migrate bool    // if true, then the configured databases are auto-migrated to fit the BOs' persistency requirements
	var codegen int     // if > 0, the server cannot be started, but code is generated instead
	var webdir string   // if codegen > 0, this is where to find the web app source code, if any
	var regen bool      // if true, then all the generated code is regenerated

	flag.StringVar(&confPath, "config", "", "the path to the config file")
	flag.StringVar(&srcdir, "srcdir", "api", "where to find all the Go code, from the project's root")
	flag.BoolVar(&migrate, "migrate", false, "activates the auto-migration of the configured databases")
	flag.IntVar(&codegen, "codegen", 0, "if > 0, runs code generation and exits; 1 = objects, 2 = classes")
	flag.StringVar(&webdir, "webdir", "webapp", "where to find all the Web app code, from the project's root")
	flag.BoolVar(&regen, "regen", false, "forces the code regeneration")
	flag.Parse()

	// reading the config file
	serverConfig := readAndCheckConfig(confPath)

	// new server
	server := &server{
		config:   serverConfig,
		instance: utils.RandomString(3), // TODO remove ?
	}

	// init the logger
	slog.SetLogLoggerLevel(slog.LevelDebug) // TODO configure

	// running the app in code generation mode, i.e. no server started here - should only be used by devs
	if codegen > 0 {
		server.runCodeGen(srcdir, codeGenLevel(codegen), webdir, regen)
	}

	// performing some checks on the code - but only in dev mode of course
	if server.IsDev() {
		server.runCodeChecks()
	}

	// init the router
	server.initRoutes()

	// initialising the DBs
	for _, dbConfig := range serverConfig.commonPart().Databases {
		initAndRegisterDB(dbConfig)
	}

	// bit of logging // TODO remove
	slog.Info(fmt.Sprintf("Instance: %s", server.instance))

	// migrating the DBs + injecting some data into the DBs
	if migrate {
		// making sure the DBs are in sync with the code
		autoMigrateDBs()

		// loading some data into the DBs
		server.loadData(true)

		os.Exit(0)
	}

	// loading some data for each instance of this server
	server.loadData(false)

	return server
}

// ------------------------------------------------------------------------------------------------
// Initialising the routes
// ------------------------------------------------------------------------------------------------

func (thisServer *server) initRoutes() {
	// no HTTP configured? Let's WARN about it
	if thisServer.config.commonPart().HTTP == nil {
		utils.Panicf("No \"HTTP\" section configured!")
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
		utils.Panicf("No path provided for the API!")
	}

	// configuring the static routes TODO not used for now
	// for _, route := range thisServer.config.commonPart().HTTP.StaticRoutes {
	// 	if fileToServe := route.ServeFile; fileToServe != "" {
	// 		thisServer.router.HandlerFunc(http.MethodGet, route.For, func(w http.ResponseWriter, r *http.Request) { // e.g.: "/"
	// 			slog.Debug(fmt.Sprintf("Serving file %s for %s", fileToServe, r.URL.Path))
	// 			http.ServeFile(w, r, fileToServe) // e.g. serving "webapp/dist/index.html"
	// 		})
	// 	} else {
	// 		path := route.For
	// 		if strings.HasSuffix(path, "*") {
	// 			path += "filepath"
	// 		}

	// 		thisServer.router.HandlerFunc(http.MethodGet, path, func(w http.ResponseWriter, r *http.Request) {
	// 			slog.Debug(fmt.Sprintf("Serving file %s from %s", r.URL.Path, route.ServeDir))
	// 			http.ServeFile(w, r, route.ServeDir+r.URL.Path) // e.g. serving index.html
	// 		})
	// 	}
	// }
}

// ------------------------------------------------------------------------------------------------
// Making the server type an HTTP server
// ------------------------------------------------------------------------------------------------
func (thisServer *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// this is a place for potential middlewares
	slog.Debug(fmt.Sprintf("%+v", req.Header))
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// TODO do better
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Add("Access-Control-Allow-Headers", "Aept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// TODO better

	thisServer.router.ServeHTTP(w, req)
}

// ------------------------------------------------------------------------------------------------
// Starting the server
// ------------------------------------------------------------------------------------------------

func (thisServer *server) Start() {
	// TODO check the configured host / port, etc

	// TODO fill the requestHandler pool

	if len(restRegistry.endpoints) == 0 {
		slog.Warn("No endpoint configured, so no starting of the HTTP server!")
	}

	// TODO set router PanicHandler

	// listening to HTTP requests (blocking process)
	port := thisServer.config.commonPart().HTTP.Port
	addr := fmt.Sprintf(":%d", port)
	slog.Info(fmt.Sprintf("Serving at: http://localhost:%d/", port))
	if errListen := http.ListenAndServe(addr, thisServer); errListen != nil && errListen != http.ErrServerClosed {
		utils.PanicErrf(errListen, "Could not start the server!")
	}

	// TODO shutdown
}

func (thisServer *server) handleFor(ep iEndpoint) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		thisServer.ServeEndpoint(ep, w, req, params)
	}
}
