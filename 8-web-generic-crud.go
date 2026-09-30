// ------------------------------------------------------------------------------------------------
// Here we implement some generic request handlers involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import (
	"errors"
	"fmt"
	"strconv"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/hstatus"
)

// doSetAutoCRUD sets up the generic CRUD endpoints for the given business object type (BOTYPE).
func doSetAutoCRUD[BOTYPE IBusinessObject, QUERYPARAMVALUES ISearchParamValues](addHandleList bool, group *EndpointGroup, loadReadBObjWith ...ILoadingConfig) {

	// dealing with the optinal loading config for reading business objects
	var readLoadingCfg ILoadingConfig
	if len(loadReadBObjWith) > 1 {
		core.PanicMsg("SetAutoCRUD: only 1 default loading config for read business objects is allowed, but %d were provided", len(loadReadBObjWith))
	}
	if len(loadReadBObjWith) > 0 {
		readLoadingCfg = loadReadBObjWith[0]
	} else {
		readLoadingCfg = modelFor((*new(BOTYPE)).GetModelName(), true).ReadWithFirstLayer()
	}

	GenericHandleCreate[BOTYPE](group)
	GenericHandleRead[BOTYPE](group, readLoadingCfg)
	GenericHandleUpdate[BOTYPE](group)
	if addHandleList {
		GenericHandleList[BOTYPE, QUERYPARAMVALUES](group)
	}
}

// SetAutoCRUD sets up the generic CRUD endpoints for the given business object type (BOTYPE),
// plus a basic list endpoint with pagination parameters
func SetAutoCRUD[BOTYPE IBusinessObject](group *EndpointGroup, loadReadBObjWith ...ILoadingConfig) {
	doSetAutoCRUD[BOTYPE, *SearchParamValues](true, group, loadReadBObjWith...)
}

// SetAutoCRUDL sets up the generic CRUD endpoints for the given business object type (BOTYPE),
// plus a basic list endpoint with a specific query parameter type
func SetAutoCRUDL[BOTYPE IBusinessObject, QV ISearchParamValues](group *EndpointGroup, loadReadBObjWith ...ILoadingConfig) {
	doSetAutoCRUD[BOTYPE, QV](true, group, loadReadBObjWith...)
}

// SetAutoCRUDS sets up the generic CRUDS endpoints for the given business object type (BOTYPE)
func SetAutoCRUDS[BO IBusinessObject, QV ISearchParamValues](group *EndpointGroup, searchQuery *query[BO, QV], loadReadBObjWith ...ILoadingConfig) {
	doSetAutoCRUD[BO, QV](false, group, loadReadBObjWith...)
	GenericHandleSearch[BO, QV](group, searchQuery)
}

// GenericHandleSearch creates a new endpoint for searching for instances of the given business object type (BOTYPE).
func GenericHandleSearch[BOTYPE IBusinessObject, QUERYPARAMVALUES ISearchParamValues](group *EndpointGroup, query IQuery) *manyForOneEndpoint[QUERYPARAMVALUES, BOTYPE] {
	ep := GetManyWithParams(
		// new (anonym) handler function here
		func(webCtx WebContext, queryParamValues QUERYPARAMVALUES) ([]BOTYPE, hstatus.Code, string) {

			// calling the Business LOgic (BLO) for business object search
			results, errSearch := SearchBusinessObjects(webCtx.GetBloContext(), query, queryParamValues, query.getLoadingConfig())
			if errSearch != nil {
				return nil, hstatus.InternalServerError, fmt.Sprintf("Failed searching for '%s' instances: %s", query.getSearchedObjectsModel().GetName(), errSearch)
			}

			// found nothing?
			if len(results) == 0 {
				return nil, hstatus.NotFound, fmt.Sprintf("No '%s' instances found for the given search parameters", query.getSearchedObjectsModel().GetName())
			}

			// return the found instances
			return allCastAndCleaned[BOTYPE](results), hstatus.OK, fmt.Sprintf("Found %d '%s' instances", len(results), query.getSearchedObjectsModel().GetName())
		},
		// passing the loading type
		query.getLoadingConfig())

	ep.Label(fmt.Sprintf("Search for %s instances", ep.getResourceModel().GetName()))
	ep.Description(fmt.Sprintf("Searches for %s instances in the database using URL query parameters", ep.getResourceModel().GetName()))
	ep.InGroup(group)

	return ep
}

// GenericHandleCreate creates a new endpoint for creating new instances of the given business object type (BOTYPE).
func GenericHandleCreate[BOTYPE IBusinessObject](group *EndpointGroup) *manyForManyEndpoint[BOTYPE, BOTYPE] {
	ep := PostManyGetMany(
		// new (anonym) handler function here
		func(webCtx WebContext, inputs []BOTYPE) ([]BOTYPE, hstatus.Code, string) {
			// early return if no inputs were provided
			if len(inputs) == 0 {
				return inputs, hstatus.BadRequest, "No input provided for creating new instances"
			}

			// checking the input first
			for _, input := range inputs {
				if errValidate := input.IsModelValid(); errValidate != nil {
					return inputs, hstatus.BadRequest, fmt.Sprintf("Invalid input for creating a new '%T' instance: %s", inputs[0], errValidate)
				}
			}

			// calling the Business LOgic (BLO) for business object creation
			if errCreate := CreateBusinessObjects(webCtx.GetBloContext(), inputs...); errCreate != nil {
				return allCleaned(inputs), hstatus.InternalServerError, fmt.Sprintf("Failed creating new '%T' instances: %s", inputs[0], errCreate)
			}

			// return the created instance after cleaning cycles
			return allCleaned(inputs), hstatus.Created, fmt.Sprintf("Created %d new '%T' instances", len(inputs), inputs[0])
		},
		// not loading anything in this
		nil)

	ep.Label(fmt.Sprintf("Create a new %s", ep.getResourceModel().GetName()))
	ep.Description(fmt.Sprintf("Performs controls and saves the given %s instance in the database", ep.getResourceModel().GetName()))
	ep.InGroup(group)

	return ep
}

