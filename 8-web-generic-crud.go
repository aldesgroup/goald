// ------------------------------------------------------------------------------------------------
// Here we implement some generic request handlers involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/hstatus"
)

// SetAutoCRUD sets up the generic CRUD endpoints for the given business object type (BOTYPE).
func SetAutoCRUD[BOTYPE IBusinessObject](group *EndpointGroup) {
	GenericHandleCreate[BOTYPE](group)

}

// SetAutoCRUDS sets up the generic CRUDS endpoints for the given business object type (BOTYPE)
func SetAutoCRUDS[BOTYPE IBusinessObject, QUERYPARAMVALUES IQueryParamsObject](group *EndpointGroup, defaultSearchQuery ...IQuery) {
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

	ep.Label(fmt.Sprintf("Create a new %s", ep.getResourceModel().getName()))
	ep.Description(fmt.Sprintf("Performs controls and saves the given %s instance in the database", ep.getResourceModel().getName()))
	ep.InGroup(group)

	return ep
}

// GenericHandleSearch creates a new endpoint for searching for instances of the given business object type (BOTYPE).
func GenericHandleSearch[BOTYPE IBusinessObject, QUERYPARAMVALUES IQueryParamsObject](group *EndpointGroup, query IQuery) *manyForOneEndpoint[QUERYPARAMVALUES, BOTYPE] {
	ep := GetManyWithParams(
		// new (anonym) handler function here
		func(webCtx WebContext, queryParamValues QUERYPARAMVALUES) ([]BOTYPE, hstatus.Code, string) {

			// calling the Business LOgic (BLO) for business object search
			results, errSearch := SearchBusinessObjects(webCtx.GetBloContext(), query, queryParamValues)
			if errSearch != nil {
				return nil, hstatus.InternalServerError, fmt.Sprintf("Failed searching for '%s' instances: %s", query.getModel().getName(), errSearch)
			}

			// return the found instances
			return allCleaned[BOTYPE](results), hstatus.OK, fmt.Sprintf("Found '%d' '%T' instances", len(results), queryParamValues)
		},
		// passing the loading type
		"")

	ep.Label(fmt.Sprintf("Search for %s instances", ep.getResourceModel().getName()))
	ep.Description(fmt.Sprintf("Searches for %s instances in the database using URL query parameters", ep.getResourceModel().getName()))
	ep.InGroup(group)

	return ep
}

// func GenericHandleRead[BOTYPE IBusinessObject](idProp IField, loadingType LoadingType) *oneForNoneEndpoint[BOTYPE] {
// 	ep := GetOne(
// 		// new (anonym) handler function here
// 		func(webCtx WebContext) (BOTYPE, hstatus.Code, string) {

// 			// boModel := GetModel[BOTYPE]()
// 			output, errRead := ReadBO(webCtx.GetBloContext(), idProp, webCtx.GetTargetRefOrID(), loadingType)
// 			if errRead != nil {
// 				return *new(BOTYPE), hstatus.InternalServerError,
// 					fmt.Sprintf("Failed reading '%s' instance '%s': %s", idProp.ownerModel().base().name, webCtx.GetTargetRefOrID(), errRead)
// 			}

// 			// return the   instance
// 			return output.(BOTYPE), hstatus.OK, fmt.Sprintf("Found the targeted '%T' instance", output)
// 		},
// 		// passing the loading type
// 		loadingType)

// 	ep.TargetWith(idProp)

// 	return ep
// }

// func GenericHandleRead[BOTYPE IBusinessObject](idProp IField, loadingType LoadingType) *oneForNoneEndpoint[BOTYPE] {
// 	ep := GetOne(
// 		// new (anonym) handler function here
// 		func(webCtx WebContext) (BOTYPE, hstatus.Code, string) {

// 			// boModel := GetModel[BOTYPE]()
// 			output, errRead := ReadBO(webCtx.GetBloContext(), idProp, webCtx.GetTargetRefOrID(), loadingType)
// 			if errRead != nil {
// 				return *new(BOTYPE), hstatus.InternalServerError,
// 					fmt.Sprintf("Failed reading '%s' instance '%s': %s", idProp.ownerModel().base().name, webCtx.GetTargetRefOrID(), errRead)
// 			}

// 			// return the   instance
// 			return output.(BOTYPE), hstatus.OK, fmt.Sprintf("Found the targeted '%T' instance", output)
// 		},
// 		// passing the loading type
// 		loadingType)

// 	ep.TargetWith(idProp)

// 	return ep
// }

// func GenericHandleUpdate[BOTYPE IBusinessObject](loadingType LoadingType) *oneForOneEndpoint[BOTYPE, BOTYPE] {
// 	ep := PutOne[BOTYPE](
// 		// new (anonym) handler function here
// 		func(webCtx WebContext, input BOTYPE) (BOTYPE, hstatus.Code, string) {
// 			if errUpdate := UpdateBO(webCtx.GetBloContext(), input, loadingType); errUpdate != nil {
// 				return *new(BOTYPE), hstatus.InternalServerError,
// 					fmt.Sprintf("Failed creating a new '%T' instance: %s", input, errUpdate)
// 			}

// 			return input, hstatus.OK, fmt.Sprintf("Update the given '%T' instance", input)
// 		},
// 		// passing the loading type
// 		loadingType)

// 	return ep
// }

// func GenericHandleDelete[BOTYPE IBusinessObject](idProp IField) *oneForNoneEndpoint[BOTYPE] {
// 	ep := DeleteOne(
// 		// new (anonym) handler function here
// 		func(webCtx WebContext) (BOTYPE, hstatus.Code, string) {
// 			// boModel := GetModel[BOTYPE]()
// 			output, errRead := DeleteBO(webCtx.GetBloContext(), idProp, webCtx.GetTargetRefOrID())
// 			if errRead != nil {
// 				return *new(BOTYPE), hstatus.InternalServerError,
// 					fmt.Sprintf("Failed reading '%s' instance '%s': %s", idProp.ownerModel().base().name, webCtx.GetTargetRefOrID(), errRead)
// 			}

// 			return output.(BOTYPE), hstatus.OK, fmt.Sprintf("Deleted the targeted '%T' instance", output)
// 		})

// 	ep.TargetWith(idProp)

// 	return ep
// }
