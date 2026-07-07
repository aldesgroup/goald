// ------------------------------------------------------------------------------------------------
// Ensuring we have a very limited use of the 'reflect' package.
// THIS FILE SHOULD BE THE ONLY PLACE WHERE WE IMPORT THE REFLECT PACKAGE!
// Basically, we only want to use reflection:
// - 1) when generating code (which then precisely allows to avoid reflection)
// - 2) in init functions
// NEVER DURING THE RUNTIME! At least not in OUR code (but some 3rd party libraries prolly do)
// ------------------------------------------------------------------------------------------------
package utils

import (
	"fmt"
	"reflect"
	"runtime"
)

// ------------------------------------------------------------------------------------------------
// structs, constructors, utils
// We define our own Type structs to have control on how we use the reflect package
// ------------------------------------------------------------------------------------------------

type GoaldType struct {
	val    reflect.Type
	fields map[string]GoaldField
}

func newType(val reflect.Type) GoaldType {
	return GoaldType{val, nil}
}

type GoaldField struct {
	val reflect.StructField
	typ GoaldType
}

func newField(val reflect.StructField) GoaldField {
	return GoaldField{val: val, typ: newType(val.Type)}
}

func TypeOf(arg any, bare bool) GoaldType {
	if bare {
		// getting the type as "MyType" rather than "*MyType"
		return newType(reflect.TypeOf(arg).Elem())
	}

	return newType(reflect.TypeOf(arg))
}

func PointerTo(arg GoaldType) GoaldType {
	return newType(reflect.PointerTo(arg.val))
}

func TypeNameOf(arg any, bare bool) string {
	if bare {
		// getting the name of "MyType" rather than "*MyType"
		return reflect.TypeOf(arg).Elem().Name()
	}

	return reflect.TypeOf(arg).Name()
}

// ------------------------------------------------------------------------------------------------
// methods - mostly proxied native methods
// ------------------------------------------------------------------------------------------------

// --- types -----------------------------------------------------------------------------------

func (t GoaldType) Field(index int) GoaldField {
	return newField(t.val.Field(index))
}

func (t GoaldType) NumField() int {
	return t.val.NumField()
}

func (t GoaldType) Implements(other GoaldType) bool {
	return t.val.Implements(other.val)
}

func (t GoaldType) Equals(other GoaldType) bool {
	return t.val == other.val
}

func (t GoaldType) Name() string {
	return t.val.Name()
}

func (t GoaldType) Kind() reflect.Kind {
	return t.val.Kind()
}

func (t GoaldType) Elem() GoaldType {
	return newType(t.val.Elem())
}

func (t GoaldType) FieldByName(name string) GoaldField {
	f, found := t.val.FieldByName(name)
	if !found {
		panic(fmt.Sprintf("no field '%s' on type '%s'", name, t.Name()))
	}
	return newField(f)
}

func (t GoaldType) String() string {
	return t.val.String()
}

func (t GoaldType) PkgPath() string {
	return t.val.PkgPath()
}

// --- fields ----------------------------------------------------------------------------------

func (f GoaldField) Type() GoaldType {
	return f.typ
}

func (f GoaldField) PkgPath() string {
	return f.val.PkgPath
}

func (f GoaldField) Val() reflect.StructField {
	return f.val
}

func (f GoaldField) IsAnonymous() bool {
	return f.val.Anonymous
}

func (f GoaldField) Name() string {
	return f.val.Name
}

func (f GoaldField) Tag() reflect.StructTag {
	return f.val.Tag
}

// ------------------------------------------------------------------------------------------------
// values
// ------------------------------------------------------------------------------------------------

type GoaldValue struct {
	val reflect.Value
}

func newValue(val reflect.Value) GoaldValue {
	return GoaldValue{val}
}

func ValueOf(arg any) GoaldValue {
	return newValue(reflect.ValueOf(arg).Elem())
}

func (thisValue GoaldValue) GetFieldValue(fieldName string) any {
	field := thisValue.val.FieldByName(fieldName)
	// TODO  if !field.CanInterface() { return nil, fmt.Errorf("cannot access unexported field: %s", fieldName) }
	return field.Interface()
}

// ------------------------------------------------------------------------------------------------
// misc dynamic stuff using reflection
// ------------------------------------------------------------------------------------------------

// returning the name of the given function
func GetFnName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// ------------------------------------------------------------------------------------------------
// proxied kinds
// ------------------------------------------------------------------------------------------------

const KindSLICE = reflect.Slice
const KindINTERFACE = reflect.Interface
const KindPTR = reflect.Ptr
const KindBOOL = reflect.Bool
const KindSTRING = reflect.String
const KindINT = reflect.Int
const KindINT64 = reflect.Int64
const KindFLOAT32 = reflect.Float32
const KindFLOAT64 = reflect.Float64
