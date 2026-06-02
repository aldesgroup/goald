// ------------------------------------------------------------------------------------------------
// Some utilities to help build classes
// ------------------------------------------------------------------------------------------------
package goald

import (
	"sort"

	"github.com/aldesgroup/goald/features/utils"
)

var (
	typeBUSINESSxOBJECT   = utils.TypeOf((*BusinessObject)(nil), true)
	typeURLxQUERYxOBJECT  = utils.TypeOf((*URLQueryParams)(nil), true)
	typeIxBUSINESSxOBJECT = utils.TypeOf((*IBusinessObject)(nil), true)
	typeIxENUM            = utils.TypeOf((*IEnum)(nil), true)
)

// // GetAllProperties returns all this class' properties
// func (model *businessObjectClass) GetAllProperties() []iBusinessObjectProperty {
// 	if model.allProperties == nil {
// 		for _, field := range model.fields {
// 			model.allProperties = append(model.allProperties, field)
// 		}

// 		for _, relationship := range model.getRelationshipsWithColumn() {
// 			model.allProperties = append(model.allProperties, relationship)
// 		}

// 		sort.SliceStable(model.allProperties, func(i, j int) bool {
// 			return model.allProperties[i].getName() < model.allProperties[j].getName()
// 		})
// 	}

// 	return model.allProperties
// }

// getPersistedProperties returns the sorted list of the properties persisted
// within the BO class' table, i.e. the persisted single Relationships + the persisted fields
func (model *businessObjectModel) getPersistedProperties() []iBusinessObjectProperty {
	if model.persistedProperties == nil {
		// how many persisted properties - fields + single Relationships - do we have ?
		// nbFields := len(model.fields)
		// size := nbFields + len(model.getRelationshipsWithColumn())

		// let's gather all the persisted properties
		// model.persistedProperties = make([]iBusinessObjectProperty, size)
		for _, field := range model.fields {
			if !field.isNotPersisted() {
				model.persistedProperties = append(model.persistedProperties, field)
			}
		}

		for _, relationship := range model.getRelationshipsWithColumn() {
			model.persistedProperties = append(model.persistedProperties, relationship)
		}

		// now, let's sort them to have a nicely sorted list of columns for each table
		// we make sure the ID column is always at 1st position
		sort.SliceStable(model.persistedProperties, func(i, j int) bool {
			property1Name := model.persistedProperties[i].getColumnName()
			property2Name := model.persistedProperties[j].getColumnName()
			if property1Name == "id" {
				return true
			}
			if property2Name == "id" {
				return false
			}

			return property1Name < property2Name
		})
	}

	return model.persistedProperties
}

// getRelationshipsWithColumn returns the sorted list of the fields that are persisted
func (model *businessObjectModel) getRelationshipsWithColumn() []*Relationship {
	// initialising it, the first time we need it
	if model.relationshipsWithColumn == nil {
		// first, we retrieve a list of IDs of the Relationships that are persisted
		relationshipsWithColumnNames := []string{}

		for relationshipName, relationship := range model.relationships {
			if relationship.needsColumn() {
				relationshipsWithColumnNames = append(relationshipsWithColumnNames, string(relationshipName))
			}
		}

		// sorting that list
		sort.Strings(relationshipsWithColumnNames)

		// creating the list of persisted relationships
		model.relationshipsWithColumn = make([]*Relationship, len(relationshipsWithColumnNames))

		// using that list to build a sorted list of persisted relationships
		for i := 0; i < len(relationshipsWithColumnNames); i++ {
			model.relationshipsWithColumn[i] = model.relationships[relationshipsWithColumnNames[i]]
		}
	}

	return model.relationshipsWithColumn
}
