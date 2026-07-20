// ------------------------------------------------------------------------------------------------
// Here we implement the generic data access instructions involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

// Generic function to insert a business object into the database
func dbInsert(dao IBusinessObjectDAO, bObj IBusinessObject) error {
	// controlling we're not trying to insert a business object that already got an ID
	if bObj.GetID() > 0 {
		return Error("Can not insert business object '%s' that already has an ID", bObj.GetClassName(bObj))
	}

	// // building the array of values to fill the query
	// values := getSQLValues(__REPLACE__)

	// // preparing the query execution context, and setting the debug values in verbose mode
	// queryExecContext := insertQuery.With(values...)
	// if dbContext.IsVerbose() {
	// 	queryExecContext.DebugWith(getMaskedSQLValues(__REPLACE__)...)
	// }

	// executing the insert request, and retrieving the new ID; this works whether the DB is
	// Postgres (using a 'RETURNING id' clause) or MySQL (using sql.Result.LastInsertId())
	id, errInsert := dao.ExecInsertQuery(bObj)
	if errInsert != nil {
		// TODO _JW$2:handle 1062 error (duplicate entry), 1048 (missing column value), 1054 (unknown column)
		return ErrorC(errInsert, "Error while performing insert query")
	}

	// we can set the ID
	bObj.setID(BObjID(id))

	// let's log the new ID, since it does not appear in the logging of the request
	// dbContext.Log().Debug("Created a new '%s' with ID: %d", __REPLACE__.GetKind(), __REPLACE__.GetID())

	// // so far, we've just handle the entities' properties and persisted single links; let's now handle the multiple links
	// if errLinks := dbInsertTableLinks(dbContext, __REPLACE__); errLinks != nil {
	// 	return errLinks
	// }

	// yeah, we dit it!
	return nil
}

// func dbLoadList(_ DaoContext, model IBusinessObjectModel) (result []IBusinessObject, err error) {
// 	for _, bObj := range mockDatabase {
// 		if model == bObj.Class() {
// 			result = append(result, bObj)
// 		}
// 	}

// 	return
// }

// func dbLoadList[ResourceType IBusinessObject](_ DaoContext, model IBusinessObjectModel) (result []ResourceType, err error) {
// 	class := getClass(model)

// 	println(class)

// 	return
// }

// func dbLoadOne(_ DaoContext, idProp IField, idPropVal string) (result IBusinessObject, err error) {
// 	// for _, bObj := range mockDatabase {
// 	// 	if idProp.ownerModel() == bObj.Class() && idPropVal == bObj.GetValueAsString(idPropVal) {
// 	// 		return bObj, nil
// 	// 	}
// 	// }

// 	return nil, Error("No '%s' found with '%s = %s'", idProp.ownerModel().base().name, idProp.GetName(), idPropVal)
// }

// func dbRemoveOne(_ DaoContext, idProp IField, idPropVal string) (result IBusinessObject, err error) {
// 	// for _, bObj := range mockDatabase {
// 	// 	if idProp.ownerModel() == bObj.Class() && idPropVal == bObj.GetValueAsString(idPropVal) {
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

// 	// input.setClassName(instore.getClassName())

// 	// mockDatabase[string(input.GetID())] = input

// 	return nil
// }
