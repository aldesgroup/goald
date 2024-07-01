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
// The registry for all the app's business objects.
// This helps registering 1 instance of each business object type, which is then used
// by code generation mechanisms to generate the business object classes, using reflection
// ------------------------------------------------------------------------------------------------

type businessObjectEntry struct {
	class      className
	lastMod    time.Time
	module     string
	srcPath    string
	instanceFn func() any // a function to instantiate 1 BO corresponding to this entry
	newSliceFn func() any // a function to instantiate an empty slice of BOs corresponding to this entry
}

var boRegistry = &struct {
	content map[className]*businessObjectEntry // all the business objects! mapped by the name
	mx      sync.Mutex
}{
	content: map[className]*businessObjectEntry{},
}

type moduleBoRegitry struct {
	module string
}

func In(module string) *moduleBoRegitry {
	return &moduleBoRegitry{module}
}

// registering happens in all the applicative packages, gence the public function
func (m *moduleBoRegitry) Register(genOneFn func() any, srcPath string, lastModification string, genSlice func() any) *moduleBoRegitry {
	boRegistry.mx.Lock()
	defer boRegistry.mx.Unlock()

	date, errParse := time.Parse(time.RFC3339, lastModification)
	utils.PanicErrf(errParse, "'%s' has an invalid date format (which is: 2006-01-02 15:04:05)", lastModification)

	entry := &businessObjectEntry{}
	entry.module = m.module
	entry.srcPath = srcPath
	entry.lastMod = date
	entry.instanceFn = genOneFn
	entry.newSliceFn = genSlice
	entry.class = className(utils.TypeNameOf(genOneFn(), true))

	// registering the business object type globally
	boRegistry.content[entry.class] = entry

	return m
}

func NewBO(clsName className) IBusinessObject {
	if boEntry := boRegistry.content[clsName]; boEntry != nil {
		return boEntry.instanceFn().(IBusinessObject)
	}

	panic("The business object registry cannot instantiate an object of class: " + clsName)
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
	loaders map[string]dataLoader
	mx      sync.Mutex
}{
	loaders: map[string]dataLoader{},
}

func RegisterDataLoader(fn dataLoader) {
	fnName := utils.GetFnName(fn)
	fnName = fnName[strings.LastIndex(fnName, ".")+1:]
	dataLoaderRegistry.mx.Lock()
	utils.PanicIff(dataLoaderRegistry.loaders[fnName] != nil, "There's already a loader registered for name '%s'", fnName)
	dataLoaderRegistry.loaders[fnName] = fn
	dataLoaderRegistry.mx.Unlock()
}
