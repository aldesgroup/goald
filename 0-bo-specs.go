// ------------------------------------------------------------------------------------------------
// The code here is about having tools to describe the properties (fields & relationships) of a
// business object class, and specifying things for these properties
// ------------------------------------------------------------------------------------------------
package goald

import (
	"sync"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// TODO endpoints should be plural
// TODO pagination all the way
// TODO generate stuff for enums ?

// ------------------------------------------------------------------------------------------------
// Specs for a business object classes
// ------------------------------------------------------------------------------------------------
type IBusinessObjectSpecs interface {
	/* public generic methods */

	SetNotPersisted() // to indicate this class has no instance persisted in a database
	SetInDB(db *DB)   // to associate the class with the DB where its instances are stored
	SetAbstract()     // to indicate this class does not model concrete business objects, but most probably a super class

	// access to generic properties (fields & relationships)
	ID() IField

	// private methods
	isNotPersisted() bool
	getInDB() *DB
	getTableName() string

	// access to the base implementation
	base() *businessObjectSpecs
	addField(field IField) IField
}

type className string

type businessObjectSpecs struct {
	name                    className                 // the corresponding class name
	fields                  map[string]IField         // the objet's simple properties
	relationships           map[string]*Relationship  // the relationships to other classes
	inDB                    *DB                       // the associated DB, if any
	inNoDB                  bool                      // if true, then no associated DB
	abstract                bool                      // if true, then is class is mainly used as a super class for others
	tableName               string                    // if persisted, the name of the corresponding DB table - should be the same as the class name most of the time
	persistedProperties     []iBusinessObjectProperty // all the properties - fields or relationships - persisted on this class
	relationshipsWithColumn []*Relationship           // all the relationships for which this class has a column in its table
	idField                 IField                    // accessor to the ID field
	usedInNativeApp         bool                      // true if this class is used in the native app
	usedInWebApp            bool                      // true if this class is used in the web app
}

func NewBusinessObjectSpecs() IBusinessObjectSpecs {
	specs := &businessObjectSpecs{
		fields:        map[string]IField{},
		relationships: map[string]*Relationship{},
	}

	// adding the generic fields
	specs.idField = NewBigIntField(specs, "ID", false)

	return specs
}

func (boClass *businessObjectSpecs) SetInDB(db *DB) {
	boClass.inNoDB = false
	boClass.inDB = db
}

func (boClass *businessObjectSpecs) SetNotPersisted() {
	boClass.inNoDB = true
	boClass.inDB = nil
}

func (boClass *businessObjectSpecs) SetAbstract() {
	boClass.abstract = true
}

func (boClass *businessObjectSpecs) getInDB() *DB {
	return boClass.inDB
}

func (boClass *businessObjectSpecs) isNotPersisted() bool {
	return boClass.inNoDB
}

func (boClass *businessObjectSpecs) isPersisted() bool {
	return !boClass.isNotPersisted()
}

func (boClass *businessObjectSpecs) ID() IField {
	return boClass.idField
}

func (boClass *businessObjectSpecs) getTableName() string {
	if boClass.tableName == "" {
		boClass.tableName = core.PascalToSnake(string(boClass.name))
	}

	return boClass.tableName
}

func (boClass *businessObjectSpecs) base() *businessObjectSpecs {
	return boClass
}

func (boClass *businessObjectSpecs) addField(field IField) IField {
	boClass.fields[field.getName()] = field

	return field
}

// ------------------------------------------------------------------------------------------------
// Business object properties, whether fields or relationships
// ------------------------------------------------------------------------------------------------
type iBusinessObjectProperty interface {
	ownerSpecs() IBusinessObjectSpecs
	setOwner(IBusinessObjectSpecs)
	getName() string
	getTypeFamily() utils.TypeFamily
	isMultiple() bool
	getColumnName() string
	isMandatory() bool
	isNotPersisted() bool
}

type businessObjectProperty struct {
	owner        IBusinessObjectSpecs // the property's owner class
	name         string               // the property's name, as declared in the struct
	typeFamily   utils.TypeFamily     // the property's type, as detected by the codegen phase
	multiple     bool                 // the property's multiplicity; false = 1, true = N
	columnName   string               // if this property - field or relationship - is persisted on the owner's table
	mandatory    bool                 // if true, then this property's value must be non-zero
	notPersisted bool                 // if true, then this property does not have a corresponding column in the BO's table
}

func (prop *businessObjectProperty) ownerSpecs() IBusinessObjectSpecs {
	return prop.owner
}

func (prop *businessObjectProperty) setOwner(owner IBusinessObjectSpecs) {
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
		prop.columnName = core.PascalToSnake(prop.name)
		if prop.typeFamily == utils.TypeFamilyRELATIONSHIPxMONOM {
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

type iNumericField interface {
	IField
	isMinSet() bool
	isMaxSet() bool
}

// base implementation
type field struct {
	businessObjectProperty
	defaultStringValue string
}

type numericField struct {
	field
	minSet bool
	maxSet bool
}

func newField(owner IBusinessObjectSpecs, name string, multiple bool, typeFamily utils.TypeFamily) field {
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

func (f *numericField) isMinSet() bool {
	return f.minSet
}

func (f *numericField) isMaxSet() bool {
	return f.maxSet
}

type BoolField struct {
	field
}

type StringField struct {
	field
	size    int
	atLeast int
}

func (sf *StringField) SetSize(size int, atLeast ...int) *field {
	sf.size = size
	if len(atLeast) > 0 {
		sf.atLeast = atLeast[0]
	}
	return &sf.field
}

type IntField struct {
	numericField
	min int
	max int
}

func (f *IntField) Min(min int) *IntField {
	f.min = min
	f.minSet = true
	return f
}

func (f *IntField) Max(max int) *IntField {
	f.max = max
	f.maxSet = true
	return f
}

type BigIntField struct {
	numericField
	min int64
	max int64
}

func (f *BigIntField) Min(min int64) *BigIntField {
	f.min = min
	f.minSet = true
	return f
}

func (f *BigIntField) Max(max int64) *BigIntField {
	f.max = max
	return f
}

type RealField struct {
	numericField
	min float32
	max float32
	// digits   int // number of digits before the decimal points, e.g. 4 in 9876.06
	// decimals int // number of digits after the decimal points, e.g. 2 in 9876.06
}

func (f *RealField) Min(min float32) *RealField {
	f.min = min
	f.minSet = true
	return f
}

func (f *RealField) Max(max float32) *RealField {
	f.max = max
	return f
}

type DoubleField struct {
	numericField
	min float64
	max float64
	// digits   int // number of digits before the decimal points, e.g. 4 in 9876.06
	// decimals int // number of digits after the decimal points, e.g. 2 in 9876.06
}

func (f *DoubleField) Min(min float64) *DoubleField {
	f.min = min
	f.minSet = true
	return f
}

func (f *DoubleField) Max(max float64) *DoubleField {
	f.max = max
	return f
}

// func (rf *realField) SetFormat(digits, decimals int) *realField {
// 	rf.digits = digits
// 	rf.decimals = decimals
// 	return rf
// }

type DateField struct {
	field
}

type EnumField struct {
	field
	enumName   string
	onlyValues []IEnum
}

func (f *EnumField) Only(values ...IEnum) *EnumField {
	f.onlyValues = values
	return f
}

func NewBoolField(owner IBusinessObjectSpecs, name string, multiple bool) *BoolField {
	return owner.addField(&BoolField{
		field: newField(owner, name, multiple, utils.TypeFamilyBOOL),
	}).(*BoolField)
}

func NewStringField(owner IBusinessObjectSpecs, name string, multiple bool) *StringField {
	return owner.addField(&StringField{
		field: newField(owner, name, multiple, utils.TypeFamilySTRING),
	}).(*StringField)
}

func NewIntField(owner IBusinessObjectSpecs, name string, multiple bool) *IntField {
	return owner.addField(&IntField{numericField: numericField{
		field: newField(owner, name, multiple, utils.TypeFamilyINT),
	}}).(*IntField)
}

func NewBigIntField(owner IBusinessObjectSpecs, name string, multiple bool) *BigIntField {
	return owner.addField(&BigIntField{numericField: numericField{
		field: newField(owner, name, multiple, utils.TypeFamilyBIGINT),
	}}).(*BigIntField)
}

func NewRealField(owner IBusinessObjectSpecs, name string, multiple bool) *RealField {
	return owner.addField(&RealField{numericField: numericField{
		field: newField(owner, name, multiple, utils.TypeFamilyREAL),
	}}).(*RealField)
}
func NewDoubleField(owner IBusinessObjectSpecs, name string, multiple bool) *DoubleField {
	return owner.addField(&DoubleField{numericField: numericField{
		field: newField(owner, name, multiple, utils.TypeFamilyDOUBLE),
	}}).(*DoubleField)
}

func NewDateField(owner IBusinessObjectSpecs, name string, multiple bool) *DateField {
	return owner.addField(&DateField{
		field: newField(owner, name, multiple, utils.TypeFamilyDATE),
	}).(*DateField)
}

func NewEnumField(owner IBusinessObjectSpecs, name string, multiple bool, enumName string) *EnumField {
	return owner.addField(&EnumField{
		field:    newField(owner, name, multiple, utils.TypeFamilyENUM),
		enumName: enumName,
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
	targets      []IBusinessObjectSpecs // the type of BO pointed by this relationship
	relationType relationshipType       // valued from the business object's init
	backRefs     []*Relationship        // valued from the business object's init
	polymorphic  bool                   // if true, then it's a polymorphic relationship
	mx           sync.Mutex             // a mutex for the operations on the slices in here
}

// Allows to declare a new relationship on a given class
func NewRelationship(owner IBusinessObjectSpecs, name string, multiple bool, targets ...IBusinessObjectSpecs) *Relationship {
	relationship := &Relationship{
		businessObjectProperty: businessObjectProperty{
			owner:    owner,
			name:     name,
			multiple: multiple,
		},
		targets:     targets,
		polymorphic: len(targets) > 1,
	}

	owner.base().relationships[name] = relationship

	return relationship
}

func (r *Relationship) addBackRef(backRef *Relationship) {
	r.mx.Lock()
	if r.polymorphic || len(r.backRefs) == 0 {
		r.backRefs = append(r.backRefs, backRef)
	}
	r.mx.Unlock()
}

// Sets a relationship as a "child to parent" one; the backref relationship is needed
func (r *Relationship) SetChildToParent(backRefRelation *Relationship) *Relationship {
	r.relationType = relationshipTypeCHILDxTOxPARENT

	// taking the opportunity here to enrich the backref relationship...
	r.addBackRef(backRefRelation)

	// ... like automatically setting on the backref the inverse relation type and this relationship as the backref
	backRefRelation.relationType = relationshipTypePARENTxTOxCHILDREN
	backRefRelation.addBackRef(r)

	return r
}

// Sets a relationship as a "parent to children" one; the backref relationship is needed
func (r *Relationship) SetSourceToTarget(backRefRelation *Relationship) *Relationship {
	r.relationType = relationshipTypeSOURCExTOxTARGET

	// taking the opportunity here to enrich the backref relationship...
	r.addBackRef(backRefRelation)

	// automatically setting on the backref the inverse relation type and this relationship as the backref
	backRefRelation.relationType = relationshipTypeTARGETxTOxSOURCE
	backRefRelation.addBackRef(r)

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
