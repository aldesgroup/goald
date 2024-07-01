// Generated file, do not edit!
package class

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the TestObject class
type testObjectClass struct {
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
	enumProp         *g.EnumField
	dateProp         *g.DateField
	otherEnumProp    *g.EnumField
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

func (t *testObjectClass) Int64Prop() *g.BigIntField {
	return t.int64Prop
}

func (t *testObjectClass) StringProp() *g.StringField {
	return t.stringProp
}

func (t *testObjectClass) IntProp() *g.IntField {
	return t.intProp
}

func (t *testObjectClass) Real32Prop() *g.RealField {
	return t.real32Prop
}

func (t *testObjectClass) Real64Prop() *g.DoubleField {
	return t.real64Prop
}

func (t *testObjectClass) BoolPropCustom() *g.BoolField {
	return t.boolPropCustom
}

func (t *testObjectClass) Int64PropCustom() *g.BigIntField {
	return t.int64PropCustom
}

func (t *testObjectClass) StringPropCustom() *g.StringField {
	return t.stringPropCustom
}

func (t *testObjectClass) IntPropCustom() *g.IntField {
	return t.intPropCustom
}

func (t *testObjectClass) Real32PropCustom() *g.RealField {
	return t.real32PropCustom
}

func (t *testObjectClass) Real64PropCustom() *g.DoubleField {
	return t.real64PropCustom
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
