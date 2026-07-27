// ------------------------------------------------------------------------------------------------
// Here is the common code for writing business objects
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"

	"github.com/aldesgroup/goald/features/reflection"
)

// ------------------------------------------------------------------------------------------------
// Interface for all the business objects - All the generic functions will rely on this
// ------------------------------------------------------------------------------------------------

type IBusinessObject interface {
	// identification
	GetClassName(thisBO IBusinessObject) className
	getClass(thisBO IBusinessObject) IClass
	GetID() BObjID
	setID(BObjID)
	GetPreID() int
	setPreID(int)

	// business logic
	ChangeBeforeInsert(BloContext) error
	IsValid(BloContext) error
	ChangeAfterInsert(BloContext) error
}

// ------------------------------------------------------------------------------------------------
// Common implementation for business objects - Should be part of any BO's inheritance
// ------------------------------------------------------------------------------------------------

// type BObjID string // probably a UUID here
type BObjID int64 // probably a UUID here

type BusinessObject struct {
	// properties common to all business objects
	ID    BObjID    `json:"id,omitempty"    io:"o*" desc:"The unique identifier of this business object"`
	Class className `json:"class,omitempty" io:"in" desc:"The name of the business object's class, sometimes used to resolve polymorphic relationships"`
	preID int       `json:"-"               io:"o*" desc:"A temporary identifier in Business Objects lists"`

	// technical stuff
	className className
	model     IBusinessObjectModel
	db        *DB
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

// ------------------------------------------------------------------------------------------------
// Special accessors
// ------------------------------------------------------------------------------------------------

func (thisBO *BusinessObject) GetClassName(thisActualBO IBusinessObject) className {
	if thisBO == nil {
		panic("nil business object")
	}

	if thisBO.className == "" {
		// one of the rare cases where we directly use reflection at runtime!
		thisBO.className = className(reflection.TypeNameOf(thisActualBO, true))
	}

	if thisBO.className == "" {
		panic("Could not find class name for business object")
	}

	return thisBO.className
}

func (thisBO *BusinessObject) getClass(thisActualBO IBusinessObject) IClass {
	return classForName(thisBO.GetClassName(thisActualBO), true)
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
