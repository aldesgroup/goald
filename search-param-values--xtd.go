// Generated file, do not edit!
package goald

import (
	//$$imports$$
	"sync"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// -----------------------------------------------------------------------------
// The SOURCE part
// -----------------------------------------------------------------------------

type SearchParamValuesModelSource struct {
	IBusinessObjectModelSource
}

func ForSearchParamValues(srcPath, lastMod string) IBusinessObjectModelSource {
	return &SearchParamValuesModelSource{IBusinessObjectModelSource: NewBusinessObjectModelSource(srcPath, "SearchParamValues", lastMod)}
}

func (this *SearchParamValuesModelSource) NewObject() any {
	return &SearchParamValues{}
}

func (this *SearchParamValuesModelSource) NewSlice() any {
	return &[]*SearchParamValues{}
}

func (this *SearchParamValuesModelSource) AppendToSlice(slicePtr any, bObj any) any {
	s := slicePtr.(*[]*SearchParamValues)
	*s = append(*s, bObj.(*SearchParamValues))
	return s
}

// -----------------------------------------------------------------------------
// The MODEL part
// -----------------------------------------------------------------------------

// static, reflect-free access to the definition of the SearchParamValues model
type SearchParamValuesModel struct {
	IBusinessObjectModel
	page     *IntField
	pageSize *IntField
}

// this is the main way to refer to the SearchParamValues model in the applicative code
func modelSearchParamValues() *SearchParamValuesModel {
	searchParamValuesOnce.Do(func() {
		searchParamValues = NewSearchParamValuesModel()

		// this helps dynamically access to the SearchParamValues model
		RegisterModel("SearchParamValues", searchParamValues)
	})

	return searchParamValues
}

// internal variables
var searchParamValues *SearchParamValuesModel
var searchParamValuesOnce sync.Once

// fully describing each of this model's properties & relationships
func NewSearchParamValuesModel() *SearchParamValuesModel {
	thisModel := &SearchParamValuesModel{IBusinessObjectModel: NewBusinessObjectModel()}
	thisModel.page = AddIntField(thisModel, "SearchParamValues", "Page", false)
	thisModel.pageSize = AddIntField(thisModel, "SearchParamValues", "PageSize", false)

	return thisModel
}

// making sure the SearchParamValues model exists at app startup
func init() {
	modelSearchParamValues()
}

// accessing all the SearchParamValues model's properties and relation	ships

func (S *SearchParamValuesModel) Page() *IntField {
	return S.page
}

func (S *SearchParamValuesModel) PageSize() *IntField {
	return S.pageSize
}

// ------------------------------------------------------------------------------------------------
// Instantiation / cache retrieval
// ------------------------------------------------------------------------------------------------

func NewSearchParamValues(id BObjID) *SearchParamValues {
	// TODO use sync.Pool?
	newSearchParamValues := &SearchParamValues{}
	newSearchParamValues.ID = id

	return newSearchParamValues
}

func GetSearchParamValuesFrom(cache *BObjCache, id BObjID) *SearchParamValues {
	if cachedSearchParamValues := cache.Get("SearchParamValues", id); cachedSearchParamValues != nil {
		return cachedSearchParamValues.(*SearchParamValues)
	}

	return nil
}

func CachedOrNewSearchParamValues(cache *BObjCache, id BObjID) *SearchParamValues {
	if cachedSearchParamValues := GetSearchParamValuesFrom(cache, id); cachedSearchParamValues != nil {
		return cachedSearchParamValues
	}

	return cache.Set(NewSearchParamValues(id))
}

// ------------------------------------------------------------------------------------------------
// Cloning
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a SearchParamValues, without using reflection
func (bo *SearchParamValues) Clone(withFields, withRelationships bool) IBusinessObject {
	clone := &SearchParamValues{}
	clone.ID = bo.ID

	if withFields {
		clone.Creation = bo.Creation
		clone.Modification = bo.Modification
		clone.Page = bo.Page
		clone.PageSize = bo.PageSize
	}

	if withRelationships {

	}

	return clone
}

// ------------------------------------------------------------------------------------------------
// Identification
// ------------------------------------------------------------------------------------------------

// getting the name of the model for a SearchParamValues, without using reflection
func (bo *SearchParamValues) GetModelName() utils.ModelName {
	return "SearchParamValues"
}

// ------------------------------------------------------------------------------------------------
// Property values <-> string conversion
// ------------------------------------------------------------------------------------------------

// getting a property's value as a string, without using reflection
func (bo *SearchParamValues) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Creation":
		return core.DateToString(bo.Creation)
	case "ID":
		return core.Int64ToString(int64(bo.ID))
	case "Modification":
		return core.DateToString(bo.Modification)
	case "Page":
		return core.IntToString(bo.Page)
	case "PageSize":
		return core.IntToString(bo.PageSize)
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (bo *SearchParamValues) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Creation":
		bo.Creation = core.StringToDate(valueAsString, "Creation")
	case "ID":
		bo.ID = BObjID(core.StringToInt64(valueAsString, "ID"))
	case "Modification":
		bo.Modification = core.StringToDate(valueAsString, "Modification")
	case "Page":
		bo.Page = core.StringToInt(valueAsString, "Page")
	case "PageSize":
		bo.PageSize = core.StringToInt(valueAsString, "PageSize")
	}

	return Error("[SetValueAsString] Unknown property: %T.%s", bo, propertyName)
}

