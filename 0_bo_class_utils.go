// ------------------------------------------------------------------------------------------------
// Some utilities to help build classes
// ------------------------------------------------------------------------------------------------
package goald

import (
	"reflect"
	"sort"
	"time"
)

var (
	typeBUSINESSxOBJECT   = reflect.TypeOf((*BusinessObject)(nil)).Elem()
	typeIxBUSINESSxOBJECT = reflect.TypeOf((*IBusinessObject)(nil)).Elem()
	typeTIME              = reflect.TypeOf((*time.Time)(nil))
	typeIxENUM            = reflect.TypeOf((*IEnum)(nil)).Elem()

// TypeAMOUNT          = reflect.TypeOf((core.Amount)(0))
// TypeDURATION        = reflect.TypeOf((time.Duration)(0))
// TypeEMAIL           = reflect.TypeOf((core.Email)(""))
// TypeIACTIONTYPE     = reflect.TypeOf((*IActionType)(nil)).Elem()
// TypeIENUM           = reflect.TypeOf((*IEnum)(nil)).Elem()
// TypeINT             = reflect.TypeOf((int)(0))
// TypeINT64           = reflect.TypeOf((int64)(0))
// TypeIUSER           = reflect.TypeOf((*IUser)(nil)).Elem()
// TypeRESOURCE        = reflect.TypeOf((*Resource)(nil)).Elem()
// TypeURL             = reflect.TypeOf((core.URL)(""))
// TypeWEBOPERATION    = reflect.TypeOf((*WebOperation)(nil)).Elem()
// TypeJSONString      = reflect.TypeOf((core.JSONString)(""))
)

// getPropertyType returns the property type of a given structfield
func getPropertyType(structField reflect.StructField) (propertyType PropertyType, multiple bool) {

	// to debug - to comment/uncomment when needed
	// if structField.Name == "Num" {
	// fmt.Printf("\n--------------------------")
	// fmt.Printf("\nName: %s ", structField.Name)
	// fmt.Printf("\nType: %s ", structField.Type)
	// fmt.Printf("\nKind: %s ", structField.Type.Kind())
	// }

	// a business object's real property must be exported, and therefore PkgPath should be empty
	// Cf. https://golang.org/pkg/reflect/#StructField
	if fieldType := structField.Type; structField.PkgPath == "" {
		// // detecting the Loadedrelationships: special case of []string, which is not allowed anywhere else
		// if structField.Name == __BUSINESS OBJEcT__FieldLOADEDrelationshipS {
		// 	return PropertyTypeSTRING, true
		// }

		// // detecting an __BUSINESS OBJEcT__ ID
		// if fieldType == Type__BUSINESS OBJEcT__ID {
		// 	return PropertyType__BUSINESS OBJEcT__ID, false
		// }

		// detecting an enum
		if fieldType.Implements(typeIxENUM) {
			return PropertyTypeENUM, false
		}

		// detecting a time
		if fieldType == typeTIME {
			return PropertyTypeDATE, false
		}

		// // detecting a URL
		// if fieldType == TypeURL {
		// 	return PropertyTypeURL, false
		// }

		// // detecting an email
		// if fieldType == TypeEMAIL {
		// 	return PropertyTypeEMAIL, false
		// }

		// // detecting an __BUSINESS OBJEcT__ reference
		// if fieldType == Type__BUSINESS OBJEcT__REFERENCE {
		// 	return PropertyType__BUSINESS OBJEcT__REFERENCE, false
		// }

		// // detecting an amount
		// if fieldType == TypeAMOUNT {
		// 	return PropertyTypeAMOUNT, false
		// }

		// // detecting a json object
		// if fieldType == TypeJSONString {
		// 	return PropertyTypeJSON, false
		// }

		// detecting the basic types here
		switch fieldKind := fieldType.Kind(); fieldKind {
		case reflect.Bool:
			return PropertyTypeBOOL, false

		case reflect.String:
			return PropertyTypeSTRING, false

			// We only allow 3 types of int, for simplicity's sake
		case reflect.Int:
			return PropertyTypeINT, false

			// We only allow 3 types of int, for simplicity's sake
		case reflect.Int64, reflect.Uint64:
			return PropertyTypeINT64, false

		case reflect.Float32:
			return PropertyTypeREAL32, false

		case reflect.Float64:
			return PropertyTypeREAL64, false
		}

		// // detecting an enum list
		// if innerSliceType := fieldType.Elem(); fieldKind == reflect.Slice && innerSliceType.Implements(TypeIENUM) {
		// 	return PropertyTypeENUM, true
		// }

		// detecting a multiple relationship to business objects
		if innerSliceType := fieldType.Elem(); fieldType.Kind() == reflect.Slice && innerSliceType.Implements(typeIxBUSINESSxOBJECT) {
			return PropertyTypeRELATIONSHIP, true
		}

		// detecting a single relationship to a business object
		if fieldType.Kind() == reflect.Ptr && fieldType.Implements(typeIxBUSINESSxOBJECT) {
			return PropertyTypeRELATIONSHIP, false
		}
	}

	// this happens with technical fields !
	return PropertyTypeUNKNOWN, false
}

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
		for _, field := range boClass.fields { // TODO only take the persisted fields
			boClass.persistedProperties = append(boClass.persistedProperties, field)
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
