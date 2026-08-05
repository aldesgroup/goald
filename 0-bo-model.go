// ------------------------------------------------------------------------------------------------
// The code here is about having tools to describe business object models, and
// their properties in a general sense
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"sort"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/reflection"
	"github.com/aldesgroup/goald/features/utils"
)

// TODO pagination all the way

// ------------------------------------------------------------------------------------------------
// Model for a business object model
// ------------------------------------------------------------------------------------------------
type IBusinessObjectModel interface {
	// public generic methods
	SetDescription(description string)                     // to set the description of the model
	SetNotPersisted()                                      // to indicate this model has no instance persisted in a database
	SetInDB(db *DB)                                        // to associate the model with the DB where its instances are stored
	SetInDbByName(dbName string)                           // to associate the model with the DB where its instances are stored, using the name of the DB instead of its instance
	SetAbstract()                                          // to indicate this model does not model concrete business objects, but most probably a super model
	AddUniqueCombination(props ...IBusinessObjectProperty) // to indicate that a combination of properties must be unique in the DB
	SetAutoCRUD()                                          // to automatically start the generic CRUD endpoints for this business object model

	// // access to generic properties (fields & relationships)
	// ID() IField

	// private methods
	setName(name utils.ModelName)
	isAbstract() bool
	isPersistedHere() bool
	isNotPersisted() bool
	getDB() *DB
	isPersisted() bool
	getTableName(withSchema bool) string
	getDescription() string
	getUniqueCombinations() map[string][]IBusinessObjectProperty
	getFields() map[string]IField
	getField(name string) IField
	addField(field IField) IField
	getRelationships() map[string]*Relationship
	getRelationship(name string) *Relationship
	setRelationship(name string, relationship *Relationship)
	setChildToParentRelationship(r *Relationship)
	isUsedInNativeApp() bool
	setUsedInNativeApp()
	isUsedInWebApp() bool
	setUsedInWebApp()
	getIdField() *BigIntField

	// utils
	getPersistedProperties() []IBusinessObjectProperty
	getRelationshipsWithColumn() []*Relationship
	getRelationshipsWithLinkTable() []*Relationship
	getRelationshipsWithRequiredBackref() []*Relationship

	// technical stuff
	resolve() IBusinessObjectModel
	getType() *reflection.GoaldType

	// this makes each model aware of its source, and capable of instantiating new business objects of its type
	IBusinessObjectModelSource
}

// ------------------------------------------------------------------------------------------------
// Generic implementation with public methods
// ------------------------------------------------------------------------------------------------

type businessObjectModel struct {
	name                             utils.ModelName                      // the corresponding model name
	source                           IBusinessObjectModelSource           // the corresponding source object
	description                      string                               // the model description
	fields                           map[string]IField                    // the objet's simple properties
	relationships                    map[string]*Relationship             // the relationships to other models
	db                               *DB                                  // the associated DB, if any
	inNoDB                           bool                                 // if true, then no associated DB
	abstract                         bool                                 // if true, then this model is mainly used as a super model for others
	tableName                        string                               // if persisted, the name of the corresponding DB table - should be the same as the model name most of the time
	persistedProperties              []IBusinessObjectProperty            // all the properties - fields or relationships - persisted on this model
	uniqueCombinations               map[string][]IBusinessObjectProperty // a combination of properties that must be unique in the DB
	idField                          *BigIntField                         // accessor to the ID field
	usedInNativeApp                  bool                                 // true if this model is used in the native app
	usedInWebApp                     bool                                 // true if this model is used in the web app
	boType                           *reflection.GoaldType                // the Go type associated with this BO model
	relationshipsWithColumn          []*Relationship                      // all the relationships for which this model has a column in its table
	relationshipsWithLinkTable       []*Relationship                      // all the relationships for which this model has a column in its table
	resolved                         bool                                 // if true, then this model has been resolved, i.e. all its properties have been detected and registered
	autoCRUD                         bool                                 // if true, then the generic CRUD endpoints will be automatically started for this business object model
	childToParentRelationship        *Relationship                        // if this model is a child in a parent-child relationship, then this is the relationship to the parent
	relationshipsWithRequiredBackref []*Relationship                      // all the relationships through which the target BOs cannot be persisted without a backref to this BO
}

