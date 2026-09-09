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

func NewDeviceBootstrapRequest(id goald.BObjID) *DeviceBootstrapRequest {
	// TODO use sync.Pool?
	newDeviceBootstrapRequest := &DeviceBootstrapRequest{}
	newDeviceBootstrapRequest.ID = id

	return newDeviceBootstrapRequest
}

func GetDeviceBootstrapRequestFrom(cache *goald.BObjCache, id goald.BObjID) *DeviceBootstrapRequest {
	if cachedDeviceBootstrapRequest := cache.Get("DeviceBootstrapRequest", id); cachedDeviceBootstrapRequest != nil {
		return cachedDeviceBootstrapRequest.(*DeviceBootstrapRequest)
	}

	return nil
}

func CachedOrNewDeviceBootstrapRequest(cache *goald.BObjCache, id goald.BObjID) *DeviceBootstrapRequest {
	if cachedDeviceBootstrapRequest := GetDeviceBootstrapRequestFrom(cache, id); cachedDeviceBootstrapRequest != nil {
		return cachedDeviceBootstrapRequest
	}

	return cache.Set(NewDeviceBootstrapRequest(id))
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a DeviceBootstrapRequest, without using reflection
func (bo *DeviceBootstrapRequest) GetModelName() utils.ModelName {
	return "DeviceBootstrapRequest"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *DeviceBootstrapRequest) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Creation":
		return core.DateToString(bo.Creation)
	case "FactoryCertPEM":
		return bo.FactoryCertPEM
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "PayloadB64":
		return bo.PayloadB64
	case "SignatureB64":
		return bo.SignatureB64
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *DeviceBootstrapRequest) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
	case "FactoryCertPEM":
		bo.FactoryCertPEM = valueAsString
	case "ID":
		bo.ID = goald.BObjID(core.StringToInt64(valueAsString, "ID"))
	case "PayloadB64":
		bo.PayloadB64 = valueAsString
	case "SignatureB64":
		bo.SignatureB64 = valueAsString
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
func (bo *DeviceBootstrapRequest) SetRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *DeviceBootstrapRequest) AddRelationshipValue(relationshipName string, value goald.IBusinessObject) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *DeviceBootstrapRequest) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *DeviceBootstrapRequest) GetSingleRelationshipValue(relationshipName string) (goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *DeviceBootstrapRequest) GetMultipleRelationshipValue(relationshipName string) ([]goald.IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, goald.Error("Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *DeviceBootstrapRequest) IsModelValid() error {
	if bo.FactoryCertPEM == "" {
		return goald.Error("'FactoryCertPEM' is mandatory and must have a non-zero value")
	}
	if bo.PayloadB64 == "" {
		return goald.Error("'PayloadB64' is mandatory and must have a non-zero value")
	}
	if bo.SignatureB64 == "" {
		return goald.Error("'SignatureB64' is mandatory and must have a non-zero value")
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *DeviceBootstrapRequest) RemoveCycles() {
}
