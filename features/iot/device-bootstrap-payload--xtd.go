// Generated file, do not edit!
package iot

import (
	"github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Instantiation / cache retrieval
// ------------------------------------------------------------------------------------------------

func NewDeviceBootstrapPayload(id goald.BObjID) *DeviceBootstrapPayload {
	// TODO use sync.Pool?
	newDeviceBootstrapPayload := &DeviceBootstrapPayload{}
	newDeviceBootstrapPayload.ID = id

	return newDeviceBootstrapPayload
}

func GetDeviceBootstrapPayloadFrom(cache *goald.BObjCache, id goald.BObjID) *DeviceBootstrapPayload {
	if cachedDeviceBootstrapPayload := cache.Get("DeviceBootstrapPayload", id); cachedDeviceBootstrapPayload != nil {
		return cachedDeviceBootstrapPayload.(*DeviceBootstrapPayload)
	}

	return nil
}

func CachedOrNewDeviceBootstrapPayload(cache *goald.BObjCache, id goald.BObjID) *DeviceBootstrapPayload {
	if cachedDeviceBootstrapPayload := GetDeviceBootstrapPayloadFrom(cache, id); cachedDeviceBootstrapPayload != nil {
		return cachedDeviceBootstrapPayload
	}

	return cache.Set(NewDeviceBootstrapPayload(id))
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a DeviceBootstrapPayload, without using reflection
func (bo *DeviceBootstrapPayload) GetModelName() utils.ModelName {
	return "DeviceBootstrapPayload"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *DeviceBootstrapPayload) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Creation":
		return core.DateToString(bo.Creation)
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "Model":
		return bo.Model
	case "Nonce":
		return bo.Nonce
	case "Serial":
		return bo.Serial
	case "Timestamp":
		return core.Int64ToString(bo.Timestamp)
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *DeviceBootstrapPayload) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "Model":
		bo.Model = valueAsString
	case "Nonce":
		bo.Nonce = valueAsString
	case "Serial":
		bo.Serial = valueAsString
	case "Timestamp":
		bo.Timestamp = core.StringToInt64(valueAsString, "Timestamp")
	}

	return goald.Error("Unknown property: %T.%s", bo, propertyName)
}

// ------------------------------------------------------------------------------------------------
// Explicit relationship access
// ------------------------------------------------------------------------------------------------


// ------------------------------------------------------------------------------------------------
// Generic relationship access
// ------------------------------------------------------------------------------------------------

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *DeviceBootstrapPayload) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *DeviceBootstrapPayload) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *DeviceBootstrapPayload) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *DeviceBootstrapPayload) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *DeviceBootstrapPayload) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *DeviceBootstrapPayload) IsModelValid() error {
	if bo.Model == "" {
		return goald.Error("'Model' is mandatory and must have a non-zero value")
	}
	if bo.Nonce == "" {
		return goald.Error("'Nonce' is mandatory and must have a non-zero value")
	}
	if bo.Serial == "" {
		return goald.Error("'Serial' is mandatory and must have a non-zero value")
	}
	if bo.Timestamp == 0 {
		return goald.Error("'Timestamp' is mandatory and must have a non-zero value")
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *DeviceBootstrapPayload) RemoveCycles() {
}
