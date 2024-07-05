// ------------------------------------------------------------------------------------------------
// Here we implement the generic business logic involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

// Controls, DB-inserts, post-treats the given BO
// TODO - quick & dirty implem for now
// func CreateBO[BOTYPE IBusinessObject](bloCtx BloContext, bObj *BOTYPE) error {
func CreateBO(bloCtx BloContext, bObj IBusinessObject) error {
	if bObj == nil {
		return nil
	}

	// the business object should not have an ID already
	if bObj.GetID() != 0 {
		return Error("Could not create object since it already has an ID (%s)", bObj.GetID())
	}

	// do we have stuff to perform on the __BOBJ__ before inserting it ?
	if err := bObj.ChangeBeforeInsert(bloCtx); err != nil {
		return ErrorC(err, "Could not create object since the pre-insert got an error")
	}

	// check of "functional / business" validity
	if err := bObj.IsValid(bloCtx); err != nil {
		return ErrorC(err, "Could not create object since it is not valid")
	}

	// setting some tracking info
	// bObj.SetCreatedByID(biContext.GetCurrentUser().GetID())
	// bObj.SetCreatedBy(biContext.GetCurrentUser().GetLabel())
	// bObj.SetCreation(core.Now())
	// bObj.Set__BOBJ__Status(__BOBJ__StatusCREATED)

	// pushing to the DB ! We're going to add a new line within the __BOBJ__'s table
	if err := dbInsert(bloCtx.GetDaoContext(), bObj); err != nil {
		return ErrorC(err, "Could not create object because of a problem with the DB")
	}

	// we have stuff to do after the insertion ? yeah ? really ? let's do it now !
	if err := bObj.ChangeAfterInsert(bloCtx); err != nil {
		return ErrorC(err, "Could not post-insert object since it got an error")
	}

	// // 'guess everything is alrite here
	// return nilreturn nil
	return nil
}

func LoadBOs[ResourceType IBusinessObject](bloCtx BloContext, boClass IBusinessObjectClass, loadingType LoadingType) ([]ResourceType, error) {
	// func LoadBOs(bloCtx BloContext, boClass IBusinessObjectClass, loadingType LoadingType) ([]ResourceType, error) {
	loadedBOs, errLoad := dbLoadList[ResourceType](bloCtx.GetDaoContext(), boClass)
	// loadedBOs, errLoad := dbLoadList(bloCtx.GetDaoContext(), boClass)

	if errLoad != nil {
		return nil, ErrorC(errLoad, "error while loading a list of '%s'", boClass.base().name)
	}

	// TODO add post read, i.e.:
	// - reading the links, using the LoadingType
	// - on each BO: setting the loadingID + check if reading is ok, then do after read changes

	return loadedBOs, nil
}

func ReadBO(bloCtx BloContext, idProp IField, idPropVal string, loadingType LoadingType) (IBusinessObject, error) {
	loadedBOs, errLoad := dbLoadOne(bloCtx.GetDaoContext(), idProp, idPropVal)

	if errLoad != nil {
		return nil, ErrorC(errLoad, "error while loading one instance of '%s' (%s)", idProp.ownerClass().base().name, idPropVal)
	}

	// TODO add post read

	return loadedBOs, nil
}

func DeleteBO(bloCtx BloContext, idProp IField, idPropVal string) (IBusinessObject, error) {
	loadedBOs, errLoad := dbRemoveOne(bloCtx.GetDaoContext(), idProp, idPropVal)

	if errLoad != nil {
		return nil, ErrorC(errLoad, "error while deleting one instance of '%s' (%s)", idProp.ownerClass().base().name, idPropVal)
	}

	// TODO add post read

	return loadedBOs, nil
}

func UpdateBO(bloCtx BloContext, input IBusinessObject, loadingType LoadingType) error {
	if errUpd := dbUpdate(bloCtx.GetDaoContext(), input); errUpd != nil {
		return ErrorC(errUpd, "error while updating one instance of '%T' (ID = %s)", input, input.GetID())
	}

	return nil
}
