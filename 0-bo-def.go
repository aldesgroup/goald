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
	GetID() BObjID                    // returning the business object's ID
	setID(BObjID)                     // setting the business object's ID
	GetCreation() *time.Time          // returning the business object's creation timestamp
	setCreation(*time.Time)           // setting the business object's creation timestamp
	GetModification() *time.Time      // returning the business object's modification timestamp
	setModification(*time.Time)       // setting the business object's modification timestamp
	GetPreID() int                    // returning the business object's pre-insert ID
	setPreID(int)                     // setting the business object's pre-insert ID
	getKey() boKey                    // returning the business object's key
	setKey(boKey)                     // setting the business object's key
	SetParent(parent IBusinessObject) // setting the business object's parent, if it has one

	// business logic
	ChangeBeforeInsert(BloContext) error                                         // performing any necessary changes before inserting the business object into the database
	IsValid(BloContext) error                                                    // checking the business object's validity according to the business logic
	ChangeAfterInsert(BloContext) error                                          // performing any necessary changes after inserting the business object into the database
	CheckAndChangeAfterRead(bloCtx BloContext) error                             // performing any necessary checks and changes after reading the business object from the database
	CheckAndChangeBeforeUpdate(bloCtx BloContext, instore IBusinessObject) error // performing any necessary checks and changes before updating the business object in the database
	ChangeAfterUpdate(bloCtx BloContext) error                                   // performing any necessary changes after updating the business object in the database

	// utilities for accessing properties and relationships without using reflection
	GetModelName() utils.ModelName                                          // returning the name of the business object's model
	Clone(withFields, withRelationships bool) IBusinessObject               // returning a copy of the business object, optionally including fields and relationships
	GetValueAsString(string) string                                         // returning a BO's field's value, given the field's name
	SetValueAsString(string, string) error                                  // setting a BO's field's value, given the field's name
	SetRelationshipValue(relName string, value IBusinessObject) error       // setting a single-valued relationship's target, given the relationship's name - without using reflection
	AddRelationshipValue(relName string, value IBusinessObject) error       // appending a target to a multi-valued relationship, given the relationship's name - without using reflection
	ClearRelationshipValue(relName string) error                            // resetting a multi-valued relationship to an empty slice, given the relationship's name - without using reflection
	GetSingleRelationshipValue(relName string) (IBusinessObject, error)     // returning the target of a single-valued relationship, given the relationship's name - without using reflection
	GetMultipleRelationshipValue(relName string) ([]IBusinessObject, error) // returning the targets of a multi-valued relationship, given the relationship's name - without using reflection
	IsModelValid() error                                                    // checking a business object's general validity - without using reflection
	RemoveCycles()                                                          // removing any cycles from the business object, without using reflection

	// work with other BOs
	DiffWith(other IBusinessObject, forLinks map[string]bool) (IBusinessObject, IBusinessObject) // creates 2 synthetic instances gathering the added and removed relationships

	// technical stuff
	getModel(this IBusinessObject) IBusinessObjectModel // returning the business object's model
	setLoaded(loadedRelationships)                      // setting the business object's loaded relationships
	getLoaded() loadedRelationships                     // returning the business object's loaded relationships
}

// ------------------------------------------------------------------------------------------------
// Common implementation for business objects - Should be part of any BO's inheritance
// ------------------------------------------------------------------------------------------------

type BObjID int64

type boKey struct {
	id  BObjID
	mdl utils.ModelName
}

type BusinessObject struct {
	// properties common to all business objects
	ID           BObjID              `json:"id,omitempty"           io:"o*" desc:"The unique identifier of this business object"`
	Loaded       loadedRelationships `json:"_loaded,omitempty"      io:"o*" desc:"The direct relationships that have been loaded for this business object"`
	Creation     *time.Time          `json:"creation,omitempty"     io:"o*" desc:"The creation timestamp of this business object"`
	Modification *time.Time          `json:"modification,omitempty" io:"o*" desc:"The modification timestamp of this business object"`
	ModelName    utils.ModelName     `json:"mdl,omitempty"          io:"in" desc:"The name of the business object's model, sometimes used to resolve polymorphic relationships"`
	preID        int                 `json:"-"                      io:"o*" desc:"A temporary identifier in Business Objects lists"`

	// technical stuff
	model IBusinessObjectModel
	key   boKey
}

