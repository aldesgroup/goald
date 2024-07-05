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
// func (boClass *businessObjectClass) GetAllProperties() []iBusinessObjectProperty {
// 	if boClass.allProperties == nil {
// 		for _, field := range boClass.fields {
// 			boClass.allProperties = append(boClass.allProperties, field)
// 		}

// 		for _, relationship := range boClass.getRelationshipsWithColumn() {
// 			boClass.allProperties = append(boClass.allProperties, relationship)
// 		}

// 		sort.SliceStable(boClass.allProperties, func(i, j int) bool {
// 			return boClass.allProperties[i].getName() < boClass.allProperties[j].getName()
// 		})
// 	}

// 	return boClass.allProperties
// }

// getPersistedProperties returns the sorted list of the properties persisted
// within the BO class' table, i.e. the persisted single Relationships + the persisted fields
func (boClass *businessObjectClass) getPersistedProperties() []iBusinessObjectProperty {
	if boClass.persistedProperties == nil {
		// how many persisted properties - fields + single Relationships - do we have ?
		// nbFields := len(boClass.fields)
		// size := nbFields + len(boClass.getRelationshipsWithColumn())

		// let's gather all the persisted properties
		// boClass.persistedProperties = make([]iBusinessObjectProperty, size)
		for _, field := range boClass.fields {
			if !field.isNotPersisted() {
				boClass.persistedProperties = append(boClass.persistedProperties, field)
			}
		}

		for _, relationship := range boClass.getRelationshipsWithColumn() {
			boClass.persistedProperties = append(boClass.persistedProperties, relationship)
		}

		// now, let's sort them to have a nicely sorted list of columns for each table
		// we make sure the ID column is always at 1st position
		sort.SliceStable(boClass.persistedProperties, func(i, j int) bool {
			property1Name := boClass.persistedProperties[i].getColumnName()
			property2Name := boClass.persistedProperties[j].getColumnName()
			if property1Name == "id" {
				return true
			}
			if property2Name == "id" {
				return false
			}

			return property1Name < property2Name
		})
	}

	return boClass.persistedProperties
}

// getRelationshipsWithColumn returns the sorted list of the fields that are persisted
func (boClass *businessObjectClass) getRelationshipsWithColumn() []*Relationship {
	// initialising it, the first time we need it
	if boClass.relationshipsWithColumn == nil {
		// first, we retrieve a list of IDs of the Relationships that are persisted
		relationshipsWithColumnNames := []string{}

		for relationshipName, relationship := range boClass.relationships {
			if relationship.needsColumn() {
				relationshipsWithColumnNames = append(relationshipsWithColumnNames, string(relationshipName))
			}
		}

		// sorting that list
		sort.Strings(relationshipsWithColumnNames)

		// creating the list of persisted relationships
		boClass.relationshipsWithColumn = make([]*Relationship, len(relationshipsWithColumnNames))

		// using that list to build a sorted list of persisted relationships
		for i := 0; i < len(relationshipsWithColumnNames); i++ {
			boClass.relationshipsWithColumn[i] = boClass.relationships[relationshipsWithColumnNames[i]]
		}
	}

	return boClass.relationshipsWithColumn
}
