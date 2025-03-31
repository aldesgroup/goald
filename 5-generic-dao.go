// ------------------------------------------------------------------------------------------------
// Here we implement the generic data access instructions involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

// var mockDatabase map[string]IBusinessObject = map[string]IBusinessObject{}

func dbInsert(_ DaoContext, bObj IBusinessObject) error {
	// uuid, errUuid := uuid.NewRandom()
	// if errUuid != nil {
	// 	return ErrorC(errUuid, "could not generate a new UUID")
	// }

	// bObj.setID(uuid.String())

	// mockDatabase[(bObj.GetID())] = bObj

	return nil
}

// func dbLoadList(_ DaoContext, boSpecs IBusinessObjectSpecs) (result []IBusinessObject, err error) {
// 	for _, bObj := range mockDatabase {
// 		if boSpecs == bObj.Class() {
// 			result = append(result, bObj)
// 		}
// 	}

// 	return
// }

func dbLoadList[ResourceType IBusinessObject](_ DaoContext, boSpecs IBusinessObjectSpecs) (result []ResourceType, err error) {
	class := getClass(boSpecs)

	println(class)

	return
}

func dbLoadOne(_ DaoContext, idProp IField, idPropVal string) (result IBusinessObject, err error) {
	// for _, bObj := range mockDatabase {
	// 	if idProp.ownerSpecs() == bObj.Class() && idPropVal == bObj.GetValueAsString(idPropVal) {
	// 		return bObj, nil
	// 	}
	// }

	return nil, Error("No '%s' found with '%s = %s'", idProp.ownerSpecs().base().name, idProp.getName(), idPropVal)
}

func dbRemoveOne(_ DaoContext, idProp IField, idPropVal string) (result IBusinessObject, err error) {
	// for _, bObj := range mockDatabase {
	// 	if idProp.ownerSpecs() == bObj.Class() && idPropVal == bObj.GetValueAsString(idPropVal) {
	// 		delete(mockDatabase, string(bObj.GetID()))
	// 		return bObj, nil
	// 	}
	// }

	return nil, Error("No '%s' found with '%s = %s'", idProp.ownerSpecs().base().name, idProp.getName(), idPropVal)
}

func dbUpdate(_ DaoContext, input IBusinessObject) error {
	// instore := mockDatabase[string(input.GetID())]
	// if instore == nil {
	// 	return Error("No object exists with ID %s", input.GetID())
	// }

	// input.setClassName(instore.getClassName())

	// mockDatabase[string(input.GetID())] = input

	return nil
}
