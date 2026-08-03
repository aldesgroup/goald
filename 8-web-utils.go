package goald

// cleaned is a utility function that removes cycles from a business object and returns it.
func cleaned[BOTYPE IBusinessObject](bObj BOTYPE) BOTYPE {
	bObj.RemoveCycles()
	return bObj
}

// allCleaned removes cycles from a slice of business objects and casts each of them to the specified type BOTYPE.
func allCleaned[BOTYPE IBusinessObject](bObjs []IBusinessObject) []BOTYPE {
	cleanedBObjs := make([]BOTYPE, len(bObjs))
	for i, bObj := range bObjs {
		bObj.RemoveCycles()
		cleanedBObjs[i] = bObj.(BOTYPE)
	}
	return cleanedBObjs
}
