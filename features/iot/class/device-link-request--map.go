// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
	"github.com/aldesgroup/goald/features/iot"
)

// getting a property's value as a string, without using reflection
func (thisClass *DeviceLinkRequestClass) GetValueAsString(bo goald.IBusinessObject, propertyName string) string {
	switch propertyName {
	case "ID":
		return core.Int64ToString(int64(bo.(*iot.DeviceLinkRequest).ID))
	case "Model":
		return bo.(*iot.DeviceLinkRequest).Model
	case "Serial":
		return bo.(*iot.DeviceLinkRequest).Serial
	case "UserFullName":
		return bo.(*iot.DeviceLinkRequest).UserFullName
	case "UserID":
		return core.IntToString(bo.(*iot.DeviceLinkRequest).UserID)
	case "VerificationCode":
		return bo.(*iot.DeviceLinkRequest).VerificationCode
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (thisClass *DeviceLinkRequestClass) SetValueAsString(bo goald.IBusinessObject, propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.(*iot.DeviceLinkRequest).ID = goald.BObjID(core.StringToInt64(valueAsString, "(*iot.DeviceLinkRequest).ID"))
	case "Model":
		bo.(*iot.DeviceLinkRequest).Model = valueAsString
	case "Serial":
		bo.(*iot.DeviceLinkRequest).Serial = valueAsString
	case "UserFullName":
		bo.(*iot.DeviceLinkRequest).UserFullName = valueAsString
	case "UserID":
		bo.(*iot.DeviceLinkRequest).UserID = core.StringToInt(valueAsString, "(*iot.DeviceLinkRequest).UserID")
	case "VerificationCode":
		bo.(*iot.DeviceLinkRequest).VerificationCode = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (thisClass *DeviceLinkRequestClass) SetRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "ENTranslation":
		targetValue, ok := value.(*i18n.Translation)
		if !ok {
			return goald.Error("Expected a value of type '*i18n.Translation' for 'DeviceLinkRequest.ENTranslation', got %T", value)
		}
		bo.(*iot.DeviceLinkRequest).ENTranslation = targetValue
		return nil
	case "ForWho":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'DeviceLinkRequest.ForWho', got %T", value)
		}
		bo.(*iot.DeviceLinkRequest).ForWho = targetValue
		return nil
	case "MainContact":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'DeviceLinkRequest.MainContact', got %T", value)
		}
		bo.(*iot.DeviceLinkRequest).MainContact = targetValue
		return nil
	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (thisClass *DeviceLinkRequestClass) AddRelationshipValue(bo goald.IBusinessObject, relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "Users":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'DeviceLinkRequest.Users', got %T", value)
		}
		bo.(*iot.DeviceLinkRequest).Users = append(bo.(*iot.DeviceLinkRequest).Users, targetValue)
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (thisClass *DeviceLinkRequestClass) ClearRelationshipValue(bo goald.IBusinessObject, relationshipName string) error {
	switch relationshipName {
	case "Users":
		bo.(*iot.DeviceLinkRequest).Users = []goald.IUser{}
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}