var _ IBusinessObject = (*BusinessObject)(nil)

// Basic accessors
func (thisBO *BusinessObject) GetID() BObjID                   { return thisBO.ID }
func (thisBO *BusinessObject) setID(id BObjID)                 { thisBO.ID = id }
func (thisBO *BusinessObject) GetCreation() *time.Time         { return thisBO.Creation }
func (thisBO *BusinessObject) setCreation(creation *time.Time) { thisBO.Creation = creation }
func (thisBO *BusinessObject) GetModification() *time.Time     { return thisBO.Modification }
func (thisBO *BusinessObject) setModification(modification *time.Time) {
	thisBO.Modification = modification
}
func (thisBO *BusinessObject) GetPreID() int                    { return thisBO.preID }
func (thisBO *BusinessObject) setPreID(preID int)               { thisBO.preID = preID }
func (thisBO *BusinessObject) getKey() boKey                    { return thisBO.key }
func (thisBO *BusinessObject) setKey(key boKey)                 { thisBO.key = key }
func (thisBO *BusinessObject) SetParent(parent IBusinessObject) { panic("unimplemented") }

// Triggers - default implems
func (thisBO *BusinessObject) ChangeBeforeInsert(BloContext) error             { return nil }
func (thisBO *BusinessObject) IsValid(BloContext) error                        { return nil }
func (thisBO *BusinessObject) ChangeAfterInsert(BloContext) error              { return nil }
func (thisBO *BusinessObject) CheckAndChangeAfterRead(bloCtx BloContext) error { return nil }
func (thisBO *BusinessObject) CheckAndChangeBeforeUpdate(bloCtx BloContext, instore IBusinessObject) error {
	return nil
}
func (thisBO *BusinessObject) ChangeAfterUpdate(bloCtx BloContext) error { return nil }

// Utilities - default implems
func (thisBO *BusinessObject) GetModelName() utils.ModelName { panic("unimplemented") }
func (thisBO *BusinessObject) Clone(withFields, withRelationships bool) IBusinessObject {
	panic("unimplemented")
}
func (thisBO *BusinessObject) GetValueAsString(string) string        { panic("unimplemented") }
func (thisBO *BusinessObject) SetValueAsString(string, string) error { panic("unimplemented") }
func (thisBO *BusinessObject) SetRelationshipValue(relName string, value IBusinessObject) error {
	panic("unimplemented")
}
func (thisBO *BusinessObject) AddRelationshipValue(relName string, value IBusinessObject) error {
	panic("unimplemented")
}
func (thisBO *BusinessObject) ClearRelationshipValue(relName string) error { panic("unimplemented") }
func (thisBO *BusinessObject) GetSingleRelationshipValue(relName string) (IBusinessObject, error) {
	panic("unimplemented")
}
func (thisBO *BusinessObject) GetMultipleRelationshipValue(relName string) ([]IBusinessObject, error) {
	panic("unimplemented")
}
func (thisBO *BusinessObject) IsModelValid() error { panic("unimplemented") }
func (thisBO *BusinessObject) RemoveCycles()       { panic("unimplemented") }

// Working with other BOs - default implems
func (thisBO *BusinessObject) DiffWith(other IBusinessObject, forLinks map[string]bool) (IBusinessObject, IBusinessObject) {
	panic("unimplemented")
}

// Technical stuff
func (thisBO *BusinessObject) setLoaded(loaded loadedRelationships) { thisBO.Loaded = loaded }
func (thisBO *BusinessObject) getLoaded() loadedRelationships       { return thisBO.Loaded }

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

// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

func KeyFor(bo IBusinessObject) boKey {
	if bo.getKey().id == 0 {
		if bo.GetID() == 0 {
			panic("cannot generate key for business object with no ID")
		}

		bo.setKey(boKey{id: bo.GetID(), mdl: bo.GetModelName()})
	}

	return bo.getKey()
}

func (key boKey) String() string {
	return fmt.Sprintf("%s-%d", key.mdl, key.id)
}
