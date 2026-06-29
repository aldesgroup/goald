// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/auth"
)

// getting a property's value as a string, without using reflection
func (thisClass *UserClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "Email":
		return bo.(*auth.User).Email
	case "FirstName":
		return bo.(*auth.User).FirstName
	case "ID":
		return core.Int64ToString(int64(bo.(*auth.User).ID))
	case "LastName":
		return bo.(*auth.User).LastName
	case "Password":
		return bo.(*auth.User).Password
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *UserClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "Email":
		bo.(*auth.User).Email = valueAsString
	case "FirstName":
		bo.(*auth.User).FirstName = valueAsString
	case "ID":
		bo.(*auth.User).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*auth.User).ID"))
	case "LastName":
		bo.(*auth.User).LastName = valueAsString
	case "Password":
		bo.(*auth.User).Password = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}
