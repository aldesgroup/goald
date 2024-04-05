// ------------------------------------------------------------------------------------------------
// The code here is about describing the business objects through "classes"
// ------------------------------------------------------------------------------------------------
package goald

import (
	"reflect"

	"github.com/aldesgroup/goald/features/utils"
)

// TODO endpoints should be plural
// TODO pagination all the way
// TODO generate stuff for enums ?

// ------------------------------------------------------------------------------------------------
// Business object classes
// ------------------------------------------------------------------------------------------------
type IBusinessObjectClass interface {
	// public generic methods
	GetInDB() *DB
	IsNotPersisted() bool
	SetInDB(db *DB)
	SetNotPersisted()

	// access to generic properties
	ID() IField

	// private methods
	getTableName() string

	// access to the base implementation
	base() *businessObjectClass
	addField(field IField) IField
}

type businessObjectClass struct {
	className               string                    // this class' name
	fields                  map[string]IField         // the objet's simple properties
	relationships           map[string]*Relationship  // the relationships to other classes
	inDB                    *DB                       // the associated DB, if any
	inNoDB                  bool                      // if true, then no associated DB
	tableName               string                    // if persisted, the name of the corresponding DB table
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
	class.idField = NewStringField(class, "ID", false)

	return class
}

func (boClass *businessObjectClass) GetInDB() *DB {
	return boClass.inDB
}

func (boClass *businessObjectClass) IsNotPersisted() bool {
	return boClass.inNoDB
}

func (boClass *businessObjectClass) isPersisted() bool {
	return !boClass.IsNotPersisted()
}

func (boClass *businessObjectClass) SetInDB(db *DB) {
	boClass.inNoDB = false
	boClass.inDB = db
}

func (boClass *businessObjectClass) SetNotPersisted() {
	boClass.inNoDB = true
	boClass.inDB = nil
}

func (boClass *businessObjectClass) ID() IField {
	return boClass.idField
}

func (boClass *businessObjectClass) getTableName() string {
	if boClass.tableName == "" {
		boClass.tableName = utils.PascalToSnake(boClass.className)
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
	getColumnName() string
}

type BusinessObjectProperty struct {
	owner      IBusinessObjectClass // the property's owner class
	name       string               // the property's name, as declared in the struct
	pType      PropertyType         // the property's type, as detected by the codegen phase
	multiple   bool                 // the property's multiplicity; false = 1, true = N
	columnName string               // if this property - field or relationship - is persisted on the owner's table
}

func (prop *BusinessObjectProperty) ownerClass() IBusinessObjectClass {
	return prop.owner
}

func (prop *BusinessObjectProperty) setOwner(owner IBusinessObjectClass) {
	prop.owner = owner
}

func (prop *BusinessObjectProperty) getName() string {
	return prop.name
}

func (prop *BusinessObjectProperty) getColumnName() string {
	if prop.columnName == "" {
		prop.columnName = utils.PascalToSnake(prop.name)
		if prop.pType == PropertyTypeRELATIONSHIP {
			prop.columnName += "_id"
		}
	}

	return prop.columnName
}

// ------------------------------------------------------------------------------------------------
// Fields (simple properties) of business object classes
// ------------------------------------------------------------------------------------------------

type IField interface {
	iBusinessObjectProperty
	StringValue(IBusinessObject) string
}

// base implementation
type field struct {
	BusinessObjectProperty
}

func newField(owner IBusinessObjectClass, name string, multiple bool, pType PropertyType) field {
	return field{
		BusinessObjectProperty: BusinessObjectProperty{
			owner:    owner,
			name:     name,
			pType:    pType,
			multiple: multiple,
		},
	}
}

func (f *field) StringValue(bObj IBusinessObject) string {
	// TODO we shall do better without reflect
	return reflect.ValueOf(bObj).Elem().FieldByName(f.name).String()
}

type BoolField struct {
	field
}

type StringField struct {
	field
	size int
}

func (sf *StringField) SetSize(size int) *StringField {
	sf.size = size
	return sf
}

type IntField struct {
	field
}

type Int64Field struct {
	field
}

type realField struct {
	field
	// digits   int // number of digits before the decimal points, e.g. 4 in 9876.06
	// decimals int // number of digits after the decimal points, e.g. 2 in 9876.06
}

// func (rf *realField) SetForma/home/jwan/Git/emeraldrt(digits, decimals int) *realField {
// 	rf.digits = digits
// 	rf.decimals = decimals
// 	return rf
// }

type Real32Field struct {
	realField
}

type Real64Field struct {
	realField
}

type DateField struct {
	field
}

type EnumField struct {
	field
}

func NewBoolField(owner IBusinessObjectClass, name string, multiple bool) *BoolField {
	return owner.addField(&BoolField{
		field: newField(owner, name, multiple, PropertyTypeBOOL),
	}).(*BoolField)
}

func NewStringField(owner IBusinessObjectClass, name string, multiple bool) *StringField {
	return owner.addField(&StringField{
		field: newField(owner, name, multiple, PropertyTypeSTRING),
	}).(*StringField)
}

func NewIntField(owner IBusinessObjectClass, name string, multiple bool) *IntField {
	return owner.addField(&IntField{
		field: newField(owner, name, multiple, PropertyTypeINT),
	}).(*IntField)
}

func NewInt64Field(owner IBusinessObjectClass, name string, multiple bool) *Int64Field {
	return owner.addField(&Int64Field{
		field: newField(owner, name, multiple, PropertyTypeINT64),
	}).(*Int64Field)
}

func NewReal32Field(owner IBusinessObjectClass, name string, multiple bool) *Real32Field {
	return owner.addField(&Real32Field{
		realField: realField{field: newField(owner, name, multiple, PropertyTypeREAL32)},
	}).(*Real32Field)
}

func NewReal64Field(owner IBusinessObjectClass, name string, multiple bool) *Real64Field {
	return owner.addField(&Real64Field{
		realField: realField{field: newField(owner, name, multiple, PropertyTypeREAL64)},
	}).(*Real64Field)
}

func NewDateField(owner IBusinessObjectClass, name string, multiple bool) *DateField {
	return owner.addField(&DateField{
		field: newField(owner, name, multiple, PropertyTypeDATE),
	}).(*DateField)
}

func NewEnumField(owner IBusinessObjectClass, name string, multiple bool) *EnumField {
	return owner.addField(&EnumField{
		field: newField(owner, name, multiple, PropertyTypeENUM),
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
	BusinessObjectProperty
	target       IBusinessObjectClass // the type of BO pointed by this relationship
	relationType relationshipType     // valued from the business object's init
	backRef      *Relationship        // valued from the business object's init
}

// Allows to declare a new relationship on a given class
func NewRelationship(owner IBusinessObjectClass, name string, multiple bool, target IBusinessObjectClass) *Relationship {
	relationship := &Relationship{
		BusinessObjectProperty: BusinessObjectProperty{
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
