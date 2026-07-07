// ------------------------------------------------------------------------------------------------
// The code here is about registering globally accessible objects.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Generic definition for the class core
// ------------------------------------------------------------------------------------------------

// A Class Core is a set of information fields that are common to all the Class objects.
type IClassCore interface {
	getClassName() className                                // class of the associated Business Object
	getLastBOMod() time.Time                                // last modification of the associated Business Object
	getModule() moduleName                                  // the application or library in which the associated BO is developed
	setModule(module moduleName)                            // setting the module
	getSrcPath() string                                     // source path of the associated Business Object
	getPackage() string                                     // the name of the package the class is from
	isInterface() bool                                      // tells if the class is a concrete one, or an interface
	AsInterface() IClassCore                                // sets the class as an interface
	GetValueAsString(IBusinessObject, string) string        // returning a BO's field's value, given the field's name
	SetValueAsString(IBusinessObject, string, string) error // setting a BO's field's value, given the field's name
}

// An internal struct that should implement IClassCore
type classCore struct {
	class     className
	lastBOMod time.Time
	module    moduleName
	srcPath   string
	intrface  bool
}

func NewClassCore(srcPath, class, lastModification string) IClassCore {
	date, errParse := time.Parse(time.RFC3339, lastModification)
	core.PanicMsgIfErr(errParse, "'%s' has an invalid date format (which is: 2006-01-02 15:04:05)", lastModification)

	return &classCore{
		class:     className(class),
		lastBOMod: date,
		srcPath:   srcPath,
	}
}

func (thisCore *classCore) getClassName() className {
	return thisCore.class
}

func (thisCore *classCore) getLastBOMod() time.Time {
	return thisCore.lastBOMod
}

func (thisCore *classCore) setModule(module moduleName) {
	thisCore.module = module
}

func (thisCore *classCore) getModule() moduleName {
	return thisCore.module
}

func (thisCore *classCore) getSrcPath() string {
	return thisCore.srcPath
}

func (thisCore *classCore) getPackage() string {
	return path.Base(thisCore.srcPath)
}

func (thisCore *classCore) isInterface() bool {
	return thisCore.intrface
}

func (thisCore *classCore) AsInterface() IClassCore {
	thisCore.intrface = true
	return thisCore
}

func (thisCore *classCore) GetValueAsString(IBusinessObject, string) string {
	panic("GetValueAsString has to be implemented by a concrete Class__UTILS__ object")
}

func (thisCore *classCore) SetValueAsString(IBusinessObject, string, string) error {
	panic("SetValueAsString has to be implemented by a concrete Class__UTILS__ object")
}

// ------------------------------------------------------------------------------------------------
// Defining and registering classes
// ------------------------------------------------------------------------------------------------

// A Class is an object associated with a specific Business Object type that
// provides automatically STATIC, generated utility methods to:
// - instantiate 1 or a slice of this BO type
// - help serializing / deserializing instances of this BO type
// - quickly perform ORM operations such as Insert(), Select(), Update(), Delete(), etc...
// - ...by containing methods such as GetSelectAllQuery(), GetInsertQuery(), etc
//
// Each Class is loosely coupled to the corresponding BO type through a registry, using
// the BO class as key.
type IClass interface {
	IClassCore

	NewObject() any // a function to instantiate 1 BO corresponding to this entry
	NewSlice() any  // a function to instantiate an empty slice of BOs corresponding to this entry
}

// The registry for all the app's business objects.
// This helps registering 1 instance of each business object type, which is then used
// by code generation mechanisms to generate the business object classes, using reflection
var classRegistry = &struct {
	items map[className]IClass // all the business objects! mapped by the name
	mx    sync.Mutex
}{
	items: map[className]IClass{},
}

type moduleName string

type moduleClassRegitry struct {
	module moduleName
}

// allows to declare a new module where to register Classes
func In(module moduleName) *moduleClassRegitry {
	return &moduleClassRegitry{module}
}

// registering happens in all the applicative packages, gence the public function
func (m *moduleClassRegitry) Register(class IClass) *moduleClassRegitry {
	classRegistry.mx.Lock()
	defer classRegistry.mx.Unlock()

	class.setModule(m.module)

	// registering the business object type globally
	classRegistry.items[class.getClassName()] = class

	return m
}

func classForName(clsName className) IClass {
	if classRegistry.items[clsName] == nil {
		panic(fmt.Sprintf("It looks like no class named '%s' has been registered, "+
			"i.e. its package has probably not been 'included', i.e. imported in the start.go file,"+
			" like this: import _ \"module_full_name/_include/package_name\"", clsName))
	}
	return classRegistry.items[clsName]
}

// 1 Class for 1 Business Object Model
func getClass(model IBusinessObjectModel) IClass {
	return classForName(model.base().name)
}

// ------------------------------------------------------------------------------------------------
// The registry for all the app's business object models
// ------------------------------------------------------------------------------------------------

var modelRegistry = struct {
	items map[className]IBusinessObjectModel
	mx    sync.Mutex
}{
	items: map[className]IBusinessObjectModel{},
}

// registering happens in the "model" package, gence the public function
func RegisterModel(name className, model IBusinessObjectModel) {
	modelRegistry.mx.Lock()

	// setting the class name
	model.base().name = name

	// making sure this class own its fields, including the inherited ones
	for _, field := range model.base().fields {
		field.setOwner(model)
	}

	// making sure this model own its relationships, including the inherited ones
	for _, relationship := range model.base().relationships {
		relationship.setOwner(model)
	}
	modelRegistry.items[name] = model
	modelRegistry.mx.Unlock()
}

func modelForName(clsName className) IBusinessObjectModel {
	// not using the MX for now, but will have to do if there's any possibility for race condition
	return modelRegistry.items[clsName]
}

// ------------------------------------------------------------------------------------------------
// Endpoints registry
// ------------------------------------------------------------------------------------------------

var restRegistry = &struct {
	endpoints []iEndpoint
	mx        sync.Mutex
}{}

// registering happens in the "goald" package, gence the private function
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

// registering happens in the "goald" package, gence the private function
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
		panic(fmt.Sprintf("No DB adapter found for database type '%s'", dbType))
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
		// db = &DB{name: dbID}
		db = &DB{}
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
	fnName := utils.GetFnName(fn)
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
