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
// func (boSpecs *businessObjectClass) GetAllProperties() []iBusinessObjectProperty {
// 	if boSpecs.allProperties == nil {
// 		for _, field := range boSpecs.fields {
// 			boSpecs.allProperties = append(boSpecs.allProperties, field)
// 		}

// 		for _, relationship := range boSpecs.getRelationshipsWithColumn() {
// 			boSpecs.allProperties = append(boSpecs.allProperties, relationship)
// 		}

// 		sort.SliceStable(boSpecs.allProperties, func(i, j int) bool {
// 			return boSpecs.allProperties[i].getName() < boSpecs.allProperties[j].getName()
// 		})
// 	}

// 	return boSpecs.allProperties
// }

// getPersistedProperties returns the sorted list of the properties persisted
// within the BO class' table, i.e. the persisted single Relationships + the persisted fields
func (boSpecs *businessObjectSpecs) getPersistedProperties() []iBusinessObjectProperty {
	if boSpecs.persistedProperties == nil {
		// how many persisted properties - fields + single Relationships - do we have ?
		// nbFields := len(boSpecs.fields)
		// size := nbFields + len(boSpecs.getRelationshipsWithColumn())

		// let's gather all the persisted properties
		// boSpecs.persistedProperties = make([]iBusinessObjectProperty, size)
		for _, field := range boSpecs.fields {
			if !field.isNotPersisted() {
				boSpecs.persistedProperties = append(boSpecs.persistedProperties, field)
			}
		}

		for _, relationship := range boSpecs.getRelationshipsWithColumn() {
			boSpecs.persistedProperties = append(boSpecs.persistedProperties, relationship)
		}

		// now, let's sort them to have a nicely sorted list of columns for each table
		// we make sure the ID column is always at 1st position
		sort.SliceStable(boSpecs.persistedProperties, func(i, j int) bool {
			property1Name := boSpecs.persistedProperties[i].getColumnName()
			property2Name := boSpecs.persistedProperties[j].getColumnName()
			if property1Name == "id" {
				return true
			}
			if property2Name == "id" {
				return false
			}

			return property1Name < property2Name
		})
	}

	return boSpecs.persistedProperties
}

// getRelationshipsWithColumn returns the sorted list of the fields that are persisted
func (boSpecs *businessObjectSpecs) getRelationshipsWithColumn() []*Relationship {
	// initialising it, the first time we need it
	if boSpecs.relationshipsWithColumn == nil {
		// first, we retrieve a list of IDs of the Relationships that are persisted
		relationshipsWithColumnNames := []string{}

		for relationshipName, relationship := range boSpecs.relationships {
			if relationship.needsColumn() {
				relationshipsWithColumnNames = append(relationshipsWithColumnNames, string(relationshipName))
			}
		}

		// sorting that list
		sort.Strings(relationshipsWithColumnNames)

		// creating the list of persisted relationships
		boSpecs.relationshipsWithColumn = make([]*Relationship, len(relationshipsWithColumnNames))

		// using that list to build a sorted list of persisted relationships
		for i := 0; i < len(relationshipsWithColumnNames); i++ {
			boSpecs.relationshipsWithColumn[i] = boSpecs.relationships[relationshipsWithColumnNames[i]]
		}
	}

	return boSpecs.relationshipsWithColumn
}