// GenericHandleRead creates a new endpoint for reading one instance of the given business object type (BOTYPE).
func GenericHandleRead[BOTYPE IBusinessObject](group *EndpointGroup, loadReadBObjWith ILoadingConfig) *oneForNoneEndpoint[BOTYPE] {

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
			output := NewBusinessObject(boModelName, BObjID(id)).(BOTYPE)

			// calling the Business LOgic (BLO) for business object reading
			if errRead := ReadBusinessObjects(webCtx.GetBloContext(), loadReadBObjWith, output); errRead != nil {
				if errors.Is(errRead, ErrBoNotFound) {
					return zero, hstatus.NotFound, fmt.Sprintf("Could not find a '%s' instance with ID '%d'", boModelName, id)
				}

				return zero, hstatus.InternalServerError, fmt.Sprintf("Failed reading '%s' instance with ID '%d': %s", boModelName, id, errRead)
			}

			return cleaned(output), hstatus.OK, fmt.Sprintf("Read a '%s' instance with ID '%d'", boModelName, id)
		},
		// passing the loading config
		loadReadBObjWith)

	ep.TargetWith(boModel.getIdField())
	ep.Label(fmt.Sprintf("Read a %s instance", ep.getResourceModel().GetName()))
	ep.Description(fmt.Sprintf("Reads a %s instance from the database using its ID", ep.getResourceModel().GetName()))
	ep.InGroup(group)

	return ep
}

// GenericHandleUpdate creates a new endpoint for updating existing instances of the given business object type (BOTYPE).
func GenericHandleUpdate[BOTYPE IBusinessObject](group *EndpointGroup) *manyForManyEndpoint[BOTYPE, BOTYPE] {
	ep := PutManyGetMany(
		// new (anonym) handler function here
		func(webCtx WebContext, inputs []BOTYPE) ([]BOTYPE, hstatus.Code, string) {
			// early return if no inputs were provided
			if len(inputs) == 0 {
				return inputs, hstatus.BadRequest, "No input provided for updating instances"
			}

			// checking the input first
			for _, input := range inputs {
				if errValidate := input.IsModelValid(); errValidate != nil {
					return inputs, hstatus.BadRequest, fmt.Sprintf("Invalid input for creating a new '%T' instance: %s", input, errValidate)
				}
			}

			// calling the Business LOgic (BLO) for business object update
			if errUpdate := UpdateBusinessObjects(webCtx.GetBloContext(), inputs...); errUpdate != nil {
				return allCleaned(inputs), hstatus.InternalServerError, fmt.Sprintf("Failed updating '%T' instances: %s", inputs[0], errUpdate)
			}

			// return the created instance
			return allCleaned(inputs), hstatus.OK, fmt.Sprintf("Updated %d '%T' instances", len(inputs), inputs[0])
		},
		// not loading anything in this
		nil)

	ep.Label(fmt.Sprintf("Update existing %s instances", ep.getResourceModel().GetName()))
	ep.Description(fmt.Sprintf("Performs controls and updates the given %s instances in the database", ep.getResourceModel().GetName()))
	ep.InGroup(group)

	return ep
}

// GenericHandleList creates a new endpoint for searching for instances of the given business object type (BOTYPE).
func GenericHandleList[BOTYPE IBusinessObject, QUERYPARAMVALUES ISearchParamValues](group *EndpointGroup) *manyForOneEndpoint[QUERYPARAMVALUES, BOTYPE] {
	boModelName := (*new(BOTYPE)).GetModelName()
	boModel := modelFor(boModelName, true)
	loadingConf := core.IfThenElse(boModel.getListLoadingConfig() != nil, boModel.getListLoadingConfig(), boModel.ReadNoRelationship())
	query := Find[BOTYPE, QUERYPARAMVALUES](loadingConf)

	ep := GetManyWithParams(
		// new (anonym) handler function here
		func(webCtx WebContext, queryParamValues QUERYPARAMVALUES) ([]BOTYPE, hstatus.Code, string) {

			// calling the Business LOgic (BLO) for business object search
			results, errSearch := SearchBusinessObjects(webCtx.GetBloContext(), query, queryParamValues, query.getLoadingConfig())
			if errSearch != nil {
				return nil, hstatus.InternalServerError, fmt.Sprintf("Failed retrieving '%s' instances: %s", query.getSearchedObjectsModel().GetName(), errSearch)
			}

			// found nothing?
			if len(results) == 0 {
				return nil, hstatus.NotFound, fmt.Sprintf("No '%s' instances found for the given parameters", query.getSearchedObjectsModel().GetName())
			}

			// return the found instances
			return allCastAndCleaned[BOTYPE](results), hstatus.OK, fmt.Sprintf("Found %d '%s' instances", len(results), query.getSearchedObjectsModel().GetName())
		},
		// passing the loading type
		query.getLoadingConfig())

	ep.Label(fmt.Sprintf("List %s instances", ep.getResourceModel().GetName()))
	ep.Description(fmt.Sprintf("Retrieves %s instances from the database", ep.getResourceModel().GetName()))
	ep.InGroup(group)

	return ep
}
