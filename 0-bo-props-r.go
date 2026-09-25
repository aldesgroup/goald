package goald

import (
	"sync"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/reflection"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Relationships with other business object models
// ------------------------------------------------------------------------------------------------

// relationshipType is used to define the type of the relationship between 2 business object models
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
	// backRefs                 []*Relationship  // valued from the business object's init
	mx                         sync.Mutex            // a mutex to avoid some race conditions
	targetModelNames           []utils.ModelName     // the names of BOs pointed by this relationship
	relationType               relationshipType      // valued from the business object's init
	backRefSlice               []*Relationship       // valued from the business object's init
	backRefMap                 map[propKey]bool      // valued from the business object's init
	columnNameForTarget        string                // the name of the target model, needed for persisted polymorphic relationships
	linkTableName              string                // the name of the link table, if this relationship is persisted in a link table
	linkTableSourceColumn      string                // the name of the column in the link table that holds the ID of the source entity
	linkTableSourceModelColumn string                // the name of the column in the link table that holds the model name of the source entity, if the backref relationship is polymorphic
	linkTableTargetColumn      string                // the name of the column in the link table that holds the ID of the target entity
	linkTableTargetModelColumn string                // the name of the column in the link table that holds the model name of the target entity, if this relationship is polymorphic
	targetBoType               *reflection.GoaldType // the reflection type of the target business object, if known
}

// Allows to declare a new monomorphic relationship on a given business object model
func AddRelationship(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool, targetName utils.ModelName) *Relationship {
	relationship := &Relationship{
		businessObjectProperty: businessObjectProperty{
			owner:       owner,
			declaringBO: declaringBO,
			name:        name,
			multiple:    multiple,
			propType:    propertyTypeRELATIONSHIPxMONOM,
		},
		targetModelNames: []utils.ModelName{targetName},
		backRefMap:       make(map[propKey]bool),
	}

	owner.setRelationship(name, relationship)

	return relationship
}

// Allows to declare a new polymorphic relationship on a given business object model
func AddPolyRelationship(owner IBusinessObjectModel, declaringBO utils.ModelName, name string, multiple bool) *Relationship {
	relationship := &Relationship{
		businessObjectProperty: businessObjectProperty{
			owner:       owner,
			declaringBO: declaringBO,
			name:        name,
			multiple:    multiple,
			propType:    propertyTypeRELATIONSHIPxPOLYM,
		},
		backRefMap: make(map[propKey]bool),
	}

	owner.setRelationship(name, relationship)

	return relationship
}

func (r *Relationship) addToBackRefs(backRef *Relationship) {
	r.mx.Lock()
	if !r.backRefMap[backRef.getKey()] {
		r.backRefMap[backRef.getKey()] = true
		r.backRefSlice = append(r.backRefSlice, backRef)
	} else {
		core.PanicMsg("Backref already exists: %s", backRef.getKey())
	}
	r.mx.Unlock()
}

// Sets a relationship as a "child to parent" one; the backref relationship is needed
func (r *Relationship) SetChildToParent(backRefRelation *Relationship) *Relationship {
	r.relationType = relationshipTypeCHILDxTOxPARENT
	r.SetRequiredInDb() // the parent must always exist and be associated to the child

	// this relationship's owner model is then the child of another model
	r.owner.setChildToParentRelationship(r)

	// taking the opportunity here to enrich the backref relationship...
	r.addToBackRefs(backRefRelation)

	// ... like automatically setting on the backref the inverse relation type and this relationship as the backref
	backRefRelation.relationType = relationshipTypePARENTxTOxCHILDREN
	backRefRelation.addToBackRefs(r)

	return r
}

// Sets a relationship as a "parent to children" one; the backref relationship is needed
func (r *Relationship) SetSourceToTarget(backRefRelation *Relationship) *Relationship {
	r.relationType = relationshipTypeSOURCExTOxTARGET

	// taking the opportunity here to enrich the backref relationship...
	r.addToBackRefs(backRefRelation)

	// automatically setting on the backref the inverse relation type and this relationship as the backref
	backRefRelation.relationType = relationshipTypeTARGETxTOxSOURCE
	backRefRelation.addToBackRefs(r)

	return r
}

// Sets a relationship as a "one way" one; is with no back ref
func (r *Relationship) SetOneWay() *Relationship {
	r.relationType = relationshipTypeONExWAY

	return r
}

