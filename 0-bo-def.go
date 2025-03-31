// ------------------------------------------------------------------------------------------------
// Here is the common code for writing business objects
// ------------------------------------------------------------------------------------------------
package goald

import "fmt"

// ------------------------------------------------------------------------------------------------
// Interface for all the business objects - All the generic functions will rely on this
// ------------------------------------------------------------------------------------------------

type IBusinessObject interface {
	// identification
	Specs() IBusinessObjectSpecs
	getClassName() className
	setClassName(className)
	GetID() BObjID
	setID(int)

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
	specs     IBusinessObjectSpecs
	className className
	ID        BObjID `json:",omitempty"`
}

var _ IBusinessObject = (*BusinessObject)(nil)

func (thisBO *BusinessObject) Specs() IBusinessObjectSpecs {
	if thisBO.specs == nil {
		thisBO.specs = specsForName(thisBO.className)
	}

	if thisBO.specs == nil {
		panic("unknown class for a business object!")
	}

	return thisBO.specs
}

/* default implementations */
func (thisBO *BusinessObject) getClassName() className             { return thisBO.className }
func (thisBO *BusinessObject) setClassName(cn className)           { thisBO.className = cn }
func (thisBO *BusinessObject) GetID() BObjID                       { return thisBO.ID }
func (thisBO *BusinessObject) setID(id int)                        { thisBO.ID = BObjID(id) }
func (thisBO *BusinessObject) ChangeBeforeInsert(BloContext) error { return nil }
func (thisBO *BusinessObject) IsValid(BloContext) error            { return nil }
func (thisBO *BusinessObject) ChangeAfterInsert(BloContext) error  { return nil }

// ------------------------------------------------------------------------------------------------
// Modelling enum types
// ------------------------------------------------------------------------------------------------

// IEnum must be implemented by every enum type
type IEnum interface {
	fmt.Stringer // each enum value has a default label
	Val() int
	Values() map[int]string
}
