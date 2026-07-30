// ------------------------------------------------------------------------------------------------
// The code here is about registering globally accessible objects.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/reflection"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Registering Business Object Models
// ------------------------------------------------------------------------------------------------

var modelRegistry = struct {
	items map[utils.ModelName]IBusinessObjectModel
	mx    sync.Mutex
}{
	items: map[utils.ModelName]IBusinessObjectModel{},
}

func RegisterModel(name utils.ModelName, model IBusinessObjectModel) {
	modelRegistry.mx.Lock()

	// setting the model name
	model.setName(name)

	// actual registration
	modelRegistry.items[name] = model
	modelRegistry.mx.Unlock()
}

// ------------------------------------------------------------------------------------------------
// Registering Business Object Model Sources
// ------------------------------------------------------------------------------------------------

// The registry for all the app's business object model sources
var sourceRegistry = &struct {
	items map[utils.ModelName]IBusinessObjectModelSource // all the business objects! mapped by the name
	mx    sync.Mutex
}{
	items: map[utils.ModelName]IBusinessObjectModelSource{},
}

type sourceModuleRegitry struct {
	module sourceModuleName
}

// allows to declare a new module where to register business object model sources
func In(module sourceModuleName) *sourceModuleRegitry {
	return &sourceModuleRegitry{module}
}

func (m *sourceModuleRegitry) Register(source IBusinessObjectModelSource) *sourceModuleRegitry {
	sourceRegistry.mx.Lock()
	defer sourceRegistry.mx.Unlock()

	source.setModule(m.module)

	// registering the business object type globally
	sourceRegistry.items[source.getName()] = source

	return m
}

// ------------------------------------------------------------------------------------------------
// Endpoints registry
// ------------------------------------------------------------------------------------------------

var restRegistry = &struct {
	endpoints []iEndpoint
	mx        sync.Mutex
}{}

// registering happens in the "goald" package, hence the private function
func registerEndpoint(ep iEndpoint) iEndpoint {
	restRegistry.mx.Lock()
	restRegistry.endpoints = append(restRegistry.endpoints, ep)
	restRegistry.mx.Unlock()
	return ep
}

// listing all the endpoints, in a sorted manner
func getSortedEndpointList() []iEndpoint {
	restRegistry.mx.Lock()
	defer restRegistry.mx.Unlock()
	sort.Slice(restRegistry.endpoints, func(i, j int) bool {
		epI := restRegistry.endpoints[i]
		epJ := restRegistry.endpoints[j]
		if epI.getPathAsString() != epJ.getPathAsString() {
			return epI.getPathAsString() < epJ.getPathAsString()
		}

		return epI.getMethod() < epJ.getMethod()
	})

	return restRegistry.endpoints
}

// ------------------------------------------------------------------------------------------------
// DB Adapters registry
// ------------------------------------------------------------------------------------------------

var dbAdapterRegistry = &struct {
	dbAdapters map[dbconn.DatabaseType]iDBAdapter
	mx         sync.Mutex
}{
	dbAdapters: map[dbconn.DatabaseType]iDBAdapter{},
}

func RegisterDbAdapter(dbAdapter iDBAdapter) iDBAdapter {
	dbAdapterRegistry.mx.Lock()
	if dbAdapterRegistry.dbAdapters[dbAdapter.DatabaseType()] != nil {
		panic(fmt.Sprintf("There's already a DB adapter registered for database type '%s'", dbAdapter.DatabaseType()))
	}
	dbAdapterRegistry.dbAdapters[dbAdapter.DatabaseType()] = dbAdapter
	dbAdapterRegistry.mx.Unlock()
	return dbAdapter
}

// returning the right adapter for the given database type, or panicking if not found
func getDbAdapter(dbType dbconn.DatabaseType) iDBAdapter {
	dbAdapterRegistry.mx.Lock()
	defer dbAdapterRegistry.mx.Unlock()

	dbAdapter := dbAdapterRegistry.dbAdapters[dbType]
	if dbAdapter == nil {
		panic(fmt.Sprintf("No DB adapter found for database type '%s'. "+
			"You should add this import: _ \"github.com/aldesgroup/goald/features/dbconn/%s\"",
			dbType, getPackageForDbServerType(dbType)))
	}

	return dbAdapter
}

// ------------------------------------------------------------------------------------------------
// DB registry
// ------------------------------------------------------------------------------------------------

var dbRegistry = &struct {
	databases map[dbconn.DbSchemaName]*DB
	mx        sync.Mutex
}{
	databases: map[dbconn.DbSchemaName]*DB{},
}

func GetDB(dbID dbconn.DbSchemaName) *DB {
	dbRegistry.mx.Lock()
	defer dbRegistry.mx.Unlock()

	db := dbRegistry.databases[dbID]
	if db == nil {
		db = &DB{name: dbID}
		dbRegistry.databases[dbID] = db
	}

	return db
}

// ------------------------------------------------------------------------------------------------
// Data loaders
// ------------------------------------------------------------------------------------------------

var dataLoaderRegistry = &struct {
	migrationLoaders map[string]dataLoader // loaders running during a migration phase
	appServerLoaders map[string]dataLoader // loaders running at each app server instance startup
	mx               sync.Mutex
}{
	migrationLoaders: map[string]dataLoader{},
	appServerLoaders: map[string]dataLoader{},
}

func RegisterDataLoader(fn dataLoader, migrationPhase bool) {
	fnName := reflection.GetFnName(fn)
	fnName = fnName[strings.LastIndex(fnName, ".")+1:]
	dataLoaderRegistry.mx.Lock()
	if migrationPhase {
		core.PanicMsgIf(dataLoaderRegistry.migrationLoaders[fnName] != nil, "There's already a migration loader registered for name '%s'", fnName)
		dataLoaderRegistry.migrationLoaders[fnName] = fn
	} else {
		core.PanicMsgIf(dataLoaderRegistry.appServerLoaders[fnName] != nil, "There's already a app server loader registered for name '%s'", fnName)
		dataLoaderRegistry.appServerLoaders[fnName] = fn
	}
	dataLoaderRegistry.mx.Unlock()
}

// ------------------------------------------------------------------------------------------------
// Data Access Objects (DAO) registry
// ------------------------------------------------------------------------------------------------

var daoRegistry = &struct {
	daos map[utils.ModelName]IBusinessObjectDAO
	mx   sync.Mutex
}{
	daos: map[utils.ModelName]IBusinessObjectDAO{},
}

func RegisterDAO(modelName utils.ModelName, dao IBusinessObjectDAO) IBusinessObjectDAO {
	daoRegistry.mx.Lock()
	if daoRegistry.daos[modelName] != nil {
		panic(fmt.Sprintf("There's already a DAO registered for model '%s'", modelName))
	}
	daoRegistry.daos[modelName] = dao
	daoRegistry.mx.Unlock()
	return dao
}

func newDaoFor(modelName utils.ModelName) IBusinessObjectDAO {
	daoRegistry.mx.Lock()
	defer daoRegistry.mx.Unlock()

	dao := daoRegistry.daos[modelName]
	if dao == nil {
		panic(fmt.Sprintf("No DAO found for model '%s'", modelName))
	}

	// a DAO is somehow its own factory
	return dao.NewDAO()
}