const BoFieldID = "ID"       // the name of the ID field, which is a special case in Goald
const boFieldPreID = "preID" // the name of the row ID field, which is a special case in Goald

func NewBusinessObjectModel() IBusinessObjectModel {
	model := &businessObjectModel{
		fields:        map[string]IField{},
		relationships: map[string]*Relationship{},
	}

	// adding the generic fields
	model.idField = AddBigIntField(model, "BusinessObject", BoFieldID, false)
	preIDField := AddIntField(model, "BusinessObject", boFieldPreID, false)
	creationField := AddDateField(model, "BusinessObject", "Creation", false)

	// some tweaking
	preIDField.setTechnical()
	creationField.SetRequiredInDb()

	return model
}

func (boModel *businessObjectModel) SetInDB(db *DB) {
	boModel.inNoDB = false
	boModel.db = db
}

func (boModel *businessObjectModel) SetInDbByName(dbName string) {
	boModel.SetInDB(GetDB(dbconn.DbSchemaName(dbName)))
}

func (boModel *businessObjectModel) SetDescription(description string) {
	boModel.description = description
}

func (boModel *businessObjectModel) SetNotPersisted() {
	boModel.inNoDB = true
	boModel.db = nil
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

func (boModel *businessObjectModel) SetAutoCRUD() {
	boModel.autoCRUD = true
}

// ------------------------------------------------------------------------------------------------
// Generic implementation with private methods
// ------------------------------------------------------------------------------------------------

func (boModel *businessObjectModel) GetName() utils.ModelName {
	return boModel.name
}

func (boModel *businessObjectModel) setName(name utils.ModelName) {
	boModel.name = name
}

func (boModel *businessObjectModel) isAbstract() bool {
	return boModel.abstract
}

func (boModel *businessObjectModel) isNotPersisted() bool {
	return boModel.inNoDB
}

func (boModel *businessObjectModel) getDB() *DB {
	return boModel.db
}

func (boModel *businessObjectModel) isPersisted() bool {
	return !boModel.isNotPersisted()
}

func (boModel *businessObjectModel) isPersistedHere() bool {
	// a bit of control here
	if boModel.isPersisted() && boModel.db == nil {
		core.PanicMsg("Model '%s' should be SetNotPersisted, SetAbstract, or associated with a DB", boModel.GetName())
	}

	return boModel.isPersisted() && boModel.db.schema != nil
}

func (boModel *businessObjectModel) getTableName(withSchema bool) string {
	if boModel.tableName == "" {
		boModel.tableName = core.PascalToSnake(string(boModel.name))
	}

	if withSchema {
		return string(boModel.db.name) + "." + boModel.tableName
	}

	return boModel.tableName
}

func (boModel *businessObjectModel) getDescription() string {
	return boModel.description
}

func (boModel *businessObjectModel) getUniqueCombinations() map[string][]IBusinessObjectProperty {
	return boModel.uniqueCombinations
}

func (boModel *businessObjectModel) getFields() map[string]IField {
	return boModel.fields
}

func (boModel *businessObjectModel) getField(name string) IField {
	return boModel.fields[name]
}

func (boModel *businessObjectModel) addField(field IField) IField {
	boModel.fields[field.GetName()] = field

	return field
}

func (boModel *businessObjectModel) getRelationships() map[string]*Relationship {
	return boModel.relationships
}

func (boModel *businessObjectModel) getRelationship(name string) *Relationship {
	return boModel.relationships[name]
}

func (boModel *businessObjectModel) setRelationship(name string, relationship *Relationship) {
	boModel.relationships[name] = relationship
}

func (boModel *businessObjectModel) setChildToParentRelationship(rel *Relationship) {
	if boModel.childToParentRelationship != nil {
		panic(fmt.Sprintf("Business object model '%s' already has a parent relationship '%s', cannot set another one '%s'",
			boModel.name, boModel.childToParentRelationship.GetName(), rel.GetName()))
	}

	boModel.childToParentRelationship = rel
}

func (boModel *businessObjectModel) isUsedInNativeApp() bool {
	return boModel.usedInNativeApp
}

func (boModel *businessObjectModel) setUsedInNativeApp() {
	boModel.usedInNativeApp = true
}

func (boModel *businessObjectModel) isUsedInWebApp() bool {
	return boModel.usedInWebApp
}

func (boModel *businessObjectModel) setUsedInWebApp() {
	boModel.usedInWebApp = true
}

func (boModel *businessObjectModel) getIdField() *BigIntField {
	return boModel.idField
}

// ------------------------------------------------------------------------------------------------
// Technical stuff
// ------------------------------------------------------------------------------------------------

func (boModel *businessObjectModel) getType() *reflection.GoaldType {
	if boModel.boType == nil {
		boType := reflection.TypeOf(boModel.NewObject(), true)
		boModel.boType = &boType
	}

	return boModel.boType
}

// ------------------------------------------------------------------------------------------------
// Business Object Model resolution =
// - allowing some kind of "inheritance" between models
// ------------------------------------------------------------------------------------------------

func (boModel *businessObjectModel) resolve() IBusinessObjectModel {
	// if the model hasn't been resolved yet
	if boModel == nil || !boModel.resolved {

		// first thing first: let's connect the model to its source, if any
		boModel.source = sourceRegistry.items[boModel.name]

		// checking all the fields, and "importing" the configuration of the super models, if any
		for _, field := range boModel.fields {
			if field.getDeclaringBO() != boModel.name {
				// making sure the "parent" model is resolved first
				if parentModel := modelFor(field.getDeclaringBO()); parentModel != nil {
					// getting the field from the parent model
					parentField := parentModel.resolve().getField(field.GetName())

					// importing the configuration of the super model's property
					boModel.importConfig(parentField, field)
				}
			}
		}

		// checking all the relationships, and "importing" the configuration of the super models, if any
		for _, rel := range boModel.relationships {
			if rel.getDeclaringBO() != boModel.name {
				// making sure the "parent" model is resolved first
				if parentModel := modelFor(rel.getDeclaringBO()); parentModel != nil {
					// getting the relationship from the parent model
					parentRel := parentModel.resolve().getRelationship(rel.GetName())

					// importing the configuration of the super model's property
					boModel.importConfig(parentRel, rel)
				}
			}
		}

		// TODO check if we have to do something with the unique combinations here

		// tagging this model as resolved
		boModel.resolved = true
	}

	return boModel
}

func (boModel *businessObjectModel) importConfig(from, to IBusinessObjectProperty) {
	// generic configuration
	if from.IsRequiredInDb() {
		to.SetRequiredInDb()
	}
	if from.isUnique() {
		to.SetUnique()
	}
	if from.isSecret() {
		to.SetSecret()
	}
	if from.isPersonal() {
		to.SetPersonal()
	}

	// specific configuration for fields
	switch field := to.(type) {
	case *StringField:
		if field.size == 0 {
			field.size = from.(*StringField).size
		}
	case *DoubleField:
		if field.totalDigits+field.decimals == 0 {
			fromField := from.(*DoubleField)
			field.totalDigits = fromField.totalDigits
			field.decimals = fromField.decimals
		}
	case *RealField:
		if field.totalDigits+field.decimals == 0 {
			fromField := from.(*RealField)
			field.totalDigits = fromField.totalDigits
			field.decimals = fromField.decimals
		}
	}

	// specific configuration for relationships
	if relationship, ok := to.(*Relationship); ok {
		if relationship.relationType == 0 {
			relationship.relationType = from.(*Relationship).relationType
		}
		if relationship.backRef == nil {
			relationship.backRef = from.(*Relationship).backRef
		}
	}
}
