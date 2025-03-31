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
	"time"
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

func (f GoaldField) IsAnonymous() bool {
	return f.val.Anonymous
}

func (f GoaldField) Name() string {
	return f.val.Name
}

// ------------------------------------------------------------------------------------------------
// global variables
// ------------------------------------------------------------------------------------------------

var (
	typeTIMExPTR = TypeOf((*time.Time)(nil), false)
)

// ------------------------------------------------------------------------------------------------
// defining type families
// ------------------------------------------------------------------------------------------------

// TypeFamily represents the type of a business object's property
type TypeFamily int

const (
	TypeFamilyUNKNOWN TypeFamily = iota - 1
	TypeFamilyBOOL
	TypeFamilySTRING
	TypeFamilyINT
	TypeFamilyBIGINT
	TypeFamilyREAL
	TypeFamilyDOUBLE
	TypeFamilyDATE
	TypeFamilyENUM
	TypeFamilyRELATIONSHIPxMONOM
	TypeFamilyRELATIONSHIPxPOLYM
)

var typeFamilies = map[int]string{
	int(TypeFamilyUNKNOWN):            "unknown",
	int(TypeFamilyBOOL):               "boolean",
	int(TypeFamilySTRING):             "string",
	int(TypeFamilyINT):                "integer",
	int(TypeFamilyBIGINT):             "bigint",
	int(TypeFamilyREAL):               "real number",
	int(TypeFamilyDOUBLE):             "real number 64",
	int(TypeFamilyDATE):               "date",
	int(TypeFamilyENUM):               "enum",
	int(TypeFamilyRELATIONSHIPxMONOM): "relationship (monomorphic)",
	int(TypeFamilyRELATIONSHIPxPOLYM): "relationship (polymorphic)",
}

func (thisProperty TypeFamily) String() string {
	return typeFamilies[int(thisProperty)]
}

// Val helps implement the IEnum interface
func (thisProperty TypeFamily) Val() int {
	return int(thisProperty)
}

// Values helps implement the IEnum interface
func (thisProperty TypeFamily) Values() map[int]string {
	return typeFamilies
}

// Tells if we have a relationship here
func (thisProperty TypeFamily) IsRelationship() bool {
	return thisProperty == TypeFamilyRELATIONSHIPxMONOM || thisProperty == TypeFamilyRELATIONSHIPxPOLYM
}

// GetTypeFamily returns the type family of a given structfield
func GetTypeFamily(field GoaldField, iBoTypeFamily, enumTypeFamily GoaldType) (TypeFamily TypeFamily, multiple bool) {

	// to debug - to comment/uncomment when needed
	// if structField.Name == "Num" {
	// fmt.Printf("\n--------------------------")
	// fmt.Printf("\nName: %s ", field.Name())
	// fmt.Printf("\nType: %s ", field.Type())
	// fmt.Printf("\nType: %s ", field.Type().val.Kind())
	// fmt.Printf("\nKind: %s ", field.Type().Kind)
	// }

	// a business object's real property must be exported, and therefore PkgPath should be empty
	// Cf. https://golang.org/pkg/reflect/#StructField
	if fieldType := field.Type(); field.val.PkgPath == "" {
		// getting the field kind
		fieldKind := fieldType.val.Kind()

		// handling the case where we have a slice in here
		if fieldKind == reflect.Slice {
			// what's in there?
			innerSliceType := fieldType.Elem()
			innerSliceKind := innerSliceType.val.Kind()

			// detecting an enum
			if innerSliceType.Implements(enumTypeFamily) {
				return TypeFamilyENUM, true
			}

			// detecting a polymorphic type, i.e. an interface; this should point to something implementing IBusinessObject
			if innerSliceKind == reflect.Interface && innerSliceType.Implements(iBoTypeFamily) {
				return TypeFamilyRELATIONSHIPxPOLYM, true
			}

			// detecting a single relationship to a business object
			if innerSliceKind == reflect.Ptr && innerSliceType.Implements(iBoTypeFamily) {
				return TypeFamilyRELATIONSHIPxMONOM, true
			}

		} else { // we have a single element here

			// detecting an enum
			if fieldType.Implements(enumTypeFamily) {
				return TypeFamilyENUM, false
			}

			// detecting a time
			if fieldType.Equals(typeTIMExPTR) {
				return TypeFamilyDATE, false
			}

			// detecting the basic types here
			switch fieldKind {
			case reflect.Bool:
				return TypeFamilyBOOL, false

			case reflect.String:
				return TypeFamilySTRING, false

			case reflect.Int:
				return TypeFamilyINT, false

			case reflect.Int64:
				return TypeFamilyBIGINT, false

			case reflect.Float32:
				return TypeFamilyREAL, false

			case reflect.Float64:
				return TypeFamilyDOUBLE, false
			}

			// detecting a polymorphic type, i.e. an interface; this should point to something implementing IBusinessObject
			if fieldKind == reflect.Interface && fieldType.Implements(iBoTypeFamily) {
				return TypeFamilyRELATIONSHIPxPOLYM, false
			}

			// detecting a single relationship to a business object
			if fieldKind == reflect.Ptr && fieldType.Implements(iBoTypeFamily) {
				return TypeFamilyRELATIONSHIPxMONOM, false
			}
		}
	}

	// this happens with technical fields !
	return TypeFamilyUNKNOWN, false
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
