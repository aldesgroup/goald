// ------------------------------------------------------------------------------------------------
// Here we implement the generic data access instructions involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import (
	"reflect"

	"github.com/google/uuid"
)

var mockDatabase map[string]IBusinessObject = map[string]IBusinessObject{}

func dbInsert(daoCtx DaoContext, bObj IBusinessObject) error {
	uuid, errUuid := uuid.NewRandom()
	if errUuid != nil {
		return ErrorC(errUuid, "could not generate a new UUID")
	}

	bObj.setID(uuid.String())
	bObj.setClassName(reflect.TypeOf(bObj).Elem().Name())

	mockDatabase[string(bObj.GetID())] = bObj

	return nil
}

func dbLoadList(daoCtx DaoContext, boClass IBusinessObjectClass) (result []IBusinessObject, err error) {
	for _, bObj := range mockDatabase {
		if boClass == bObj.Class() {
			result = append(result, bObj)
		}
	}

	return
}

func dbLoadOne(daoCtx DaoContext, idProp IField, idPropVal string) (result IBusinessObject, err error) {
	for _, bObj := range mockDatabase {
		if idProp.ownerClass() == bObj.Class() && idPropVal == idProp.StringValue(bObj) {
			return bObj, nil
		}
	}

	return nil, Error("No '%s' found with '%s = %s'", idProp.ownerClass().base().className, idProp.getName(), idPropVal)
}

func dbRemoveOne(daoCtx DaoContext, idProp IField, idPropVal string) (result IBusinessObject, err error) {
	for _, bObj := range mockDatabase {
		if idProp.ownerClass() == bObj.Class() && idPropVal == idProp.StringValue(bObj) {
			delete(mockDatabase, string(bObj.GetID()))
			return bObj, nil
		}
	}

	return nil, Error("No '%s' found with '%s = %s'", idProp.ownerClass().base().className, idProp.getName(), idPropVal)
}

func dbUpdate(daoCtx DaoContext, input IBusinessObject) error {
	instore := mockDatabase[string(input.GetID())]
	if instore == nil {
		return Error("No object exists with ID %s", input.GetID())
	}

	input.setClassName(instore.getClassName())

	mockDatabase[string(input.GetID())] = input

	return nil
}