// returns true if the owner of this relationship is responsible for persisting it in the DB
func (r *Relationship) isDirectlyPersisted() bool {
	return r.relationType == relationshipTypeSOURCExTOxTARGET ||
		r.relationType == relationshipTypeCHILDxTOxPARENT ||
		r.relationType == relationshipTypeONExWAY
}

// returns true if this relationship should it be persisted by the owner using a column on its table
func (r *Relationship) needsColumn() bool {
	return !r.multiple && r.isDirectlyPersisted()
}

// returns true if this relationship should it be persisted by the owner using a link table
func (r *Relationship) needsLinkTable() bool {
	return r.multiple && r.isDirectlyPersisted()
}

// returns the name of the column that goes along with the column containing the ID for a related business object,
// when the relationship is polymorphic, i.e. when we also need to store the BO's model name in the same table
func (r *Relationship) getColumnNameForTargetModel() string {
	if r.columnNameForTarget == "" && r.IsPolymorphic() {
		r.columnNameForTarget = core.PascalToSnake(r.name) + suffixMDL
	}

	return r.columnNameForTarget
}

func (r *Relationship) IsPolymorphic() bool {
	return r.propType == propertyTypeRELATIONSHIPxPOLYM
}

// returns the list of the names of the BOs pointed by this relationship, resolving it if needed
// WARNING: don't use it at init time, as the underlying models might not be fully initialized yet
func (r *Relationship) getTargetModelNames() []utils.ModelName {
	if !r.IsPolymorphic() {
		return r.targetModelNames
	}

	if len(r.targetModelNames) == 0 {
		r.targetModelNames = r.resolveTargetNames()
	}

	return r.targetModelNames
}

// returns the name of the unique target BO pointed by this relationship, or an error message if there is no unique target
// WARNING: don't use it at init time, as the underlying models might not be fully initialized yet
func (r *Relationship) getUniqueTargetName() utils.ModelName {
	targetModelNames := r.getTargetModelNames()
	if len(targetModelNames) != 1 {
		core.PanicMsg("No unique target for relationship '%s' on model '%s'! We found: %+v", r.name, r.owner.GetName(), targetModelNames)
	}

	return targetModelNames[0]
}

// returns the model of the unique target BO pointed by this relationship, or nil if there is no unique target
func (r *Relationship) getUniqueTargetModel() IBusinessObjectModel {
	return modelFor(r.getUniqueTargetName(), true)
}

var interfaceImplementations = map[utils.ModelName][]utils.ModelName{}

// returns the list of the names of the BOs pointed by this relationship, or the list of the types implementing the interface, if it's a polymorphic relationship
// WARNING: don't use it at init time, as the underlying models might not be fully initialized yet
func (r *Relationship) resolveTargetNames() []utils.ModelName {
	// only doing this once, and caching the result for later use
	if len(interfaceImplementations[r.owner.GetName()]) == 0 {
		// now, let's look for all the business object models implementing this interface, and gather them
		for _, model := range core.GetSortedValues(modelRegistry.items) {
			if !model.isInterface() && !model.isAbstract() {
				if boType := reflection.TypeOf(model.NewObject(), false); boType.Implements(r.getTargetBoType()) {
					interfaceImplementations[r.owner.GetName()] = append(interfaceImplementations[r.owner.GetName()], model.GetName())
				}
			}
		}
	}

	return interfaceImplementations[r.owner.GetName()]
}

// returns the reflection type of the target business object of this relationship, if known
// WARNING: don't use it at init time, as the underlying models might not be fully initialized yet
func (r *Relationship) getTargetBoType() reflection.GoaldType {
	if r.targetBoType == nil {
		targetType := r.getStructField().Type() // e.g. Type ITargetObj or []ITargetObj
		if r.multiple {
			targetType = targetType.Elem()
			r.targetBoType = &targetType // e.g. Type ITargetObj
		} else {
			r.targetBoType = &targetType
		}
	}
	return *r.targetBoType
}

// returns the name of the foreign key constraint for this relationship, if it is persisted in the database
func (r *Relationship) getForeignKeyName() string {
	return prefixFK + r.owner.getTableName(false) + "__" + r.getColumnName()
}

