package goald

// cleaned is a utility function that removes cycles from a business object and returns it.
func cleaned[BOTYPE IBusinessObject](bo BOTYPE) BOTYPE {
	bo.RemoveCycles()
	return bo
}
