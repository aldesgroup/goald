// ------------------------------------------------------------------------------------------------
// Here is the common code for writing business objects
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Interface for all the business objects - All the generic functions will rely on this
// ------------------------------------------------------------------------------------------------

type IBusinessObject interface {
	// identification
	GetID() BObjID
	setID(BObjID)
	GetPreID() int
	setPreID(int)

	// business logic
	ChangeBeforeInsert(BloContext) error
	IsValid(BloContext) error
	ChangeAfterInsert(BloContext) error

	// utilities for accessing properties and relationships without using reflection
	ClassName() utils.ClassName                                             // returning the name of the business object's class
	GetValueAsString(string) string                                         // returning a BO's field's value, given the field's name
	SetValueAsString(string, string) error                                  // setting a BO's field's value, given the field's name
	SetRelationshipValue(relName string, value IBusinessObject) error       // setting a single-valued relationship's target, given the relationship's name - without using reflection
	AddRelationshipValue(relName string, value IBusinessObject) error       // appending a target to a multi-valued relationship, given the relationship's name - without using reflection
	ClearRelationshipValue(relName string) error                            // resetting a multi-valued relationship to an empty slice, given the relationship's name - without using reflection
	GetMultipleRelationshipValue(relName string) ([]IBusinessObject, error) // returning the targets of a multi-valued relationship, given the relationship's name - without using reflection
	IsModelValid() error                                                    // checking a business object's general validity - without using reflection

	// technical stuff
	getClass() IClass
	getModel() IBusinessObjectModel
}

// ------------------------------------------------------------------------------------------------
// Common implementation for business objects - Should be part of any BO's inheritance
// ------------------------------------------------------------------------------------------------

// type BObjID string // probably a UUID here
type BObjID int64 // probably a UUID here

type BusinessObject struct {
	// properties common to all business objects
	ID    BObjID          `json:"id,omitempty"    io:"o*" desc:"The unique identifier of this business object"`
	Class utils.ClassName `json:"class,omitempty" io:"in" desc:"The name of the business object's class, sometimes used to resolve polymorphic relationships"`
	preID int             `json:"-"               io:"o*" desc:"A temporary identifier in Business Objects lists"`

	// technical stuff
	class IClass
	db    *DB
}

var _ IBusinessObject = (*BusinessObject)(nil)

// Basic accessors
func (thisBO *BusinessObject) GetID() BObjID      { return thisBO.ID }
func (thisBO *BusinessObject) setID(id BObjID)    { thisBO.ID = id }
func (thisBO *BusinessObject) GetPreID() int      { return thisBO.preID }
func (thisBO *BusinessObject) setPreID(preID int) { thisBO.preID = preID }

// Triggers - default implems
func (thisBO *BusinessObject) ChangeBeforeInsert(BloContext) error { return nil }
func (thisBO *BusinessObject) IsValid(BloContext) error            { return nil }
func (thisBO *BusinessObject) ChangeAfterInsert(BloContext) error  { return nil }

// Utilities - default implems
func (thisBO *BusinessObject) ClassName() utils.ClassName            { panic("unimplemented") }
func (thisBO *BusinessObject) GetValueAsString(string) string        { panic("unimplemented") }
func (thisBO *BusinessObject) SetValueAsString(string, string) error { panic("unimplemented") }
func (thisBO *BusinessObject) SetRelationshipValue(relName string, value IBusinessObject) error {
	panic("unimplemented")
}
func (thisBO *BusinessObject) AddRelationshipValue(relName string, value IBusinessObject) error {
	panic("unimplemented")
}
func (thisBO *BusinessObject) ClearRelationshipValue(relName string) error { panic("unimplemented") }
func (thisBO *BusinessObject) GetMultipleRelationshipValue(relName string) ([]IBusinessObject, error) {
	panic("unimplemented")
}
func (thisBO *BusinessObject) IsModelValid() error { panic("unimplemented") }

// ------------------------------------------------------------------------------------------------
// Special accessors
// ------------------------------------------------------------------------------------------------

func (thisBO *BusinessObject) getClass() IClass {
	if thisBO.class == nil {
		thisBO.class = classFor(thisBO.ClassName(), true)
	}
	return thisBO.class
}

func (thisBO *BusinessObject) getModel() IBusinessObjectModel {
	return thisBO.getClass().getModel()
}

// ------------------------------------------------------------------------------------------------
// Modelling enum types
// ------------------------------------------------------------------------------------------------

// IEnum must be implemented by every enum type
type IEnum interface {
	fmt.Stringer            // each enum value has a default label
	Val() int               // each enum value has an integer value
	Values() map[int]string // each enum has a set of values associated with labels
}

// ------------------------------------------------------------------------------------------------
// Generic definition for the core of each business object class
// ------------------------------------------------------------------------------------------------

// A Class is an object associated with a specific Business Object type,
// and that allows to instantiate new objects of that type, and to get the class' metadata
type IClass interface {
	// technical properties
	getClassName() utils.ClassName  // class of the associated Business Object
	getLastBOMod() time.Time        // last modification of the associated Business Object
	getModule() moduleName          // the application or library in which the associated BO is developed
	setModule(module moduleName)    // setting the module
	getSrcPath() string             // source path of the associated Business Object
	getPackage() string             // the name of the package the class is from
	isInterface() bool              // tells if the class is a concrete one, or an interface
	isFromDir(dirName string) bool  // tells if the class is from the given package
	getModel() IBusinessObjectModel // the model of the associated Business Object

	// public methods
	AsInterface() IClass // sets the class as an interface
	NewObject() any      // a function to instantiate 1 BO corresponding to this entry
	NewSlice() any       // a function to instantiate an empty slice of BOs corresponding to this entry
}

// An internal struct that should implement IClass
type baseClass struct {
	class     utils.ClassName
	lastBOMod time.Time
	module    moduleName
	srcPath   string
	intrface  bool
	model     IBusinessObjectModel
}

func NewClass(srcPath, class, lastModification string) IClass {
	date, errParse := time.Parse(time.RFC3339, lastModification)
	core.PanicMsgIfErr(errParse, "'%s' has an invalid date format (which is: 2006-01-02 15:04:05)", lastModification)

	return &baseClass{
		class:     utils.ClassName(class),
		lastBOMod: date,
		srcPath:   srcPath,
	}
}

func (thisClass *baseClass) getClassName() utils.ClassName {
	return thisClass.class
}

func (thisClass *baseClass) getLastBOMod() time.Time {
	return thisClass.lastBOMod
}

func (thisClass *baseClass) setModule(module moduleName) {
	thisClass.module = module
}

func (thisClass *baseClass) getModule() moduleName {
	return thisClass.module
}

func (thisClass *baseClass) getSrcPath() string {
	return thisClass.srcPath
}

func (thisClass *baseClass) getPackage() string {
	return path.Base(thisClass.srcPath)
}

func (thisClass *baseClass) isInterface() bool {
	return thisClass.intrface
}

func (thisClass *baseClass) isFromDir(dirName string) bool {
	return thisClass.getModule() == getCurrentModuleName() && thisClass.getPackage() == dirName
}

func (thisClass *baseClass) getModel() IBusinessObjectModel {
	if thisClass.model == nil {
		thisClass.model = modelRegistry.items[thisClass.getClassName()]
	}
	return thisClass.model
}

func (thisClass *baseClass) AsInterface() IClass {
	thisClass.intrface = true
	return thisClass
}

// NewObject implements [IClass].
func (thisClass *baseClass) NewObject() any {
	panic("unimplemented")
}

// NewSlice implements [IClass].
func (thisClass *baseClass) NewSlice() any {
	panic("unimplemented")
}
