package goald

// ----------------------------------------------------------------------------
// Base implem
// ----------------------------------------------------------------------------

// default implem for ISearchParamValues
type SearchParamValues struct {
	BusinessObject
	Page     int `json:"page,omitempty"     io:"in" desc:"the page number for pagination (default is 0)"`
	PageSize int `json:"pageSize,omitempty" io:"in" desc:"the number of items per page for pagination (default is controlled by the server)"`
}

func init() {
	modelSearchParamValues().SetNotPersisted()
	modelSearchParamValues().SetDescription("Base query parameters to list business objects")
}

// ----------------------------------------------------------------------------
// BLO part - should normally be in another file but it's small so it's here
// ----------------------------------------------------------------------------

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
