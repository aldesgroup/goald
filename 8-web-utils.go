package goald

// cleaned is a utility function that removes cycles from a business object and returns it.
func cleaned[BOTYPE IBusinessObject](bObj BOTYPE) BOTYPE {
	bObj.RemoveCycles()
	return bObj
}

// allCleaned is a utility function that removes cycles from a slice of business objects and returns it.
func allCleaned[BOTYPE IBusinessObject](bObjs []BOTYPE) []BOTYPE {
	for _, bObj := range bObjs {
		bObj.RemoveCycles()
	}
	return bObjs
}

// allCastAndCleaned removes cycles from a slice of business objects and casts each of them to the specified type BOTYPE.
func allCastAndCleaned[BOTYPE IBusinessObject](bObjs []IBusinessObject) []BOTYPE {
	castBObjs := make([]BOTYPE, len(bObjs))
	for i, bObj := range bObjs {
		bObj.RemoveCycles()
		castBObjs[i] = bObj.(BOTYPE)
	}
	return castBObjs
}
