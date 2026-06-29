// ------------------------------------------------------------------------------------------------
// Here is the starting point of any Goald app: the initialisation of a server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/logging"
	"github.com/julienschmidt/httprouter"
)

// TODO add version to BO, isFieldValid (& links ?), patching methods, etc

// ------------------------------------------------------------------------------------------------
// Initialisation
// ------------------------------------------------------------------------------------------------

// This object contains properties used for code generation
type codegenParams struct {
	srcdir       string // if codegen > 0, this is where to find the go source code
	othersrcdirs string // if codegen > 0, this is where to find secondary go source code directories, that we also watch with aldev
	codegen      int    // if > 0, the server cannot be started, but code is generated instead
	docpath      string // if codegen > 0, the path of the API doc file to generate, i.e. data/api-doc.yaml
	webdir       string // if codegen > 0, this is where to find the web app source code, if any
	nativedir    string // if codegen > 0, this is where to find the native app source code, if any
	regen        bool   // if true and codegen > 0, then all the generated code is regenerated
	bindir       string // if codegen > 0, this is where to find the compilated code
	servers      string // if codegen > 0, the list of the remote server URLs, separated by a pipe, for the API doc generation
}

// This function should be called in each Goald-based app
func NewServer() ServerContext {
	// reading the program's arguments
	var confPath string // the path to the config file
	var migrate bool    // if true, then the configured databases are auto-migrated to fit the BOs' persistency requirements
	var cgParams = &codegenParams{}

	flag.StringVar(&confPath, "config", "", "the path to the config file")
	flag.BoolVar(&migrate, "migrate", false, "activates the auto-migration of the configured databases")
	flag.StringVar(&cgParams.srcdir, "srcdir", "api", "where to find all the Go code, from the project's root")
	flag.StringVar(&cgParams.othersrcdirs, "othersrcdirs", "", "where to find secondary go source code directories, eg. path/to/dir1,dir2,etc")
	flag.IntVar(&cgParams.codegen, "codegen", 0, "if > 0, runs code generation and exits; 1 = objects, 2 = classes")
	flag.StringVar(&cgParams.docpath, "docpath", "", "the path of the API doc file to generate, i.e. data/api-doc.yaml")
	flag.StringVar(&cgParams.webdir, "webdir", "webapp", "where to find all the Web app code, from the project's root")
	flag.StringVar(&cgParams.nativedir, "nativedir", "webapp", "where to find all the Native app code, from the project's root")
	flag.BoolVar(&cgParams.regen, "regen", false, "forces the code regeneration")
	flag.StringVar(&cgParams.bindir, "bindir", "bin", "where to find the compilated code")
	flag.StringVar(&cgParams.servers, "servers", "", "the list of the remote server URLs, separated by a pipe, for the API doc generation; "+
		"eg. sandbox:http://dev.example.com|staging:http://qa.example.com|production:http://prd.example.com")
	flag.Parse()

	// reading the config file
	serverConfig := readAndCheckConfig(confPath)
	loggerConfig := serverConfig.base().Logging

	// new instance ID for the server
	instanceID := core.RandomString(3)

	// new server
	server := &server{
		ILogger:  logging.NewLogger(loggerConfig.resolvedLogLevel, instanceID, loggerConfig.Type, loggerConfig.FNames),
		instance: instanceID,
		config:   serverConfig,
	}
	server.Info("New server")

	// running the app in code generation mode, i.e. no server started here - should only be used by devs
	if cgParams.codegen > 0 {
		server.runCodeGen(cgParams)

		os.Exit(0)
	}

	// initialising the DB servers
	if migrate {
		for _, dbConfig := range serverConfig.base().DBServers {
			getDbAdapter(dbConfig.Type).InitDbServer(server, dbConfig)
		}
	}

	// connecting the DB schemas
	for _, dbConfig := range serverConfig.base().DBServers {
		for _, dbSchema := range dbConfig.Schemas {
			server.connectDbSchema(dbSchema)
		}
	}

	server.Debug("coucou")

	// migrating the DBs + injecting some data into the DBs
	if migrate {
		// making sure the DBs are in sync with the code
		server.autoMigrateDBs()

		// loading some data into the DBs
		server.loadData(true)

		os.Exit(0)
	}

	// init the router
	server.initRoutes()

	// loading some data for each instance of this server
	server.loadData(false)

	return server
}

// ------------------------------------------------------------------------------------------------
// Initialising the routes
// ------------------------------------------------------------------------------------------------

const apiPath = "/rest"

func (thisServer *server) initRoutes() {
	// new router
	thisServer.router = httprouter.New()
	thisServer.router.RedirectTrailingSlash = false

	// serving the API doc
	thisServer.Info("Serving: GET /doc/api")
	thisServer.router.Handle(http.MethodGet, "/doc/api", serveDocForAPI)

	// locally, we also serve the API doc from the root path, for easier access
	if thisServer.IsLocal() {
		thisServer.Info("Serving: GET /")
		thisServer.router.Handle(http.MethodGet, "/", serveDocForAPI)
	}

	// configuring & adding the REST API endpoints - should we have to serve an API
	for _, endpoint := range restRegistry.endpoints {
		thisServer.Info(fmt.Sprintf("Serving: %s", endpoint.getPathAsString()))
		thisServer.router.Handle(endpoint.getMethod(), endpoint.getOperationPath(true), thisServer.handleFor(endpoint))
	}
}

// ------------------------------------------------------------------------------------------------
// Making the server type an HTTP server
// ------------------------------------------------------------------------------------------------
func (thisServer *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// this is a place for potential middlewares
	thisServer.Debug(fmt.Sprintf("%+v", req.Header))

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// TODO do better / probably through the config
	if thisServer.IsLocal() || thisServer.IsSandbox() {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	// w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	// w.Header().Add("Access-Control-Allow-Headers", "Aept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// add protection :https://medium.com/@rahulreza920/go-1-25-is-released-faster-smarter-and-safer-9c97ff8b493d

	// TODO better

	thisServer.router.ServeHTTP(w, req)
}

// ------------------------------------------------------------------------------------------------
// Starting the server
// ------------------------------------------------------------------------------------------------

func (thisServer *server) Start() {
	// TODO check the configured host / port, etc
	core.PanicMsgIf(thisServer.config.base().Port == 0, "No --port provided!")

	// TODO fill the requestHandler pool

	if len(restRegistry.endpoints) == 0 {
		thisServer.Warn("No endpoint configured, so no starting of the HTTP server!")
		return
	}

	// TODO set router PanicHandler

	// listening to HTTP requests (blocking process)
	addr := fmt.Sprintf(":%d", thisServer.config.base().Port)
	thisServer.Info(fmt.Sprintf("Serving at: http://localhost:%d/", thisServer.config.base().Port))
	if errListen := http.ListenAndServe(addr, thisServer); errListen != nil && errListen != http.ErrServerClosed {
		core.PanicMsgIfErr(errListen, "Could not start the server!")
	}

	// TODO shutdown
}

func (thisServer *server) handleFor(ep iEndpoint) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		thisServer.ServeEndpoint(ep, w, req, params)
	}
}
