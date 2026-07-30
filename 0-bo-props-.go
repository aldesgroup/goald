// ------------------------------------------------------------------------------------------------
// The code here is about describing the 2 types of properties, i.e. the fields - and their
// derivatives - and the relationships with other business objects
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"time"
	"unicode"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/reflection"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Business object properties, whether fields or relationships
// ------------------------------------------------------------------------------------------------

type IBusinessObjectProperty interface {
	// public methods
	GetName() string  // the property's name, as declared in the struct
	IsMultiple() bool // the property's multiplicity; false = 1, true = N
	SetUnique()       // to set this property as unique in the database (UNIQUE)
	SetSecret()       // to set this property as secret
	SetPersonal()     // to set this property as personal

	// private methods
	ownerModel() IBusinessObjectModel       // the property's owner model
	setOwner(IBusinessObjectModel)          // to set the property's owner model
	getPropertyType() propertyType          // the property's type, as detected by the codegen phase
	getColumnName() string                  // if this property - field or relationship - is persisted on the owner's table, then this is the name of the corresponding column
	isNotPersisted() bool                   // if true, then this property does not have a corresponding column in the BO's table
	getStructField() *reflection.GoaldField // the struct field corresponding to this property, as detected by the codegen phase
	getTag(tagName string) string           // to get the value of a given tag on the struct field corresponding to this property
	isMandatoryInput() bool                 // if true, then this property is a mandatory input for the associated endpoint
	isPureOutput() bool                     // if true, then this property is a pure output for the associated endpoint
	IsRequiredInDb() bool                   // if true, then this property is required in the database (NOT NULL)
	SetRequiredInDb()                       // to set this property as required in the database (NOT NULL)
	isUnique() bool                         // if true, then this property is unique in the database (UNIQUE)
	getUniqueConstraintName() string        // to get the name of the unique constraint associated with this property, if any
	isSecret() bool                         // if true, then this property is secret
	isPersonal() bool                       // if true, then this property is personal
	getDeclaringBO() utils.ModelName        // the name of the business object that actually declares this property (instead of inheriting it from a super model)
	setTechnical()                          // to set this property as technical
}

type ioType string

const ioTypeINPUT ioType = "in"
const ioTypeINPUTxMANDATORY ioType = "i*"
const ioTypePURExOUTPUT ioType = "o*"

type businessObjectProperty struct {
	owner        IBusinessObjectModel   // the property's owner model
	declaringBO  utils.ModelName        // the name of the business object that actually declares this property (instead of inheriting it from a super model)
	name         string                 // the property's name, as declared in the struct
	propType     propertyType           // the property's type, as detected by the codegen phase
	multiple     bool                   // the property's multiplicity; false = 1, true = N
	columnName   string                 // if this property - field or relationship - is persisted on the owner's table
	notPersisted bool                   // if true, then this property does not have a corresponding column in the BO's table
	structField  *reflection.GoaldField // the struct field corresponding to this property, as detected by the codegen phase
	requiredInDb bool                   // if true, then this property is required in the database (NOT NULL)
	unique       bool                   // if true, then this property is unique in the database (UNIQUE)
	secret       bool                   // if true, then this property is secret
	personal     bool                   // if true, then this property is personal
	technical    bool                   // if true, then this property is purely technical
}

func (prop *businessObjectProperty) ownerModel() IBusinessObjectModel {
	return prop.owner
}

func (prop *businessObjectProperty) setOwner(owner IBusinessObjectModel) {
	prop.owner = owner
}

func (prop *businessObjectProperty) GetName() string {
	return prop.name
}

func (prop *businessObjectProperty) getPropertyType() propertyType {
	return prop.propType
}

func (prop *businessObjectProperty) getColumnName() string {
	if prop.columnName == "" {
		prop.columnName = core.PascalToSnake(prop.name)
		if prop.propType.IsRelationship() {
			prop.columnName += suffixID
		}
		if prop.technical {
			prop.columnName = "_" + prop.columnName
		}
	}

	return prop.columnName
}

func (prop *businessObjectProperty) IsMultiple() bool {
	return prop.multiple
}

func (prop *businessObjectProperty) isNotPersisted() bool {
	return prop.notPersisted
}

func (prop *businessObjectProperty) getStructField() *reflection.GoaldField {
	if prop.structField == nil {
		structField := prop.ownerModel().getType().FieldByName(prop.name)
		prop.structField = &structField
	}
	return prop.structField
}

func (prop *businessObjectProperty) getTag(tagName string) string {
	return prop.getStructField().Tag().Get(tagName)
}

func (prop *businessObjectProperty) isMandatoryInput() bool {
	return prop.getTag("io") == string(ioTypeINPUTxMANDATORY)
}

func (prop *businessObjectProperty) isPureOutput() bool {
	return prop.getTag("io") == string(ioTypePURExOUTPUT)
}

func (prop *businessObjectProperty) IsRequiredInDb() bool {
	return prop.requiredInDb
}

func (prop *businessObjectProperty) SetRequiredInDb() {
	prop.requiredInDb = true
}

func (prop *businessObjectProperty) isUnique() bool {
	return prop.unique
}

func (prop *businessObjectProperty) SetUnique() {
	prop.unique = true
}

func (prop *businessObjectProperty) getUniqueConstraintName() string {
	return prefixUK + prop.ownerModel().getTableName(false) + "__" + prop.getColumnName()
}

func (prop *businessObjectProperty) SetSecret() {
	prop.secret = true
}

