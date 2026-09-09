// ------------------------------------------------------------------------------------------------
// Some utilities to help use business object models
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"sort"

	"github.com/aldesgroup/goald/features/reflection"
	"github.com/aldesgroup/goald/features/utils"
)

var (
	typeBUSINESSxOBJECT     = reflection.TypeOf((*BusinessObject)(nil), true)
	typeSEARCHxPARAMxVALUES = reflection.TypeOf((*SearchParamValues)(nil), true)
	typeIxBUSINESSxOBJECT   = reflection.TypeOf((*IBusinessObject)(nil), true)
	typeIxENUM              = reflection.TypeOf((*IEnum)(nil), true)
)

// ------------------------------------------------------------------------------------------------
// Retrieving models / instantiating business objects
// ------------------------------------------------------------------------------------------------

func modelFor(modelName utils.ModelName, failIfNil ...bool) IBusinessObjectModel {
	if modelRegistry.items[modelName] == nil && (len(failIfNil) > 0 && failIfNil[0]) {
		panic(fmt.Sprintf("It looks like no model named '%s' has been registered, "+
			"i.e. its package has probably not been 'included', i.e. imported in the start.go file,"+
			" like this: import _ \"module_full_name/_include/package_name\"", modelName))
	}
	return modelRegistry.items[modelName]
}

func NewBusinessObject(modelName utils.ModelName, id BObjID) IBusinessObject {
	bObj := modelFor(modelName, true).NewObject().(IBusinessObject)
	bObj.setID(BObjID(id))
	return bObj
}

// ------------------------------------------------------------------------------------------------
// Some "views" on a model's fields and properties
// ------------------------------------------------------------------------------------------------

// getPersistedProperties returns the sorted list of the properties persisted
// within the BO model's table, i.e. the persisted single Relationships + the persisted fields
func (model *businessObjectModel) getPersistedProperties() []IBusinessObjectProperty {
	// initialising it, the first time we need it
	if model.persistedProperties == nil {
		// let's gather all the persisted properties - the fields first
		for _, field := range model.fields {
			// special case of the rowID field, which is only needed when the associated DB handles the RETURNING clause
			if field.GetName() == boFieldPreID {
				if model.db.is.SupportsReturningID() {
					model.persistedProperties = append(model.persistedProperties, field)
				}
			} else {
				if !field.isNotPersisted() {
					model.persistedProperties = append(model.persistedProperties, field)
				}
			}
		}

		// and the relationships
		for _, relationship := range model.getRelationshipsWithColumn() {
			model.persistedProperties = append(model.persistedProperties, relationship)
		}

		// now, let's sort them to have a nicely sorted list of columns for each table
		// we make sure the ID column is always at 1st position
		sort.SliceStable(model.persistedProperties, func(i, j int) bool {
			property1Name := model.persistedProperties[i].GetName()
			property2Name := model.persistedProperties[j].GetName()
			if property1Name == BoFieldID {
				return true
			}
			if property2Name == BoFieldID {
				return false
			}
			if property1Name == boFieldPreID {
				return true
			}
			if property2Name == boFieldPreID {
				return false
			}

			return property1Name < property2Name
		})
	}

	return model.persistedProperties
}

// getRelationshipsWithColumn returns the sorted list of the relationships that are persisted using a column in the BO model's table
func (model *businessObjectModel) getRelationshipsWithColumn() []*Relationship {
	// initialising it, the first time we need it
	if model.relationshipsWithColumn == nil {
		// we retrieve a list of the Relationships that are directly persisted
		for _, relationship := range model.relationships {
			if relationship.needsColumn() {
				model.relationshipsWithColumn = append(model.relationshipsWithColumn, relationship)
			}
		}

		// then, sorting that list
		sort.SliceStable(model.relationshipsWithColumn, func(i, j int) bool {
			return model.relationshipsWithColumn[i].GetName() < model.relationshipsWithColumn[j].GetName()
		})
	}

	return model.relationshipsWithColumn
}

// getRelationshipsWithLinkTable returns the sorted list of the relationships that are persisted using a link table in the BO model's table
func (model *businessObjectModel) getRelationshipsWithLinkTable() []*Relationship {
	// initialising it, the first time we need it
	if model.relationshipsWithLinkTable == nil {
		// first, we retrieve a list of the Relationships that are persisted through a link table
		for _, relationship := range model.relationships {
			if relationship.needsLinkTable() {
				model.relationshipsWithLinkTable = append(model.relationshipsWithLinkTable, relationship)
			}
		}

		// then, sorting that list
		sort.SliceStable(model.relationshipsWithLinkTable, func(i, j int) bool {
			return model.relationshipsWithLinkTable[i].GetName() < model.relationshipsWithLinkTable[j].GetName()
		})
	}

	return model.relationshipsWithLinkTable
}

func (model *businessObjectModel) getRelationshipsWithRequiredBackref() []*Relationship {
	// initialising it, the first time we need it
	if model.relationshipsWithRequiredBackref == nil {
		// first, we retrieve a list of the relationships that have a required back reference
		for _, relationship := range model.relationships {
			if relationship.backRef != nil && relationship.backRef.IsRequiredInDb() {
				model.relationshipsWithRequiredBackref = append(model.relationshipsWithRequiredBackref, relationship)
			}
		}

		// then, sorting that list
		sort.SliceStable(model.relationshipsWithRequiredBackref, func(i, j int) bool {
			return model.relationshipsWithRequiredBackref[i].GetName() < model.relationshipsWithRequiredBackref[j].GetName()
		})

	}

	return model.relationshipsWithRequiredBackref
}
