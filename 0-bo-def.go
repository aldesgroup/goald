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
	Class() IBusinessObjectClass
	getClassName() className
	setClassName(className)
	GetID() BObjID
	setID(string)
	GetLabel() string

	// business logic
	ChangeBeforeInsert(BloContext) error
	IsValid(BloContext) error
	ChangeAfterInsert(BloContext) error

	// utils
	GetValueAsString(propName string) string
	SetValueAsString(propName, value string) error
}

// ------------------------------------------------------------------------------------------------
// Common implementation for business objects - Should be part of any BO's inheritance
// ------------------------------------------------------------------------------------------------

type BObjID string // probably a UUID here

type BusinessObject struct {
	class     IBusinessObjectClass
	className className
	ID        BObjID `json:"ID"`
	Label     string `json:"Label,omitempty"`
}

var _ IBusinessObject = (*BusinessObject)(nil)

func (thisBO *BusinessObject) Class() IBusinessObjectClass {
	if thisBO.class == nil {
		thisBO.class = classForName(thisBO.className)
	}

	if thisBO.class == nil {
		panic("unknown class for a business object!")
	}

	return thisBO.class
}

/* default implementations */
func (thisBO *BusinessObject) getClassName() className                 { return thisBO.className }
func (thisBO *BusinessObject) setClassName(cn className)               { thisBO.className = cn }
func (thisBO *BusinessObject) GetID() BObjID                           { return thisBO.ID }
func (thisBO *BusinessObject) setID(id string)                         { thisBO.ID = BObjID(id) }
func (thisBO *BusinessObject) GetLabel() string                        { return thisBO.Label }
func (thisBO *BusinessObject) ChangeBeforeInsert(BloContext) error     { return nil }
func (thisBO *BusinessObject) IsValid(BloContext) error                { return nil }
func (thisBO *BusinessObject) ChangeAfterInsert(BloContext) error      { return nil }
func (thisBO *BusinessObject) GetValueAsString(propName string) string { return "- not implemented -" }
func (thisBO *BusinessObject) SetValueAsString(string, string) error   { return Error("not implemented") }

// ------------------------------------------------------------------------------------------------
// Modelling enum types
// ------------------------------------------------------------------------------------------------

// IEnum must be implemented by every enum type
type IEnum interface {
	fmt.Stringer // each enum value has a default label
	Val() int
	Values() map[int]string
}

// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

func InstanceOf(class IBusinessObjectClass) BusinessObject {
	return BusinessObject{
		class: class,
	}
}
