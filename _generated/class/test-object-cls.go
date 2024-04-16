// Generated file, do not edit!
package class

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the TestObject class
type testObjectClass struct {
	g.IBusinessObjectClass
	boolProp      *g.BoolField
	int64Prop     *g.Int64Field
	stringProp    *g.StringField
	intProp       *g.IntField
	real32Prop    *g.Real32Field
	real64Prop    *g.Real64Field
	enumProp      *g.EnumField
	dateProp      *g.DateField
	otherEnumProp *g.EnumField
}

// this is the main way to refer to the TestObject class in the applicative code
func TestObject() *testObjectClass {
	return testObject
}

// internal variables
var (
	testObject     *testObjectClass
	testObjectOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newTestObjectClass() *testObjectClass {
	newClass := &testObjectClass{IBusinessObjectClass: g.NewClass()}
	newClass.boolProp = g.NewBoolField(newClass, "BoolProp", false)
	newClass.int64Prop = g.NewInt64Field(newClass, "Int64Prop", false)
	newClass.stringProp = g.NewStringField(newClass, "StringProp", false)
	newClass.intProp = g.NewIntField(newClass, "IntProp", false)
	newClass.real32Prop = g.NewReal32Field(newClass, "Real32Prop", false)
	newClass.real64Prop = g.NewReal64Field(newClass, "Real64Prop", false)
	newClass.enumProp = g.NewEnumField(newClass, "EnumProp", false)
	newClass.dateProp = g.NewDateField(newClass, "DateProp", false)
	newClass.otherEnumProp = g.NewEnumField(newClass, "OtherEnumProp", false)

	return newClass
}

// making sure the TestObject class exists at app startup
func init() {
	testObjectOnce.Do(func() {
		testObject = newTestObjectClass()
	})

	// this helps dynamically access to the TestObject class
	g.RegisterClass("TestObject", testObject)
}

// accessing all the TestObject class' properties and relationships

func (t *testObjectClass) BoolProp() *g.BoolField {
	return t.boolProp
}

func (t *testObjectClass) Int64Prop() *g.Int64Field {
	return t.int64Prop
}

func (t *testObjectClass) StringProp() *g.StringField {
	return t.stringProp
}

func (t *testObjectClass) IntProp() *g.IntField {
	return t.intProp
}

func (t *testObjectClass) Real32Prop() *g.Real32Field {
	return t.real32Prop
}

func (t *testObjectClass) Real64Prop() *g.Real64Field {
	return t.real64Prop
}

func (t *testObjectClass) EnumProp() *g.EnumField {
	return t.enumProp
}

func (t *testObjectClass) DateProp() *g.DateField {
	return t.dateProp
}

func (t *testObjectClass) OtherEnumProp() *g.EnumField {
	return t.otherEnumProp
}
