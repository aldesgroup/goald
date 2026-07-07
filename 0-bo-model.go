// ------------------------------------------------------------------------------------------------
// The code here is about having tools to describe business object models, and
// their properties in a general sense
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"sort"
	"time"
	"unicode"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/utils"
)

// TODO pagination all the way

// ------------------------------------------------------------------------------------------------
// Model for a business object model
// ------------------------------------------------------------------------------------------------
type IBusinessObjectModel interface {
	/* public generic methods */

	SetDescription(description string)                     // to set the description of the model
	SetNotPersisted()                                      // to indicate this class has no instance persisted in a database
	SetInDB(db *DB)                                        // to associate the class with the DB where its instances are stored
	SetInDbByName(dbName string)                           // to associate the class with the DB where its instances are stored, using the name of the DB instead of its instance
	SetAbstract()                                          // to indicate this class does not model concrete business objects, but most probably a super class
	AddUniqueCombination(props ...IBusinessObjectProperty) // to indicate that a combination of properties must be unique in the DB

	// access to generic properties (fields & relationships)
	ID() IField

	// private methods
	isNotPersisted() bool
	getInDB() *DB
	getTableName(withSchema bool) string

	// access to the base implementation
	base() *businessObjectModel
	addField(field IField) IField
	getType() *utils.GoaldType
}

type className string

type businessObjectModel struct {
	name                       className                            // the corresponding class name
	description                string                               // the model description
	fields                     map[string]IField                    // the objet's simple properties
	relationships              map[string]*Relationship             // the relationships to other classes
	inDB                       *DB                                  // the associated DB, if any
	inNoDB                     bool                                 // if true, then no associated DB
	abstract                   bool                                 // if true, then is class is mainly used as a super class for others
	tableName                  string                               // if persisted, the name of the corresponding DB table - should be the same as the class name most of the time
	persistedProperties        []IBusinessObjectProperty            // all the properties - fields or relationships - persisted on this class
	uniqueCombinations         map[string][]IBusinessObjectProperty // a combination of properties that must be unique in the DB
	idField                    IField                               // accessor to the ID field
	usedInNativeApp            bool                                 // true if this class is used in the native app
	usedInWebApp               bool                                 // true if this class is used in the web app
	boType                     *utils.GoaldType                     // the Go type associated with this BO model
	relationshipsWithColumn    []*Relationship                      // all the relationships for which this class has a column in its table
	relationshipsWithLinkTable []*Relationship                      // all the relationships for which this class has a column in its table
}

const BoFieldID = "ID" // the name of the ID field, which is a special case in Goald

func NewBusinessObjectModel() IBusinessObjectModel {
	model := &businessObjectModel{
		fields:        map[string]IField{},
		relationships: map[string]*Relationship{},
	}

	// adding the generic fields
	model.idField = NewBigIntField(model, BoFieldID, false)

	return model
}

func (boModel *businessObjectModel) SetInDB(db *DB) {
	boModel.inNoDB = false
	boModel.inDB = db
}

func (boModel *businessObjectModel) SetInDbByName(dbName string) {
	boModel.inNoDB = false
	boModel.inDB = GetDB(dbconn.DbSchemaName(dbName))
}

func (boModel *businessObjectModel) SetDescription(description string) {
	boModel.description = description
}

func (boModel *businessObjectModel) SetNotPersisted() {
	boModel.inNoDB = true
	boModel.inDB = nil
}

func (boModel *businessObjectModel) SetAbstract() {
	boModel.abstract = true
}

func (boModel *businessObjectModel) AddUniqueCombination(props ...IBusinessObjectProperty) {
	// bit of control
	core.PanicMsgIf(len(props) <= 1, "AddUniqueCombination() should be used with at least 2 properties of model '%s'", boModel.name)

	// sorting the props
	sort.Slice(props, func(i, j int) bool {
		return props[i].GetName() < props[j].GetName()
	})

	// building the ck constraint name
	constraintName := string(prefixCK + boModel.getTableName(false))
	for _, prop := range props {
		constraintName = constraintName + "_" + core.PascalToShort(prop.GetName())
	}

	// adding this constraint to the model
	if boModel.uniqueCombinations == nil {
		boModel.uniqueCombinations = map[string][]IBusinessObjectProperty{}
	}
	boModel.uniqueCombinations[constraintName] = props
}

func (boModel *businessObjectModel) getInDB() *DB {
	return boModel.inDB
}

func (boModel *businessObjectModel) isNotPersisted() bool {
	return boModel.inNoDB
}

func (boModel *businessObjectModel) isPersisted() bool {
	return !boModel.isNotPersisted()
}

func (boModel *businessObjectModel) ID() IField {
	return boModel.idField
}

func (boModel *businessObjectModel) getTableName(withSchema bool) string {
	if boModel.tableName == "" {
		boModel.tableName = core.PascalToSnake(string(boModel.name))
	}

	if withSchema {
		return string(boModel.inDB.config.Name) + "." + boModel.tableName
	}

	return boModel.tableName
}