func (prop *businessObjectProperty) isSecret() bool {
	return prop.secret
}

func (prop *businessObjectProperty) SetPersonal() {
	prop.personal = true
}

func (prop *businessObjectProperty) isPersonal() bool {
	return prop.personal
}

func (prop *businessObjectProperty) getDeclaringBO() utils.ModelName {
	return prop.declaringBO
}

func (prop *businessObjectProperty) setTechnical() {
	prop.technical = true
}

// ------------------------------------------------------------------------------------------------
// Utils: defining the property types Goald supports, and the way to detect them with reflection
// ------------------------------------------------------------------------------------------------

// propertyType represents the type of a business object's property
type propertyType int

const (
	propertyTypeUNKNOWN propertyType = iota
	propertyTypeBOOL
	propertyTypeSTRING
	propertyTypeINT
	propertyTypeBIGINT
	propertyTypeREAL
	propertyTypeDOUBLE
	propertyTypeDATE
	propertyTypeENUM
	propertyTypeRELATIONSHIPxMONOM
	propertyTypeRELATIONSHIPxPOLYM
)

var typeFamilies = map[int]string{
	int(propertyTypeUNKNOWN):            "unknown",
	int(propertyTypeBOOL):               "boolean",
	int(propertyTypeSTRING):             "string",
	int(propertyTypeINT):                "integer",
	int(propertyTypeBIGINT):             "bigint",
	int(propertyTypeREAL):               "real number",
	int(propertyTypeDOUBLE):             "real number 64",
	int(propertyTypeDATE):               "date",
	int(propertyTypeENUM):               "enum",
	int(propertyTypeRELATIONSHIPxMONOM): "relationship (monomorphic)",
	int(propertyTypeRELATIONSHIPxPOLYM): "relationship (polymorphic)",
}

var (
	typeTIMExPTR = reflection.TypeOf((*time.Time)(nil), false)
)

func (thispropertyType propertyType) String() string {
	return typeFamilies[int(thispropertyType)]
}

// Val helps implement the IEnum interface
func (thispropertyType propertyType) Val() int {
	return int(thispropertyType)
}

// Values helps implement the IEnum interface
func (thispropertyType propertyType) Values() map[int]string {
	return typeFamilies
}

// Tells if we have a relationship here
func (thispropertyType propertyType) IsRelationship() bool {
	return thispropertyType == propertyTypeRELATIONSHIPxMONOM || thispropertyType == propertyTypeRELATIONSHIPxPOLYM
}

// detectPropertyType returns the type family of a given structfield
func detectPropertyType(field reflection.GoaldField, iBopropertyType, enumpropertyType reflection.GoaldType) (propertyType propertyType, multiple bool) {
	// a business object's real property must be exported, and therefore PkgPath should be empty
	// Cf. https://golang.org/pkg/reflect/#StructField
	if fieldType := field.Type(); field.PkgPath() == "" {
		// getting the field kind
		fieldKind := fieldType.Kind()

		// handling the case where we have a slice in here
		if fieldKind == reflection.KindSLICE {
			// what's in there?
			innerSliceType := fieldType.Elem()
			innerSliceKind := innerSliceType.Kind()

			// detecting an enum
			if innerSliceType.Implements(enumpropertyType) {
				return propertyTypeENUM, true
			}

			// detecting a polymorphic type, i.e. an interface; this should point to something implementing IBusinessObject
			if innerSliceKind == reflection.KindINTERFACE && innerSliceType.Implements(iBopropertyType) {
				return propertyTypeRELATIONSHIPxPOLYM, true
			}

			// detecting a single relationship to a business object
			if innerSliceKind == reflection.KindPTR && innerSliceType.Implements(iBopropertyType) {
				return propertyTypeRELATIONSHIPxMONOM, true
			}

		} else { // we have a single element here

			// detecting an enum
			if fieldType.Implements(enumpropertyType) {
				return propertyTypeENUM, false
			}

			// detecting a time
			if fieldType.Equals(typeTIMExPTR) {
				return propertyTypeDATE, false
			}

			// detecting the basic types here
			switch fieldKind {
			case reflection.KindBOOL:
				return propertyTypeBOOL, false

			case reflection.KindSTRING:
				return propertyTypeSTRING, false

			case reflection.KindINT:
				return propertyTypeINT, false

			case reflection.KindINT64:
				return propertyTypeBIGINT, false

			case reflection.KindFLOAT32:
				return propertyTypeREAL, false

			case reflection.KindFLOAT64:
				return propertyTypeDOUBLE, false
			}

			// detecting a polymorphic type, i.e. an interface; this should point to something implementing IBusinessObject
			if fieldKind == reflection.KindINTERFACE && fieldType.Implements(iBopropertyType) {
				return propertyTypeRELATIONSHIPxPOLYM, false
			}

			// detecting a single relationship to a business object
			if fieldKind == reflection.KindPTR && fieldType.Implements(iBopropertyType) {
				return propertyTypeRELATIONSHIPxMONOM, false
			}
		}
	}

	// detect if the field name starts with an uppercase letter, this should be allowed for non-handled property types
	if unicode.IsUpper(rune(field.Name()[0])) {
		panic(fmt.Sprintf("This exported property ('%s') should have a type handled by Goald. "+
			"If this is a relationship, make sure you're using a pointer to it: *somepackage.SomeModel", field.Name()))
	}

	// this happens with technical fields !
	return propertyTypeUNKNOWN, false
}
