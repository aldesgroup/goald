package goald

import (
	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Business objects for containing search parameter values, which are involved in queries
// ------------------------------------------------------------------------------------------------

// particular business object used for containing search parameter values
type ISearchParamValues interface {
	IBusinessObject
	DoBeforeSearch(bloCtx BloContext) error
	getPage() int
	getPageSize() int
}

// ----------------------------------------------------------------------------
// Base implem
// ----------------------------------------------------------------------------

// default implem for ISearchParamValues
type SearchParamValues struct {
	BusinessObject
	Page     int `json:"page,omitempty"     io:"in" desc:"the page number for pagination (default is 0)"`
	PageSize int `json:"pageSize,omitempty" io:"in" desc:"the number of items per page for pagination (default is controlled by the server)"`
}

// default implem
func (this *SearchParamValues) DoBeforeSearch(bloCtx BloContext) error {
	// default implementation does nothing
	return nil
}

func (this *SearchParamValues) getPage() int {
	return this.Page
}

func (this *SearchParamValues) getPageSize() int {
	return this.PageSize
}

// unlike every other ISearchParamValues implementation, SearchParamValues is used "as-is" (i.e. not
// embedded within a generated subtype) whenever a BOTYPE is listed without any custom search params
// (see GenericHandleList), so it needs to fully implement IBusinessObject by itself, like a generated
// business object would.

func (this *SearchParamValues) GetModelName() utils.ModelName {
	return "SearchParamValues"
}

func (this *SearchParamValues) Clone(withFields, withRelationships bool) IBusinessObject {
	clone := &SearchParamValues{}
	clone.ID = this.ID

	if withFields {
		clone.Page = this.Page
		clone.PageSize = this.PageSize
	}

	return clone
}

func (this *SearchParamValues) GetValueAsString(propertyName string) string {
	switch propertyName {
	case "Page":
		return core.IntToString(this.Page)
	case "PageSize":
		return core.IntToString(this.PageSize)
	default:
		return "unknown property: " + propertyName
	}
}

func (this *SearchParamValues) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
	case "Page":
		this.Page = core.StringToInt(valueAsString, "Page")
	case "PageSize":
		this.PageSize = core.StringToInt(valueAsString, "PageSize")
	default:
		return Error("Unknown property: %T.%s", this, propertyName)
	}

	return nil
}

func (this *SearchParamValues) IsModelValid() error {
	return nil
}

func (this *SearchParamValues) RemoveCycles() {
	// no relationship to clean up
}

func (this *SearchParamValues) SetRelationshipValue(relationshipName string, value IBusinessObject) error {
	return Error("Unknown or non-single-valued relationship: %T.%s", this, relationshipName)
}

func (this *SearchParamValues) AddRelationshipValue(relationshipName string, value IBusinessObject) error {
	return Error("Unknown or non-multi-valued relationship: %T.%s", this, relationshipName)
}

func (this *SearchParamValues) ClearRelationshipValue(relationshipName string) error {
	return Error("Unknown or non-multi-valued relationship: %T.%s", this, relationshipName)
}

func (this *SearchParamValues) GetSingleRelationshipValue(relationshipName string) (IBusinessObject, error) {
	return nil, Error("Unknown or non-multi-valued relationship: %T.%s", this, relationshipName)
}

func (this *SearchParamValues) GetMultipleRelationshipValue(relationshipName string) ([]IBusinessObject, error) {
	return nil, Error("Unknown or non-multi-valued relationship: %T.%s", this, relationshipName)
}

func (this *SearchParamValues) DiffWith(other IBusinessObject, forLinks map[string]bool) (IBusinessObject, IBusinessObject) {
	return &SearchParamValues{}, &SearchParamValues{}
}

// ----------------------------------------------------------------------------
// Model & source registration - making SearchParamValues directly usable
// ----------------------------------------------------------------------------

type searchParamValuesModelSource struct {
	IBusinessObjectModelSource
}

func (this *searchParamValuesModelSource) NewObject() any {
	return &SearchParamValues{}
}

func (this *searchParamValuesModelSource) NewSlice() any {
	return &[]*SearchParamValues{}
}

func (this *searchParamValuesModelSource) AppendToSlice(slicePtr any, bObj any) any {
	s := slicePtr.(*[]*SearchParamValues)
	*s = append(*s, bObj.(*SearchParamValues))
	return s
}

func init() {
	RegisterModel("SearchParamValues", NewSearchParamValuesModel())

	In("goald").Register(&searchParamValuesModelSource{
		IBusinessObjectModelSource: NewBusinessObjectModelSource("goald", "SearchParamValues", "2026-09-16T00:00:00+00:00"),
	})
}

// ----------------------------------------------------------------------------
// BO model for it
// ----------------------------------------------------------------------------

// model associated with it
type ISearchParamValuesModel interface {
	IBusinessObjectModel
}

// constructor for this model
func NewSearchParamValuesModel() ISearchParamValuesModel {
	model := &struct {
		businessObjectModel
		page     *IntField
		pageSize *IntField
	}{
		businessObjectModel: businessObjectModel{
			fields:        map[string]IField{},
			relationships: map[string]*Relationship{},
			inNoDB:        true,
		},
	}

	model.page = AddIntField(model, "SearchParamValues", "Page", false)
	model.pageSize = AddIntField(model, "SearchParamValues", "PageSize", false)
	model.SetDescription("Base parameters for fetching several objects from the DB")

	return model
}
