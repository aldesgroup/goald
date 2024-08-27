// ------------------------------------------------------------------------------------------------
// The code here is about describing the business objects through "classes"
// ------------------------------------------------------------------------------------------------
package goald

import (
	"github.com/aldesgroup/goald/features/utils"
)

// TODO endpoints should be plural
// TODO pagination all the way
// TODO generate stuff for enums ?

// ------------------------------------------------------------------------------------------------
// Business object classes
// ------------------------------------------------------------------------------------------------
type IBusinessObjectClass interface {
	/* public generic methods */

	SetNotPersisted() // to indicate this class has no instance persisted in a database
	SetInDB(db *DB)   // to associate the class with the DB where its instances are stored

	// access to generic properties (fields & relationships)
	ID() IField

	// private methods
	isNotPersisted() bool
	getInDB() *DB
	getTableName() string

	// access to the base implementation
	base() *businessObjectClass
	addField(field IField) IField
}

type className string

type businessObjectClass struct {
	name                    className                 // this class' name
	fields                  map[string]IField         // the objet's simple properties
	relationships           map[string]*Relationship  // the relationships to other classes
	inDB                    *DB                       // the associated DB, if any
	inNoDB                  bool                      // if true, then no associated DB
	tableName               string                    // if persisted, the name of the corresponding DB table - should be the same as the class name most of the time
	persistedProperties     []iBusinessObjectProperty // all the properties - fields or relationships - persisted on this class
	relationshipsWithColumn []*Relationship           // all the relationships for which this class has a column in its table
	idField                 IField                    // accessor to the ID field
	// allProperties           []iBusinessObjectProperty // all the properties - fields or relationships
}

func NewClass() IBusinessObjectClass {
	class := &businessObjectClass{
		fields:        map[string]IField{},
		relationships: map[string]*Relationship{},
	}

	// adding the generic fields
	class.idField = NewBigIntField(class, "ID", false)

	return class
}

func (boClass *businessObjectClass) SetInDB(db *DB) {
	boClass.inNoDB = false
	boClass.inDB = db
}

func (boClass *businessObjectClass) SetNotPersisted() {
	boClass.inNoDB = true
	boClass.inDB = nil
}

func (boClass *businessObjectClass) getInDB() *DB {
	return boClass.inDB
}

func (boClass *businessObjectClass) isNotPersisted() bool {
	return boClass.inNoDB
}

func (boClass *businessObjectClass) isPersisted() bool {
	return !boClass.isNotPersisted()
}

func (boClass *businessObjectClass) ID() IField {
	return boClass.idField
}

func (boClass *businessObjectClass) getTableName() string {
	if boClass.tableName == "" {
		boClass.tableName = utils.PascalToSnake(string(boClass.name))
	}

	return boClass.tableName
}

func (boClass *businessObjectClass) base() *businessObjectClass {
	return boClass
}

func (boClass *businessObjectClass) addField(field IField) IField {
	boClass.fields[field.getName()] = field

	return field
}

// ------------------------------------------------------------------------------------------------
// Business object properties, wether fields or relationships
// ------------------------------------------------------------------------------------------------
type iBusinessObjectProperty interface {
	ownerClass() IBusinessObjectClass
	setOwner(IBusinessObjectClass)
	getName() string
	getTypeFamily() utils.TypeFamily
	isMultiple() bool
	getColumnName() string
	isMandatory() bool
	isNotPersisted() bool
}

type businessObjectProperty struct {
	owner        IBusinessObjectClass // the property's owner class
	name         string               // the property's name, as declared in the struct
	typeFamily   utils.TypeFamily     // the property's type, as detected by the codegen phase
	multiple     bool                 // the property's multiplicity; false = 1, true = N
	columnName   string               // if this property - field or relationship - is persisted on the owner's table
	mandatory    bool                 // if true, then this property's value must be non-zero
	notPersisted bool                 // if true, then this property does not have a corresponding column in the BO's table
}

func (prop *businessObjectProperty) ownerClass() IBusinessObjectClass {
	return prop.owner
}

