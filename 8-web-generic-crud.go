// ------------------------------------------------------------------------------------------------
// Here we implement some generic request handlers involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"strconv"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/hstatus"
)

// SetAutoCRUD sets up the generic CRUD endpoints for the given business object type (BOTYPE).
func SetAutoCRUD[BOTYPE IBusinessObject](group *EndpointGroup) {
	GenericHandleCreate[BOTYPE](group)
	GenericHandleRead[BOTYPE](group)
}

// SetAutoCRUDS sets up the generic CRUDS endpoints for the given business object type (BOTYPE)
func SetAutoCRUDS[BOTYPE IBusinessObject, QUERYPARAMVALUES ISearchParamValues](group *EndpointGroup, defaultSearchQuery ...*query[BOTYPE, QUERYPARAMVALUES]) {
	SetAutoCRUD[BOTYPE](group)
	if len(defaultSearchQuery) > 1 {
		core.PanicMsg("SetAutoCRUD: only 1 default search query is allowed, but %d were provided", len(defaultSearchQuery))
	}
	if len(defaultSearchQuery) > 0 {
		GenericHandleSearch[BOTYPE, QUERYPARAMVALUES](group, defaultSearchQuery[0])
	}
}

// GenericHandleCreate creates a new endpoint for creating a new instance of the given business object type (BOTYPE).
func GenericHandleCreate[BOTYPE IBusinessObject](group *EndpointGroup) *oneForOneEndpoint[BOTYPE, BOTYPE] {
	ep := PostOneGetOne(
		// new (anonym) handler function here
		func(webCtx WebContext, input BOTYPE) (BOTYPE, hstatus.Code, string) {
			// checking the input first
			if errValidate := input.IsModelValid(); errValidate != nil {
				return cleaned(input), hstatus.BadRequest, fmt.Sprintf("Invalid input for creating a new '%T' instance: %s", input, errValidate)
			}

			// calling the Business LOgic (BLO) for business object creation
			if errCreate := CreateBusinessObjects(webCtx.GetBloContext(), input); errCreate != nil {
				return cleaned(input), hstatus.InternalServerError,
					fmt.Sprintf("Failed creating a new '%T' instance: %s", input, errCreate)
			}

			// return the created instance
			return cleaned(input), hstatus.Created, fmt.Sprintf("Created a new '%T' instance", input)
		},
		// passing the loading type
		"")

	ep.Label(fmt.Sprintf("Create a new %s", ep.getResourceModel().GetName()))
	ep.Description(fmt.Sprintf("Performs controls and saves the given %s instance in the database", ep.getResourceModel().GetName()))
	ep.InGroup(group)

	return ep
}

// GenericHandleSearch creates a new endpoint for searching for instances of the given business object type (BOTYPE).
func GenericHandleSearch[BOTYPE IBusinessObject, QUERYPARAMVALUES ISearchParamValues](group *EndpointGroup, query IQuery) *manyForOneEndpoint[QUERYPARAMVALUES, BOTYPE] {
	ep := GetManyWithParams(
		// new (anonym) handler function here
		func(webCtx WebContext, queryParamValues QUERYPARAMVALUES) ([]BOTYPE, hstatus.Code, string) {

			// calling the Business LOgic (BLO) for business object search
			results, errSearch := SearchBusinessObjects(webCtx.GetBloContext(), query, queryParamValues)
			if errSearch != nil {
				return nil, hstatus.InternalServerError, fmt.Sprintf("Failed searching for '%s' instances: %s", query.getSearchedObjectsModel().GetName(), errSearch)
			}

			// return the found instances
			return allCleaned[BOTYPE](results), hstatus.OK, fmt.Sprintf("Found '%d' '%T' instances", len(results), queryParamValues)
		},
		// passing the loading type
		"")

	ep.Label(fmt.Sprintf("Search for %s instances", ep.getResourceModel().GetName()))
	ep.Description(fmt.Sprintf("Searches for %s instances in the database using URL query parameters", ep.getResourceModel().GetName()))
	ep.InGroup(group)

	return ep
}

// GenericHandleRead creates a new endpoint for reading one instance of the given business object type (BOTYPE).
func GenericHandleRead[BOTYPE IBusinessObject](group *EndpointGroup) *oneForNoneEndpoint[BOTYPE] {

	boModelName := (*new(BOTYPE)).GetModelName()
	boModel := modelFor(boModelName, true)

	ep := GetOne(
		// new (anonym) handler function here
		func(webCtx WebContext) (BOTYPE, hstatus.Code, string) {
			// retrieving the ID from the URL path
			idAsString := webCtx.GetResourceRefOrID()

			// checking the ID is valid
			var zero BOTYPE
			if idAsString == "" {
				return zero, hstatus.BadRequest, fmt.Sprintf("Missing ID in request path for reading a '%s' instance", boModelName)
			}
			id, errConv := strconv.ParseInt(idAsString, 10, 64)
			if errConv != nil {
				return zero, hstatus.BadRequest, fmt.Sprintf("Invalid ID '%s' in request path for reading a '%s' instance: %s", idAsString, boModelName, errConv)
			}

			// the expected output
			output := NewBusinessObject(string(boModelName), id).(BOTYPE)

			// calling the Business LOgic (BLO) for business object reading
			if errRead := ReadBusinessObjects(webCtx.GetBloContext(), output); errRead != nil {
				return zero, hstatus.InternalServerError, fmt.Sprintf("Failed reading '%s' instance with ID '%d': %s", boModelName, id, errRead)
			}

			return cleaned(output), hstatus.OK, fmt.Sprintf("Read a '%s' instance with ID '%d'", boModelName, id)
		},
		// passing the loading type
		"")

	ep.TargetWith(boModel.getIdField())
	ep.Label(fmt.Sprintf("Read a %s instance", ep.getResourceModel().GetName()))
	ep.Description(fmt.Sprintf("Reads a %s instance from the database using its ID", ep.getResourceModel().GetName()))
	ep.InGroup(group)

	return ep
}
