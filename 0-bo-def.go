// ------------------------------------------------------------------------------------------------
// Here is the common code for writing business objects
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Interface for all the business objects - All the generic functions will rely on this
// ------------------------------------------------------------------------------------------------

type IBusinessObject interface {
	// identification
	GetID() BObjID
	setID(BObjID)
	GetCreation() *time.Time
	setCreation(*time.Time)
	GetPreID() int
	setPreID(int)

	// business logic
	ChangeBeforeInsert(BloContext) error
	IsValid(BloContext) error
	ChangeAfterInsert(BloContext) error

	// utilities for accessing properties and relationships without using reflection
	GetModelName() utils.ModelName                                          // returning the name of the business object's model
	GetValueAsString(string) string                                         // returning a BO's field's value, given the field's name
	SetValueAsString(string, string) error                                  // setting a BO's field's value, given the field's name
	SetRelationshipValue(relName string, value IBusinessObject) error       // setting a single-valued relationship's target, given the relationship's name - without using reflection
	AddRelationshipValue(relName string, value IBusinessObject) error       // appending a target to a multi-valued relationship, given the relationship's name - without using reflection
	ClearRelationshipValue(relName string) error                            // resetting a multi-valued relationship to an empty slice, given the relationship's name - without using reflection
	GetMultipleRelationshipValue(relName string) ([]IBusinessObject, error) // returning the targets of a multi-valued relationship, given the relationship's name - without using reflection
	IsModelValid() error                                                    // checking a business object's general validity - without using reflection
	RemoveCycles()                                                          // removing any cycles from the business object, without using reflection

	// technical stuff
	getModel(this IBusinessObject) IBusinessObjectModel
}

// ------------------------------------------------------------------------------------------------
// Common implementation for business objects - Should be part of any BO's inheritance
// ------------------------------------------------------------------------------------------------

// type BObjID string // probably a UUID here
type BObjID int64 // probably a UUID here

type BusinessObject struct {
	// properties common to all business objects
	ID        BObjID          `json:"id,omitempty"       io:"o*" desc:"The unique identifier of this business object"`
	Creation  *time.Time      `json:"creation,omitempty" io:"o*" desc:"The creation timestamp of this business object"`
	ModelName utils.ModelName `json:"mdl,omitempty"      io:"in" desc:"The name of the business object's model, sometimes used to resolve polymorphic relationships"`
	preID     int             `json:"-"                  io:"o*" desc:"A temporary identifier in Business Objects lists"`

	// technical stuff
	model IBusinessObjectModel
}

var _ IBusinessObject = (*BusinessObject)(nil)

// Basic accessors
func (thisBO *BusinessObject) GetID() BObjID                   { return thisBO.ID }
func (thisBO *BusinessObject) setID(id BObjID)                 { thisBO.ID = id }
func (thisBO *BusinessObject) GetCreation() *time.Time         { return thisBO.Creation }
func (thisBO *BusinessObject) setCreation(creation *time.Time) { thisBO.Creation = creation }
func (thisBO *BusinessObject) GetPreID() int                   { return thisBO.preID }
func (thisBO *BusinessObject) setPreID(preID int)              { thisBO.preID = preID }

// Triggers - default implems
func (thisBO *BusinessObject) ChangeBeforeInsert(BloContext) error { return nil }
func (thisBO *BusinessObject) IsValid(BloContext) error            { return nil }
func (thisBO *BusinessObject) ChangeAfterInsert(BloContext) error  { return nil }

// Utilities - default implems
func (thisBO *BusinessObject) GetModelName() utils.ModelName         { panic("unimplemented") }
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
func (thisBO *BusinessObject) RemoveCycles()       { panic("unimplemented") }

// ------------------------------------------------------------------------------------------------
// Special accessors
// ------------------------------------------------------------------------------------------------

func (thisBO *BusinessObject) getModel(this IBusinessObject) IBusinessObjectModel {
	if thisBO.model == nil {
		thisBO.model = modelFor(this.GetModelName(), true)
	}
	return thisBO.model
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
