// ------------------------------------------------------------------------------------------------
// Here are the enums used for building business object classes
// ------------------------------------------------------------------------------------------------
package goald

// ------------------------------------------------------------------------------------------------
// the type a business object property can have
// ------------------------------------------------------------------------------------------------

// PropertyType represents the type of a business object's property
type PropertyType int

const (
	// PropertyTypeUNKNOWN : when the property type is not recognized
	PropertyTypeUNKNOWN PropertyType = iota - 1

	// PropertyTypeBOOL : for boolean properties
	PropertyTypeBOOL

	// PropertyTypeSTRING : for string properties
	PropertyTypeSTRING

	// PropertyTypeINT : for integer properties
	PropertyTypeINT

	// PropertyTypeINT64 : for integer properties
	PropertyTypeINT64

	// PropertyTypeREAL32 : for 32-bits real number properties
	PropertyTypeREAL32

	// PropertyTypeREAL64 : for 64-bits real number properties
	PropertyTypeREAL64

	// PropertyTypeDATE : for date properties
	PropertyTypeDATE

	// PropertyTypeENUM : for enum properties
	PropertyTypeENUM

	// PropertyTypeRELATIONSHIP : for relationships to other entities
	PropertyTypeRELATIONSHIP

	// // PropertyTypeENTITYREFERENCE : for entity references
	// PropertyTypeENTITYREFERENCE

	// // PropertyTypeURL : for URL string properties
	// PropertyTypeURL

	// // PropertyTypeEMAIL : for email properties
	// PropertyTypeEMAIL

	// // PropertyTypeENTITYID : for entity IDs
	// PropertyTypeENTITYID

	// // PropertyTypeAMOUNT : for amount of money (which is a float64)
	// PropertyTypeAMOUNT

	// // PropertyTypeJSON : for JSON strings
	// PropertyTypeJSON
)

var propertyTypes = map[int]string{
	int(PropertyTypeBOOL):         "boolean",
	int(PropertyTypeDATE):         "date",
	int(PropertyTypeUNKNOWN):      "unknown",
	int(PropertyTypeSTRING):       "string",
	int(PropertyTypeINT):          "integer",
	int(PropertyTypeREAL32):       "real number",
	int(PropertyTypeREAL64):       "real number (double precision)",
	int(PropertyTypeENUM):         "enum",
	int(PropertyTypeRELATIONSHIP): "relationship",
	// int(PropertyTypeENTITYREFERENCE): "entity reference",
	// int(PropertyTypeURL):             "URL",
	// int(PropertyTypeAMOUNT):          "amount",
	// int(PropertyTypeEMAIL):           "email",
	// int(PropertyTypeJSON):            "json string",
}

func (thisProperty PropertyType) String() string {
	return propertyTypes[int(thisProperty)]
}

// Val helps implement the IEnum interface
func (thisProperty PropertyType) Val() int {
	return int(thisProperty)
}

// Values helps implement the IEnum interface
func (thisProperty PropertyType) Values() map[int]string {
	return propertyTypes
}

// // codeName returns the property's name, as written in the code
// func (thisProperty PropertyType) codeName() string {
// 	switch thisProperty {
// 	case PropertyTypeBOOL:
// 		return "PropertyTypeBOOL"
// 	case PropertyTypeDATE:
// 		return "PropertyTypeDATE"
// 	case PropertyTypeUNKNOWN:
// 		return "PropertyTypeUNKNOWN"
// 	case PropertyTypeSTRING:
// 		return "PropertyTypeSTRING"
// 	case PropertyTypeINT:
// 		return "PropertyTypeINT"
// 	case PropertyTypeREAL32:
// 		return "PropertyTypeREAL32"
// 	case PropertyTypeREAL64:
// 		return "PropertyTypeREAL64"
// 	case PropertyTypeENUM:
// 		return "PropertyTypeENUM"
// 	case PropertyTypeRELATIONSHIP:
// 		return "PropertyTypeRELATIONSHIP"
// 	default:
// 		utils.Panicf("Unhandle property type %d", thisProperty)
// 		return ""
// 	}
// }