func (prop *businessObjectProperty) setOwner(owner IBusinessObjectClass) {
	prop.owner = owner
}

func (prop *businessObjectProperty) getName() string {
	return prop.name
}

func (prop *businessObjectProperty) getTypeFamily() utils.TypeFamily {
	return prop.typeFamily
}

func (prop *businessObjectProperty) getColumnName() string {
	if prop.columnName == "" {
		prop.columnName = utils.PascalToSnake(prop.name)
		if prop.typeFamily == utils.TypeFamilyRELATIONSHIP {
			prop.columnName += "_id"
		}
	}

	return prop.columnName
}

func (prop *businessObjectProperty) isMultiple() bool {
	return prop.multiple
}

func (prop *businessObjectProperty) isMandatory() bool {
	return prop.mandatory
}

func (prop *businessObjectProperty) isNotPersisted() bool {
	return prop.notPersisted
}

// ------------------------------------------------------------------------------------------------
// Fields (simple properties) of business object classes
// ------------------------------------------------------------------------------------------------

type IField interface {
	iBusinessObjectProperty
	isBuiltIn() bool
	getDefaultValue() string
	// SetDefaultValue(string) IField
}

// base implementation
type field struct {
	businessObjectProperty
	defaultStringValue string
}

func newField(owner IBusinessObjectClass, name string, multiple bool, typeFamily utils.TypeFamily) field {
	return field{
		businessObjectProperty: businessObjectProperty{
			owner:      owner,
			name:       name,
			typeFamily: typeFamily,
			multiple:   multiple,
		},
	}
}

func (f *field) SetMandatory() *field {
	f.mandatory = true
	return f
}

func (f *field) SetNotPersisted() *field {
	f.notPersisted = true
	return f
}

func (f *field) isBuiltIn() bool {
	return false
}

func (f *field) SetDefaultValue(val string) *field {
	f.defaultStringValue = val
	return f
}

func (f *field) getDefaultValue() string {
	return f.defaultStringValue
}

type BoolField struct {
	field
}

type StringField struct {
	field
	size int
}

func (sf *StringField) SetSize(size int) *field {
	sf.size = size
	return &sf.field
}

type IntField struct {
	field
}

type BigIntField struct {
	field
}

type RealField struct {
	field
	// digits   int // number of digits before the decimal points, e.g. 4 in 9876.06
	// decimals int // number of digits after the decimal points, e.g. 2 in 9876.06
}

type DoubleField struct {
	field
	// digits   int // number of digits before the decimal points, e.g. 4 in 9876.06
	// decimals int // number of digits after the decimal points, e.g. 2 in 9876.06
}

// func (rf *realField) SetForma/home/jwan/Git/emeraldrt(digits, decimals int) *realField {
// 	rf.digits = digits
// 	rf.decimals = decimals
// 	return rf
// }

type DateField struct {
	field
}

type EnumField struct {
	field
}

func NewBoolField(owner IBusinessObjectClass, name string, multiple bool) *BoolField {
	return owner.addField(&BoolField{
		field: newField(owner, name, multiple, utils.TypeFamilyBOOL),
	}).(*BoolField)
}

func NewStringField(owner IBusinessObjectClass, name string, multiple bool) *StringField {
	return owner.addField(&StringField{
		field: newField(owner, name, multiple, utils.TypeFamilySTRING),
	}).(*StringField)
}

func NewIntField(owner IBusinessObjectClass, name string, multiple bool) *IntField {
	return owner.addField(&IntField{
		field: newField(owner, name, multiple, utils.TypeFamilyINT),
	}).(*IntField)
}

func NewBigIntField(owner IBusinessObjectClass, name string, multiple bool) *BigIntField {
	return owner.addField(&BigIntField{
		field: newField(owner, name, multiple, utils.TypeFamilyBIGINT),
	}).(*BigIntField)
}