// ------------------------------------------------------------------------------------------------
// Explicit relationship access
// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
// Generic relationship access
// ------------------------------------------------------------------------------------------------

// setting this SearchParamValues's parent
func (bo *SearchParamValues) SetParent(parent IBusinessObject) {
	// no parent for this model
}

// setting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *SearchParamValues) SetRelationshipValue(relationshipName string, value IBusinessObject) error {
	switch relationshipName {

	}

	return Error("[SetRelationshipValue] Unknown or non-single-valued relationship: %T.%s", bo, relationshipName)
}

// appending a target to a multi-valued relationship, given the relationship's name, without using reflection
func (bo *SearchParamValues) AddRelationshipValue(relationshipName string, value IBusinessObject) error {
	switch relationshipName {

	}

	return Error("[AddRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// resetting a multi-valued relationship to an empty slice, given the relationship's name, without using reflection
func (bo *SearchParamValues) ClearRelationshipValue(relationshipName string) error {
	switch relationshipName {

	}

	return Error("[ClearRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a single-valued relationship's target, given the relationship's name, without using reflection
func (bo *SearchParamValues) GetSingleRelationshipValue(relationshipName string) (IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, Error("[GetSingleRelationshipValue]Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// getting a multi-valued relationship's targets, given the relationship's name, without using reflection
func (bo *SearchParamValues) GetMultipleRelationshipValue(relationshipName string) ([]IBusinessObject, error) {
	switch relationshipName {

	}

	return nil, Error("[GetMultipleRelationshipValue] Unknown or non-multi-valued relationship: %T.%s", bo, relationshipName)
}

// ------------------------------------------------------------------------------------------------
// Model validity check
// ------------------------------------------------------------------------------------------------

// checking a business object's general validity, without using reflection
func (bo *SearchParamValues) IsModelValid() error {

	return nil
}

// ------------------------------------------------------------------------------------------------
// Diffing
// ------------------------------------------------------------------------------------------------

// Creates 2 synthetic instances gathering the added and removed relationships
func (bo *SearchParamValues) DiffWith(other IBusinessObject, forLinks map[string]bool) (IBusinessObject, IBusinessObject) {
	added := NewSearchParamValues(bo.ID)
	removed := NewSearchParamValues(bo.ID)

	return added, removed
}

// ------------------------------------------------------------------------------------------------
// Misc utils
// ------------------------------------------------------------------------------------------------

// removing any cycles from the business object, without using reflection
func (bo *SearchParamValues) RemoveCycles() {
}

// setting the models names on all the business objects associated with this one
func (bo *SearchParamValues) SetModelNames() {
}
