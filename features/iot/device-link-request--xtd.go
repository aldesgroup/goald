// Generated file, do not edit!
package iot

import (
	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
	"github.com/aldesgroup/goald/features/utils"
)

// getting the name of the class for a DeviceLinkRequest, without using reflection
func (bo *DeviceLinkRequest) ClassName() utils.ClassName {
	return "DeviceLinkRequest"
}

// getting a property's value as a string, without using reflection
func (bo *DeviceLinkRequest) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "Model":
		return bo.Model
	case "Serial":
		return bo.Serial
	case "UserFullName":
		return bo.UserFullName
	case "UserID":
		return core.IntToString(bo.UserID)
	case "VerificationCode":
		return bo.VerificationCode
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *DeviceLinkRequest) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "Model":
		bo.Model = valueAsString
	case "Serial":
		bo.Serial = valueAsString
	case "UserFullName":
		bo.UserFullName = valueAsString
	case "UserID":
		bo.UserID = core.StringToInt(valueAsString, "UserID")
	case "VerificationCode":
		bo.VerificationCode = valueAsString
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *DeviceLinkRequest) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "ENTranslation":
		targetValue, ok := value.(*i18n.Translation)
		if !ok {
			return goald.Error("Expected a value of type '*i18n.Translation' for 'DeviceLinkRequest.ENTranslation', got %T", value)
		}
		bo.ENTranslation = targetValue
		return nil
	case "ForWho":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'DeviceLinkRequest.ForWho', got %T", value)
		}
		bo.ForWho = targetValue
		return nil
	case "MainContact":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'DeviceLinkRequest.MainContact', got %T", value)
		}
		bo.MainContact = targetValue
		return nil
	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *DeviceLinkRequest) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {
	case "Users":
		targetValue, ok := value.(goald.IUser)
		if !ok {
			return goald.Error("Expected a value of type 'goald.IUser' for 'DeviceLinkRequest.Users', got %T", value)
		}
		bo.Users = append(bo.Users, targetValue)
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *DeviceLinkRequest) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {
	case "Users":
		bo.Users = []goald.IUser{}
		return nil
	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *DeviceLinkRequest) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {
	case "Users":
		users := make([]goald.IBusinessObject, len(bo.Users))
		for i, target := range bo.Users {
			users[i] = target
		}
		return users, nil
	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// checking a business object's general validity, without using reflection
func (bo *DeviceLinkRequest) IsModelValid() error {
	if bo.Model == "" {
		return goald.Error("'Model' is mandatory and must have a non-zero value")
	}
	if bo.Serial == "" {
		return goald.Error("'Serial' is mandatory and must have a non-zero value")
	}

	return nil
}
