// ------------------------------------------------------------------------------------------------
// Ensuring we have a very limited use of the 'reflect' package.
// THIS FILE SHOULD BE THE ONLY PLACE WHERE WE IMPORT THE REFLECT PACKAGE!
// Basically, we only want to use reflection:
// - 1) when generating code (which precisely allows to avoid reflection)
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
	fields map[string]*GoaldField
}

func newType(val reflect.Type) *GoaldType {
	return &GoaldType{val, nil}
}

type GoaldField struct {
	val reflect.StructField
	typ *GoaldType
}

func newField(val reflect.StructField) *GoaldField {
	return &GoaldField{val: val, typ: newType(val.Type)}
}

func TypeOf(arg any, bare bool) *GoaldType {
	if bare {
		// getting the name of "MyType" rather than "*MyType"
		return newType(reflect.TypeOf(arg).Elem())
	}

	return newType(reflect.TypeOf(arg))
}

func PointerTo(arg *GoaldType) *GoaldType {
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

func (t *GoaldType) Field(index int) *GoaldField {
	return newField(t.val.Field(index))
}

func (t *GoaldType) NumField() int {
	return t.val.NumField()
}

func (t *GoaldType) Implements(other *GoaldType) bool {
	return t.val.Implements(other.val)
}

func (t *GoaldType) Equals(other *GoaldType) bool {
	return t.val == other.val
}

func (t *GoaldType) Name() string {
	return t.val.Name()
}

func (t *GoaldType) Elem() *GoaldType {
	return newType(t.val.Elem())
}

func (t *GoaldType) FieldByName(name string) *GoaldField {
	f, found := t.val.FieldByName(name)
	if !found {
		panic(fmt.Sprintf("no field '%s' on type '%s'", name, t.Name()))
	}
	return newField(f)
}

func (t *GoaldType) String() string {
	return t.val.String()
}

func (t *GoaldType) PkgPath() string {
	return t.val.PkgPath()
}

// --- fields ----------------------------------------------------------------------------------

func (f *GoaldField) Type() *GoaldType {
	return f.typ
}

// func

// enumType := bObjectEntry.getAllProperties()[fieldName].Type
// enumTypeString := enumType.Name() // e.g.: MyEnumType
// if enumPkg := enumType.PkgPath(); enumPkg != bObjectEntry.bObjType.PkgPath() {
// 	enumTypeString = enumType.String() /// e.g.: thatpackage.MyEnumType
// 	importsMap[enumPkg] = true
// }

func (f *GoaldField) IsAnonymous() bool {
	return f.val.Anonymous
}

func (f *GoaldField) Name() string {
	return f.val.Name
}

// WARNING: concurrent access not handled for now
// func (thisBOEntry *businessObjectEntry) getAllProperties() map[string]reflect.StructField {
// 	if thisBOEntry.properties == nil {
// 		thisBOEntry.properties = make(map[string]reflect.StructField)
// 		for _, field := range getAllFields(thisBOEntry.bObjType) {
// 			thisBOEntry.properties[field.Name] = field
// 		}
// 	}

// 	return thisBOEntry.properties
// }

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
	TypeFamilyRELATIONSHIP
)

var TypeFamilys = map[int]string{
	int(TypeFamilyUNKNOWN):      "unknown",
	int(TypeFamilyBOOL):         "boolean",
	int(TypeFamilySTRING):       "string",
	int(TypeFamilyINT):          "integer",
	int(TypeFamilyBIGINT):       "bigint",
	int(TypeFamilyREAL):         "real number",
	int(TypeFamilyDOUBLE):       "real number 64",
	int(TypeFamilyDATE):         "date",
	int(TypeFamilyENUM):         "enum",
	int(TypeFamilyRELATIONSHIP): "relationship",
}

func (thisProperty TypeFamily) String() string {
	return TypeFamilys[int(thisProperty)]
}

// Val helps implement the IEnum interface
func (thisProperty TypeFamily) Val() int {
	return int(thisProperty)
}

// Values helps implement the IEnum interface
func (thisProperty TypeFamily) Values() map[int]string {
	return TypeFamilys
}

// GetTypeFamily returns the type family of a given structfield
func GetTypeFamily(field *GoaldField, iBoTypeFamily, enumTypeFamily *GoaldType) (TypeFamily TypeFamily, multiple bool) {

	// to debug - to comment/uncomment when needed
	// if structField.Name == "Num" {
	// fmt.Printf("\n--------------------------")
	// fmt.Printf("\nName: %s ", structField.Name)
	// fmt.Printf("\nType: %s ", structField.Type)
	// fmt.Printf("\nKind: %s ", structField.Type.Kind())
	// }

	// a business object's real property must be exported, and therefore PkgPath should be empty
	// Cf. https://golang.org/pkg/reflect/#StructField
	if fieldType := field.Type(); field.val.PkgPath == "" {
		// detecting an enum
		if fieldType.Implements(enumTypeFamily) {
			return TypeFamilyENUM, false
		}

		// detecting a time
		if fieldType.Equals(typeTIMExPTR) {
			return TypeFamilyDATE, false
		}

		// detecting the basic types here
		switch fieldKind := fieldType.val.Kind(); fieldKind {
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

			// // detecting an enum list
			// if innerSliceType := fieldType.Elem(); fieldKind == reflect.Slice && innerSliceType.Implements(TypeIENUM) {
			// 	return TypeFamilyENUM, true
			// }
		}

		// detecting a multiple relationship to business objects
		if innerSliceType := fieldType.Elem(); fieldType.val.Kind() == reflect.Slice && innerSliceType.Implements(iBoTypeFamily) {
			return TypeFamilyRELATIONSHIP, true
		}

		// detecting a single relationship to a business object
		if fieldType.val.Kind() == reflect.Ptr && fieldType.Implements(iBoTypeFamily) {
			return TypeFamilyRELATIONSHIP, false
		}
	}

	// this happens with technical fields !
	return TypeFamilyUNKNOWN, false
}

// func getAllFields(bObjType reflect.Type) (fields []reflect.StructField) {
// 	for i := 0; i < bObjType.NumField(); i++ {
// 		field := bObjType.Field(i)
// 		if field.Anonymous {
// 			fields = append(fields, getAllFields(field.Type)...)
// 		} else {
// 			fields = append(fields, field)
// 		}
// 	}

// 	return
// }

// ------------------------------------------------------------------------------------------------
// misc dynamic stuff using reflection
// ------------------------------------------------------------------------------------------------

// returning the name of the given function
func GetFnName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}
