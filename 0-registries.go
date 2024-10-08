// ------------------------------------------------------------------------------------------------
// The code here is about registering globally accessible objects.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"strings"
	"sync"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Generic definition for class utils
// ------------------------------------------------------------------------------------------------

// A ClassUtils core is a set of information fields that are common to all the Class Utils objects.
type IClassUtilsCore interface {
	completeCoreWith(className, moduleName)                 // complete the class utils internal state
	getClass() className                                    // class of the associated Business Object
	getLastBOMod() time.Time                                // last modification of the associated Business Object
	getModule() moduleName                                  // the application or library in which the associated BO is developed
	getSrcPath() string                                     // source path of the associated Business Object
	GetValueAsString(IBusinessObject, string) string        // returning a BO's field's value, given the field's name
	SetValueAsString(IBusinessObject, string, string) error // setting a BO's field's value, given the field's name
}

type classUtilsCore struct {
	class     className
	lastBOMod time.Time
	module    moduleName
	srcPath   string
}

func NewClassUtilsCore(srcPath, lastModification string) IClassUtilsCore {
	date, errParse := time.Parse(time.RFC3339, lastModification)
	utils.PanicErrf(errParse, "'%s' has an invalid date format (which is: 2006-01-02 15:04:05)", lastModification)

	return &classUtilsCore{
		lastBOMod: date,
		srcPath:   srcPath,
	}
}

func (thisCore *classUtilsCore) getClass() className {
	return thisCore.class
}

func (thisCore *classUtilsCore) getLastBOMod() time.Time {
	return thisCore.lastBOMod
}

func (thisCore *classUtilsCore) getModule() moduleName {
	return thisCore.module
}

func (thisCore *classUtilsCore) getSrcPath() string {
	return thisCore.srcPath
}

func (thisCore *classUtilsCore) completeCoreWith(class className, module moduleName) {
	thisCore.class = class
	thisCore.module = module
}

func (thisCore *classUtilsCore) GetValueAsString(IBusinessObject, string) string {
	panic("GetValueAsString has to be implemented by a concrete ClassUtils object")
}

func (thisCore *classUtilsCore) SetValueAsString(IBusinessObject, string, string) error {
	panic("SetValueAsString has to be implemented by a concrete ClassUtils object")
}

// ------------------------------------------------------------------------------------------------
// Defining and registering class utils
// ------------------------------------------------------------------------------------------------

// A ClassUtils is an object associated with a specific Business Object type that
// provides automatically genrated utility methods to:
// - instantiate 1 or a slice of this BO type
// - help serializing / deserializing instances of this BO type
// - quickly perform ORM operations such as Insert(), Select(), Update(), Delete(), etc...
// - ...by containing methods such as GetSelectAllQuery(), GetInsertQuery(), etc
//
// Each ClassUtils is loosely coupled to the corresponding BO type through a registry, using
// the BO class as key.
type IClassUtils interface {
	IClassUtilsCore
	NewObject() any // a function to instantiate 1 BO corresponding to this entry
	NewSlice() any  // a function to instantiate an empty slice of BOs corresponding to this entry
}

// The registry for all the app's business objects.
// This helps registering 1 instance of each business object type, which is then used
// by code generation mechanisms to generate the business object classes, using reflection
var classUtilsRegistry = &struct {
	content map[className]IClassUtils // all the business objects! mapped by the name
	mx      sync.Mutex
}{
	content: map[className]IClassUtils{},
}

type moduleName string

type moduleClassUtilsRegitry struct {
	module moduleName
}

// allows to declare a new module where to register ClassUtils
func In(module moduleName) *moduleClassUtilsRegitry {
	return &moduleClassUtilsRegitry{module}
}

// registering happens in all the applicative packages, gence the public function
func (m *moduleClassUtilsRegitry) Register(classUtils IClassUtils) *moduleClassUtilsRegitry {
	classUtilsRegistry.mx.Lock()
	defer classUtilsRegistry.mx.Unlock()

	// clsUtilsCore.module = m.module
	class := className(utils.TypeNameOf(classUtils.NewObject(), true))
	classUtils.completeCoreWith(class, m.module)

	// registering the business object type globally
	classUtilsRegistry.content[class] = classUtils
	return m
}

// 1 ClassUtils for 1 Business Object Class
func getClassUtils(class IBusinessObjectClass) IClassUtils {
	return classUtilsRegistry.content[class.base().name]
}

// ------------------------------------------------------------------------------------------------
// The registry for all the app's business object classes
// ------------------------------------------------------------------------------------------------

var classRegistry = struct {
	classes map[className]IBusinessObjectClass
	mx      sync.Mutex
}{
	classes: map[className]IBusinessObjectClass{},
}

// registering happens in the "class" package, gence the public function
func RegisterClass(name className, class IBusinessObjectClass) {
	classRegistry.mx.Lock()

	// setting the class name
	class.base().name = className(name)

	// making sure this class own its fields, including the inherited ones
	for _, field := range class.base().fields {
		field.setOwner(class)
	}

	// making sure this class own its relationships, including the inherited ones
	for _, relationship := range class.base().relationships {
		relationship.setOwner(class)
	}
	classRegistry.classes[name] = class
	classRegistry.mx.Unlock()
}

func classForName(clsName className) IBusinessObjectClass {
	// not using the MX for now, but will have to do if there's any possibility for race condition
	return classRegistry.classes[clsName]
}

func getAllClasses() map[className]IBusinessObjectClass {
	return classRegistry.classes
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

// ------------------------------------------------------------------------------------------------
// DB registry
// ------------------------------------------------------------------------------------------------

var dbRegistry = &struct {
	databases map[DatabaseID]*DB
	mx        sync.Mutex
}{
	databases: map[DatabaseID]*DB{},
}

func initAndRegisterDB(config *dbConfig) {
	dbRegistry.mx.Lock()
	defer dbRegistry.mx.Unlock()

	// init the instance if needed
	db := dbRegistry.databases[config.DbID]
	if db == nil {
		db = &DB{}
	}

	// init the DB driver
	db.DB, db.adapter = openDB(config)
	db.config = config

	// back into the registry (not needed if already done
	dbRegistry.databases[config.DbID] = db
}

func GetDB(dbID DatabaseID) *DB {
	dbRegistry.mx.Lock()
	defer dbRegistry.mx.Unlock()

	db := dbRegistry.databases[dbID]
	if db == nil {
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
		utils.PanicIff(dataLoaderRegistry.migrationLoaders[fnName] != nil, "There's already a migration loader registered for name '%s'", fnName)
		dataLoaderRegistry.migrationLoaders[fnName] = fn
	} else {
		utils.PanicIff(dataLoaderRegistry.appServerLoaders[fnName] != nil, "There's already a app server loader registered for name '%s'", fnName)
		dataLoaderRegistry.appServerLoaders[fnName] = fn
	}
	dataLoaderRegistry.mx.Unlock()
}
