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