func (boModel *businessObjectModel) base() *businessObjectModel {
	return boModel
}

func (boModel *businessObjectModel) addField(field IField) IField {
	boModel.fields[field.GetName()] = field

	return field
}

func (boModel *businessObjectModel) getType() *utils.GoaldType {
	if boModel.boType == nil {
		boType := utils.TypeOf(getClass(boModel).NewObject(), true)
		boModel.boType = &boType
	}

	return boModel.boType
}

// ------------------------------------------------------------------------------------------------
// utils: defining the property types Goald supports, and the way to detect them with reflection
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
	typeTIMExPTR = utils.TypeOf((*time.Time)(nil), false)
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
func detectPropertyType(field utils.GoaldField, iBopropertyType, enumpropertyType utils.GoaldType) (propertyType propertyType, multiple bool) {
	// a business object's real property must be exported, and therefore PkgPath should be empty
	// Cf. https://golang.org/pkg/reflect/#StructField
	if fieldType := field.Type(); field.PkgPath() == "" {
		// getting the field kind
		fieldKind := fieldType.Kind()

		// handling the case where we have a slice in here
		if fieldKind == utils.KindSLICE {
			// what's in there?
			innerSliceType := fieldType.Elem()
			innerSliceKind := innerSliceType.Kind()

			// detecting an enum
			if innerSliceType.Implements(enumpropertyType) {
				return propertyTypeENUM, true
			}

			// detecting a polymorphic type, i.e. an interface; this should point to something implementing IBusinessObject
			if innerSliceKind == utils.KindINTERFACE && innerSliceType.Implements(iBopropertyType) {
				return propertyTypeRELATIONSHIPxPOLYM, true
			}

			// detecting a single relationship to a business object
			if innerSliceKind == utils.KindPTR && innerSliceType.Implements(iBopropertyType) {
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
			case utils.KindBOOL:
				return propertyTypeBOOL, false

			case utils.KindSTRING:
				return propertyTypeSTRING, false

			case utils.KindINT:
				return propertyTypeINT, false

			case utils.KindINT64:
				return propertyTypeBIGINT, false

			case utils.KindFLOAT32:
				return propertyTypeREAL, false

			case utils.KindFLOAT64:
				return propertyTypeDOUBLE, false
			}

			// detecting a polymorphic type, i.e. an interface; this should point to something implementing IBusinessObject
			if fieldKind == utils.KindINTERFACE && fieldType.Implements(iBopropertyType) {
				return propertyTypeRELATIONSHIPxPOLYM, false
			}

			// detecting a single relationship to a business object
			if fieldKind == utils.KindPTR && fieldType.Implements(iBopropertyType) {
				return propertyTypeRELATIONSHIPxMONOM, false
			}
		}
	}

	// detect if the field name starts with an uppercase letter, this should be allowed for non-handled property types
	if unicode.IsUpper(rune(field.Name()[0])) {
		panic(fmt.Sprintf("This exported property ('%s') should have a type handled by Goald. "+
			"If this is a relationship, make sure you're using a pointer to it: *somepackage.SomeClass", field.Name()))
	}

	// this happens with technical fields !
	return propertyTypeUNKNOWN, false
}

// ------------------------------------------------------------------------------------------------
// Business object properties, whether fields or relationships
// ------------------------------------------------------------------------------------------------
type IBusinessObjectProperty interface {
	ownerModel() IBusinessObjectModel
	setOwner(IBusinessObjectModel)
	GetName() string
	getPropertyType() propertyType
	IsMultiple() bool
	getColumnName() string
	isNotPersisted() bool
	getStructField() *utils.GoaldField
	getTag(tagName string) string
	isMandatoryInput() bool
	isPureOutput() bool
	IsRequiredInDb() bool
	SetRequiredInDb()
	isUnique() bool
	SetUnique()
	getUniqueConstraintName() string
}

type ioType string

const ioTypeINPUT ioType = "in"
const ioTypeINPUTxMANDATORY ioType = "i*"
const ioTypePURExOUTPUT ioType = "o*"

type businessObjectProperty struct {
	owner        IBusinessObjectModel // the property's owner class
	name         string               // the property's name, as declared in the struct
	propType     propertyType         // the property's type, as detected by the codegen phase
	multiple     bool                 // the property's multiplicity; false = 1, true = N
	columnName   string               // if this property - field or relationship - is persisted on the owner's table
	notPersisted bool                 // if true, then this property does not have a corresponding column in the BO's table
	structField  *utils.GoaldField    // the struct field corresponding to this property, as detected by the codegen phase
	requiredInDb bool                 // if true, then this property is required in the database (NOT NULL)
	unique       bool                 // if true, then this property is unique in the database (UNIQUE)
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
	}

	return prop.columnName
}

func (prop *businessObjectProperty) IsMultiple() bool {
	return prop.multiple
}

func (prop *businessObjectProperty) isNotPersisted() bool {
	return prop.notPersisted
}

func (prop *businessObjectProperty) getStructField() *utils.GoaldField {
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
