package goald

import (
	"path"
	"runtime/debug"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Generic definition for the business object model sources
// ------------------------------------------------------------------------------------------------

// A Business Object Model Source is an object associated with a specific Business Object Model.
// Business Object Model Sources are responsible for mainly these things:
// - knowing which module and package a business object is from
// - knowing the last modification date of a business object model
// - instantiating new business objects of a given model
type IBusinessObjectModelSource interface {
	// technical properties
	GetName() utils.ModelName          // the name of the model
	getLastBOMod() time.Time           // last modification of the associated Business Object
	getModule() sourceModuleName       // the application or library in which the associated BO is developed
	setModule(module sourceModuleName) // setting the module
	getSrcPath() string                // source path of the associated Business Object
	getPackage() string                // the name of the package the model is from
	isInterface() bool                 // tells if the model is a concrete one, or an interface
	isFromDir(dirName string) bool     // tells if the model is from the given package

	// public methods
	AsInterface() IBusinessObjectModelSource // sets the model as an interface
	NewObject() any                          // a function to instantiate 1 BO corresponding to this entry
	NewSlice() any                           // a function to instantiate an empty slice of BOs corresponding to this entry
}

// ------------------------------------------------------------------------------------------------
// Base implementation and constructor using it
// ------------------------------------------------------------------------------------------------

// An internal struct that should implement IBusinessObjectModelSource
type baseBusinessObjectModelSource struct {
	modelName utils.ModelName
	lastBOMod time.Time
	module    sourceModuleName
	srcPath   string
	intrface  bool
}

func NewBusinessObjectModelSource(srcPath string, modelName utils.ModelName, lastModification string) IBusinessObjectModelSource {
	date, errParse := time.Parse(time.RFC3339, lastModification)
	core.PanicMsgIfErr(errParse, "'%s' has an invalid date format (which is: 2006-01-02 15:04:05)", lastModification)

	return &baseBusinessObjectModelSource{
		modelName: modelName,
		lastBOMod: date,
		srcPath:   srcPath,
	}
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (this *baseBusinessObjectModelSource) GetName() utils.ModelName {
	return this.modelName
}

func (this *baseBusinessObjectModelSource) getLastBOMod() time.Time {
	return this.lastBOMod
}

func (this *baseBusinessObjectModelSource) setModule(module sourceModuleName) {
	this.module = module
}

func (this *baseBusinessObjectModelSource) getModule() sourceModuleName {
	return this.module
}

func (this *baseBusinessObjectModelSource) getSrcPath() string {
	return this.srcPath
}

func (this *baseBusinessObjectModelSource) getPackage() string {
	return path.Base(this.srcPath)
}

func (this *baseBusinessObjectModelSource) isInterface() bool {
	return this.intrface
}

func (this *baseBusinessObjectModelSource) isFromDir(dirName string) bool {
	return this.getModule() == getCurrentSourceModuleName() && this.getPackage() == dirName
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------

func (this *baseBusinessObjectModelSource) AsInterface() IBusinessObjectModelSource {
	this.intrface = true
	return this
}

// NewObject implements [IBusinessObjectModelSource].
func (this *baseBusinessObjectModelSource) NewObject() any {
	panic("unimplemented")
}

// NewSlice implements [IBusinessObjectModelSource].
func (this *baseBusinessObjectModelSource) NewSlice() any {
	panic("unimplemented")
}

// ------------------------------------------------------------------------------------------------
// Utils helping with code modules
// ------------------------------------------------------------------------------------------------

type sourceModuleName string

var (
	currentSourceModule     string
	currentSourceModuleName sourceModuleName
)

// returns this module's path, e.g. "github.com/aldesgroup/goald"
func getCurrentSourceModule() string {
	if currentSourceModule == "" {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			core.PanicMsg("Could not read the build infos!")
		}

		currentSourceModule = bi.Main.Path
	}

	return currentSourceModule
}

// returns this module's name, e.g. "goald"
func getCurrentSourceModuleName() sourceModuleName {
	if currentSourceModuleName == "" {
		currentSourceModuleName = sourceModuleName(path.Base(getCurrentSourceModule()))
	}

	return currentSourceModuleName
}

// ------------------------------------------------------------------------------------------------
// Proxying the source methods to the model
// ------------------------------------------------------------------------------------------------

func (boModel *businessObjectModel) getLastBOMod() time.Time {
	return boModel.source.getLastBOMod()
}
func (boModel *businessObjectModel) getModule() sourceModuleName {
	return boModel.source.getModule()
}
func (boModel *businessObjectModel) setModule(module sourceModuleName) {
	boModel.source.setModule(module)
}
func (boModel *businessObjectModel) getSrcPath() string {
	return boModel.source.getSrcPath()
}
func (boModel *businessObjectModel) getPackage() string {
	return boModel.source.getPackage()
}
func (boModel *businessObjectModel) isInterface() bool {
	return boModel.source.isInterface()
}
func (boModel *businessObjectModel) isFromDir(dirName string) bool {
	return boModel.source.isFromDir(dirName)
}
func (boModel *businessObjectModel) AsInterface() IBusinessObjectModelSource {
	return boModel.source.AsInterface()
}
func (boModel *businessObjectModel) NewObject() any {
	return boModel.source.NewObject()
}
func (boModel *businessObjectModel) NewSlice() any {
	return boModel.source.NewSlice()
}
