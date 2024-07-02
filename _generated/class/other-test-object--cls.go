// Generated file, do not edit!
package class

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the OtherTestObject class
type otherTestObjectClass struct {
	g.IBusinessObjectClass
	boolProp         *g.BoolField
	int64Prop        *g.BigIntField
	stringProp       *g.StringField
	intProp          *g.IntField
	real32Prop       *g.RealField
	real64Prop       *g.DoubleField
	boolPropCustom   *g.BoolField
	int64PropCustom  *g.BigIntField
	stringPropCustom *g.StringField
	intPropCustom    *g.IntField
	real32PropCustom *g.RealField
	real64PropCustom *g.DoubleField
	dateProp         *g.DateField
	otherEnumProp    *g.EnumField
}

// this is the main way to refer to the OtherTestObject class in the applicative code
func OtherTestObject() *otherTestObjectClass {
	return otherTestObject
}

// internal variables
var (
	otherTestObject     *otherTestObjectClass
	otherTestObjectOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newOtherTestObjectClass() *otherTestObjectClass {
	newClass := &otherTestObjectClass{IBusinessObjectClass: g.NewClass()}
	newClass.boolProp = g.NewBoolField(newClass, "BoolProp", false)
	newClass.int64Prop = g.NewBigIntField(newClass, "Int64Prop", false)
	newClass.stringProp = g.NewStringField(newClass, "StringProp", false)
	newClass.intProp = g.NewIntField(newClass, "IntProp", false)
	newClass.real32Prop = g.NewRealField(newClass, "Real32Prop", false)
	newClass.real64Prop = g.NewDoubleField(newClass, "Real64Prop", false)
	newClass.boolPropCustom = g.NewBoolField(newClass, "BoolPropCustom", false)
	newClass.int64PropCustom = g.NewBigIntField(newClass, "Int64PropCustom", false)
	newClass.stringPropCustom = g.NewStringField(newClass, "StringPropCustom", false)
	newClass.intPropCustom = g.NewIntField(newClass, "IntPropCustom", false)
	newClass.real32PropCustom = g.NewRealField(newClass, "Real32PropCustom", false)
	newClass.real64PropCustom = g.NewDoubleField(newClass, "Real64PropCustom", false)
	newClass.dateProp = g.NewDateField(newClass, "DateProp", false)
	newClass.otherEnumProp = g.NewEnumField(newClass, "OtherEnumProp", false)

	return newClass
}

// making sure the OtherTestObject class exists at app startup
func init() {
	otherTestObjectOnce.Do(func() {
		otherTestObject = newOtherTestObjectClass()
	})

	// this helps dynamically access to the OtherTestObject class
	g.RegisterClass("OtherTestObject", otherTestObject)
}

// accessing all the OtherTestObject class' properties and relationships

func (o *otherTestObjectClass) BoolProp() *g.BoolField {
	return o.boolProp
}

func (o *otherTestObjectClass) Int64Prop() *g.BigIntField {
	return o.int64Prop
}

func (o *otherTestObjectClass) StringProp() *g.StringField {
	return o.stringProp
}

func (o *otherTestObjectClass) IntProp() *g.IntField {
	return o.intProp
}

func (o *otherTestObjectClass) Real32Prop() *g.RealField {
	return o.real32Prop
}

func (o *otherTestObjectClass) Real64Prop() *g.DoubleField {
	return o.real64Prop
}

func (o *otherTestObjectClass) BoolPropCustom() *g.BoolField {
	return o.boolPropCustom
}

func (o *otherTestObjectClass) Int64PropCustom() *g.BigIntField {
	return o.int64PropCustom
}

func (o *otherTestObjectClass) StringPropCustom() *g.StringField {
	return o.stringPropCustom
}

func (o *otherTestObjectClass) IntPropCustom() *g.IntField {
	return o.intPropCustom
}

func (o *otherTestObjectClass) Real32PropCustom() *g.RealField {
	return o.real32PropCustom
}

func (o *otherTestObjectClass) Real64PropCustom() *g.DoubleField {
	return o.real64PropCustom
}

func (o *otherTestObjectClass) DateProp() *g.DateField {
	return o.dateProp
}

func (o *otherTestObjectClass) OtherEnumProp() *g.EnumField {
	return o.otherEnumProp
}