func NewRealField(owner IBusinessObjectClass, name string, multiple bool) *RealField {
	return owner.addField(&RealField{
		field: newField(owner, name, multiple, utils.TypeFamilyREAL),
	}).(*RealField)
}
func NewDoubleField(owner IBusinessObjectClass, name string, multiple bool) *DoubleField {
	return owner.addField(&DoubleField{
		field: newField(owner, name, multiple, utils.TypeFamilyDOUBLE),
	}).(*DoubleField)
}

func NewDateField(owner IBusinessObjectClass, name string, multiple bool) *DateField {
	return owner.addField(&DateField{
		field: newField(owner, name, multiple, utils.TypeFamilyDATE),
	}).(*DateField)
}

func NewEnumField(owner IBusinessObjectClass, name string, multiple bool) *EnumField {
	return owner.addField(&EnumField{
		field: newField(owner, name, multiple, utils.TypeFamilyENUM),
	}).(*EnumField)
}

// ------------------------------------------------------------------------------------------------
// Relationships with other business object classes
// ------------------------------------------------------------------------------------------------

// relationshipType is used to define the type of the relationship between 2 classes
type relationshipType int

const (
	// relationshipTypeONExWAY : the entity owning the link is pointing to a target entity
	// There's no backref in this case, but there is in all other cases
	relationshipTypeONExWAY relationshipType = 1 + iota

	// relationshipTypeSOURCExTOxTARGET : the entity owning the link is pointing to a target entity, retaining its ID in DB
	relationshipTypeSOURCExTOxTARGET

	// relationshipTypeTARGETxTOxSOURCE : the entity owning the link is pointed by another entity, from another table
	relationshipTypeTARGETxTOxSOURCE

	// relationshipTypePARENTxTOxCHILDREN : the entity owning the link is pointed by children entities
	relationshipTypePARENTxTOxCHILDREN

	// relationshipTypeCHILDxTOxPARENT : the entity owning the link points to a parent entity
	relationshipTypeCHILDxTOxPARENT
)

type Relationship struct {
	businessObjectProperty
	target       IBusinessObjectClass // the type of BO pointed by this relationship
	relationType relationshipType     // valued from the business object's init
	backRef      *Relationship        // valued from the business object's init
}

// Allows to declare a new relationship on a given class
func NewRelationship(owner IBusinessObjectClass, name string, multiple bool, target IBusinessObjectClass) *Relationship {
	relationship := &Relationship{
		businessObjectProperty: businessObjectProperty{
			owner:    owner,
			name:     name,
			multiple: multiple,
		},
		target: target,
	}

	// owner.addRelationship(name, relationship)
	owner.base().relationships[name] = relationship

	return relationship
}

// Sets a relationship as a "child to parent" one; the backref relationship is needed
func (r *Relationship) SetChildToParent(backRefRelation *Relationship) *Relationship {
	r.relationType = relationshipTypeCHILDxTOxPARENT
	r.backRef = backRefRelation

	// automatically setting on the backref the inverse relation type and this relationship as the backref
	backRefRelation.relationType = relationshipTypePARENTxTOxCHILDREN
	backRefRelation.backRef = r

	return r
}

// Sets a relationship as a "parent to children" one; the backref relationship is needed
func (r *Relationship) SetSourceToTarget(backRefRelation *Relationship) *Relationship {
	r.relationType = relationshipTypeSOURCExTOxTARGET
	r.backRef = backRefRelation

	// automatically setting on the backref the inverse relation type and this relationship as the backref
	backRefRelation.relationType = relationshipTypeTARGETxTOxSOURCE
	backRefRelation.backRef = r

	return r
}

// Sets a relationship as a "one way" one; is with no back ref
func (r *Relationship) SetOneWay() *Relationship {
	r.relationType = relationshipTypeONExWAY

	return r
}

// returns true if this relationships, should it be persisted, needs a column on its owner's table for it
func (r *Relationship) needsColumn() bool {
	if r.multiple {
		return false
	}

	return r.relationType == relationshipTypeSOURCExTOxTARGET ||
		r.relationType == relationshipTypeCHILDxTOxPARENT ||
		r.relationType == relationshipTypeONExWAY
}
