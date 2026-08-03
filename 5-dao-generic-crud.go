// ------------------------------------------------------------------------------------------------
// Here we implement the generic data access instructions involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

// Generic function to insert a business object into the database
func dbInsert(dao IBusinessObjectDAO, bObjs ...IBusinessObject) error {
	if len(bObjs) == 0 {
		return nil
	}

	// controlling we're not trying to insert a business object that already got an ID
	for _, bObj := range bObjs {
		if bObj.GetID() > 0 {
			return Error("Can not insert a business object '%s' that already has an ID", bObj.GetModelName())
		}
	}

	// executing the insert request, which should result in all the BOs getting an ID
	rowToIDMap, errInsert := dao.ExecInsertQuery(bObjs...)
	if errInsert != nil {
		// TODO _JW$2:handle 1062 error (duplicate entry), 1048 (missing column value), 1054 (unknown column)
		return ErrorC(errInsert, "Error while performing insert query for '%s'", dao.getModel().getName())
	}

	// consolidating the DB IDs back into the business objects
	for _, bObj := range bObjs {
		bObj.setID(BObjID(rowToIDMap[bObj.GetPreID()]))
	}

	// so far, we've just handle the entities' properties and persisted single links; let's now handle the multiple links
	if errLinks := dao.ExecInsertLinksQueries(bObjs...); errLinks != nil {
		return ErrorC(errLinks, "Error while performing insert links for '%s'", dao.getModel().getName())
	}

	// yeah, we dit it!
	return nil
}

// Generic function to select business objects from the database
func dbSelect(dao IBusinessObjectDAO, queryName queryName, values IQueryParamsObject) (result []IBusinessObject, err error) {
	// calling the right DAO method
	bObjs, errSelect := dao.ExecSelectQuery(queryName, values)
	if errSelect != nil {
		return nil, ErrorC(errSelect, "Error while performing select query '%s' for '%s'", queryName, dao.getModel().getName())
	}

	return bObjs, nil
}

// func dbLoadList(_ DaoContext, model IBusinessObjectModel) (result []IBusinessObject, err error) {
// 	for _, bObj := range mockDatabase {
// 		if model == bObj.Model() {
// 			result = append(result, bObj)
// 		}
// 	}

// 	return
// }

// func dbLoadList[ResourceType IBusinessObject](_ DaoContext, model IBusinessObjectModel) (result []ResourceType, err error) {
// 	Model := getModel(model)

// 	println(Model)

// 	return
// }

// func dbLoadOne(_ DaoContext, idProp IField, idPropVal string) (result IBusinessObject, err error) {
// 	// for _, bObj := range mockDatabase {
// 	// 	if idProp.ownerModel() == bObj.Model() && idPropVal == bObj.GetValueAsString(idPropVal) {
// 	// 		return bObj, nil
// 	// 	}
// 	// }

// 	return nil, Error("No '%s' found with '%s = %s'", idProp.ownerModel().base().name, idProp.GetName(), idPropVal)
// }

// func dbRemoveOne(_ DaoContext, idProp IField, idPropVal string) (result IBusinessObject, err error) {
// 	// for _, bObj := range mockDatabase {
// 	// 	if idProp.ownerModel() == bObj.Model() && idPropVal == bObj.GetValueAsString(idPropVal) {
// 	// 		delete(mockDatabase, string(bObj.GetID()))
// 	// 		return bObj, nil
// 	// 	}
// 	// }

// 	return nil, Error("No '%s' found with '%s = %s'", idProp.ownerModel().base().name, idProp.GetName(), idPropVal)
// }

// func dbUpdate(_ DaoContext, input IBusinessObject) error {
// 	// instore := mockDatabase[string(input.GetID())]
// 	// if instore == nil {
// 	// 	return Error("No object exists with ID %s", input.GetID())
// 	// }

// 	// input.setModelName(instore.getModelName())

// 	// mockDatabase[string(input.GetID())] = input

// 	return nil
// }
