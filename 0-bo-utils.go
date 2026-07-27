// ------------------------------------------------------------------------------------------------
// Some utilities to help build classes
// ------------------------------------------------------------------------------------------------
package goald

import (
	"sort"

	"github.com/aldesgroup/goald/features/reflection"
)

var (
	typeBUSINESSxOBJECT   = reflection.TypeOf((*BusinessObject)(nil), true)
	typeURLxQUERYxOBJECT  = reflection.TypeOf((*URLQueryParams)(nil), true)
	typeIxBUSINESSxOBJECT = reflection.TypeOf((*IBusinessObject)(nil), true)
	typeIxENUM            = reflection.TypeOf((*IEnum)(nil), true)
)

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
		// first, we retrieve a list of IDs of the Relationships that are directly persisted
		// i.e. through a column in this BO model's table, and not through a link table nor a foreign table
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

// getRelationshipsWithLinkTable returns the sorted list of the relationships that are persisted using a link table in the BO model's table
func (model *businessObjectModel) getRelationshipsWithLinkTable() []*Relationship {
	// initialising it, the first time we need it
	if model.relationshipsWithLinkTable == nil {
		// first, we retrieve a list of IDs of the Relationships that are directly persisted
		// i.e. through a link table in this BO model's table, and not through a column nor a foreign table
		relationshipsWithLinkTableNames := []string{}

		for relationshipName, relationship := range model.relationships {
			if relationship.needsLinkTable() {
				relationshipsWithLinkTableNames = append(relationshipsWithLinkTableNames, string(relationshipName))
			}
		}

		// sorting that list
		sort.Strings(relationshipsWithLinkTableNames)

		// creating the list of persisted relationships
		model.relationshipsWithLinkTable = make([]*Relationship, len(relationshipsWithLinkTableNames))

		// using that list to build a sorted list of persisted relationships
		for i := 0; i < len(relationshipsWithLinkTableNames); i++ {
			model.relationshipsWithLinkTable[i] = model.relationships[relationshipsWithLinkTableNames[i]]
		}
	}

	return model.relationshipsWithLinkTable
}