// returns the name of the link table for this relationship, if it is persisted in the database
func (r *Relationship) getLinkTableName() string {
	if r.linkTableName == "" {
		if r.IsMultiple() {
			if r.isPolymorphicSource() {
				// if the backref is polymorphic, then we consider the business object declaring this relationship
				r.linkTableName = prefixLINK + core.PascalToSnake(string(r.declaringBO)) + "__" + core.PascalToSnake(r.name)
			} else {
				r.linkTableName = prefixLINK + r.owner.getTableName(false) + "__" + core.PascalToSnake(r.name)
			}
		} else {
			core.PanicMsg("getLinktableName() should not be calle for a non-multiple relationship (%s.%s)", r.owner.GetName(), r.GetName())
		}
	}

	return r.linkTableName
}

// returns the name of the column in the link table that holds the ID of the source entity
func (r *Relationship) getLinkTableSourceColumn() (string, string) {
	if r.linkTableSourceColumn == "" {
		if r.isPolymorphicSource() {
			r.linkTableSourceColumn = prefixSOURCE + core.PascalToSnake(string(r.declaringBO)) + suffixID
			r.linkTableSourceModelColumn = prefixSOURCE + core.PascalToSnake(string(r.declaringBO)) + suffixMDL
		} else {
			r.linkTableSourceColumn = prefixSOURCE + r.owner.getTableName(false) + suffixID
		}
	}

	return r.linkTableSourceColumn, r.linkTableSourceModelColumn
}

// returns the name of the column in the link table that holds the ID of the target entity
func (r *Relationship) getLinkTableTargetColumn() (string, string) {
	if r.linkTableTargetColumn == "" {
		r.linkTableTargetColumn = prefixTARGET + r.getColumnName()
		if r.IsPolymorphic() {
			r.linkTableTargetModelColumn = prefixTARGET + r.getColumnNameForTargetModel()
		}
	}

	return r.linkTableTargetColumn, r.linkTableTargetModelColumn
}

// ------------------------------------------------------------------------------------------------
// Checks on backref relationships
// ------------------------------------------------------------------------------------------------

// applies a getter on all back references of this relationship; ensures all back references return the same result
func (r *Relationship) getOnAllBackRefs[RESULT comparable](def RESULT, cherckFnName string, getFn func(backRef *Relationship) RESULT) RESULT {
	if len(r.backRefSlice) == 0 {
		return def
	}
	result := getFn(r.backRefSlice[0])
	for i, backRef := range r.backRefSlice {
		if i == 0 {
			continue
		}
		if getFn(backRef) != result {
			core.PanicMsg("We have a problem here: %[1]s(%[2]s) is %[3]t while %[1]s(%[4]s) is not",
				cherckFnName, r.backRefSlice[0], result, backRef)
		}
	}
	return result
}

// indicates that this relationship's backref relationships are all polymorphic
func (r *Relationship) isPolymorphicSource() bool {
	return r.getOnAllBackRefs(false, "IsPolymorphic", func(backRef *Relationship) bool {
		return backRef.IsPolymorphic()
	})
}

// indicates that this relationship's backref relationships are all multiple-valued
func (r *Relationship) isMultipleSource() bool {
	return r.getOnAllBackRefs(false, "IsMultiple", func(backRef *Relationship) bool {
		return backRef.IsMultiple()
	})
}

// indicates that this relationship's backref relationships are all required in the database
func (r *Relationship) isSourceRequiredInDB() bool {
	return r.getOnAllBackRefs(false, "IsRequiredInDb", func(backRef *Relationship) bool {
		return backRef.IsRequiredInDb()
	})
}

// indicates that this relationship's backref relationships are all persisted
func (r *Relationship) isIndirectlyPersisted() bool {
	return r.getOnAllBackRefs(false, "isDirectlyPersisted", func(backRef *Relationship) bool {
		return backRef.isDirectlyPersisted()
	})
}

// indicates that this relationship's backref relationships all need a link table
func (r *Relationship) needsReverseLinkTable() bool {
	return r.getOnAllBackRefs(false, "needsLinkTable", func(backRef *Relationship) bool {
		return backRef.needsLinkTable()
	})
}

// retrieves the common name for this relationship's backref relationships
func (r *Relationship) getBackRefName() string {
	return r.getOnAllBackRefs("", "needsLinkTable", func(backRef *Relationship) string {
		return backRef.GetName()
	})
}
