// ------------------------------------------------------------------------------------------------
// The code here is about registering globally accessible objects.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// The registry for all the app's business objects.
// This helps registering 1 instance of each business object type, which is then used
// by code generation mechanisms to generate the business object classes, using reflection
// ------------------------------------------------------------------------------------------------

type businessObjectEntry struct {
	name     string
	lastMod  time.Time
	bObjType reflect.Type
	module   string
	srcPath  string
}

var boRegistry = &struct {
	content map[string]*businessObjectEntry // all the business objects! mapped by the name
	mx      sync.Mutex
}{
	content: map[string]*businessObjectEntry{},
}

// registering happens in all the applicative packages, gence the public function
func Register(bObj IBusinessObject, module, srcPath string, lastModification string) {
	boRegistry.mx.Lock()
	defer boRegistry.mx.Unlock()

	date, errParse := time.Parse(time.RFC3339, lastModification)
	utils.PanicErrf(errParse, "'%s' has an invalid date format (which is: 2006-01-02 15:04:05)", lastModification)

	entry := &businessObjectEntry{}
	entry.bObjType = reflect.TypeOf(bObj).Elem()
	entry.name = entry.bObjType.Name()
	entry.module = module
	entry.srcPath = srcPath
	entry.lastMod = date

	// registering the business object type globally
	boRegistry.content[entry.name] = entry
}

// ------------------------------------------------------------------------------------------------
// The registry for all the app's business object classes
// ------------------------------------------------------------------------------------------------

var classRegistry = struct {
	classes map[string]IBusinessObjectClass
	mx      sync.Mutex
}{
	classes: map[string]IBusinessObjectClass{},
}

// registering happens in the "class" package, gence the public function
func RegisterClass(name string, class IBusinessObjectClass) {
	classRegistry.mx.Lock()

	// setting the class name
	class.base().className = name

	// making sure this class own its fields, including the inherited ones
	for _, field := range class.base().fields {
		field.setOwner(class)
		// println("- " + field.getName() + " belongs to " + field.ownerClass().base().className)
	}

	// making sure this class own its relationships, including the inherited ones
	for _, relationship := range class.base().relationships {
		relationship.setOwner(class)
	}
	classRegistry.classes[name] = class
	classRegistry.mx.Unlock()
}

func ClassForName(name string) IBusinessObjectClass {
	// not using the MX for now, but will have to do if there's any possibility for race condition
	return classRegistry.classes[name]
}

func getAllClasses() map[string]IBusinessObjectClass {
	return classRegistry.classes
}

// func GetClass[BOTYPE IBusinessObject]() IBusinessObjectClass {
// 	return ClassForName(reflect.TypeOf(new(BOTYPE)).Elem().Elem().Name())
// }

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
	db.DB = openDB(config)
	db.config = config
	db.adapter = getAdapter(config.DbType)

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
	loaders map[string]dataLoader
	mx      sync.Mutex
}{
	loaders: map[string]dataLoader{},
}

func RegisterDataLoader(fn dataLoader) {
	fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	fnName = fnName[strings.LastIndex(fnName, ".")+1:]
	dataLoaderRegistry.mx.Lock()
	utils.PanicIff(dataLoaderRegistry.loaders[fnName] != nil, "There's already a loader registered for name '%s'", fnName)
	dataLoaderRegistry.loaders[fnName] = fn
	dataLoaderRegistry.mx.Unlock()
}
