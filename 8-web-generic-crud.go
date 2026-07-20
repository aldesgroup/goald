// ------------------------------------------------------------------------------------------------
// Here we implement some generic request handlers involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"

	"github.com/aldesgroup/goald/features/hstatus"
)

// SetAutoCRUD sets up the generic CRUD endpoints for the given business object type (BOTYPE).
func SetAutoCRUD[BOTYPE IBusinessObject](group *EndpointGroup) {
	GenericHandleCreate[BOTYPE](group)
}

// GenericHandleCreate creates a new endpoint for creating a new instance of the given business object type (BOTYPE).
func GenericHandleCreate[BOTYPE IBusinessObject](group *EndpointGroup) *oneForOneEndpoint[BOTYPE, BOTYPE] {
	ep := PostOneGetOne(
		// new (anonym) handler function here
		func(webCtx WebContext, input BOTYPE) (BOTYPE, hstatus.Code, string) {

			// calling the Business LOgic (BLO) for business object creation
			if errCreate := CreateBO(webCtx.GetBloContext(), input); errCreate != nil {
				return input, hstatus.InternalServerError,
					fmt.Sprintf("Failed creating a new '%T' instance: %s", input, errCreate)
			}

			// return the created instance
			return input, hstatus.Created, fmt.Sprintf("Created a new '%T' instance", input)
		},
		// passing the loading type
		"")

	ep.Label(fmt.Sprintf("Create a new %s", ep.getResourceClass()))
	ep.Description(fmt.Sprintf("Performs controls and saves the given %s instance in the database", ep.getResourceClass()))
	ep.InGroup(group)

	return ep
}

// func GenericHandleRead[BOTYPE IBusinessObject](idProp IField, loadingType LoadingType) *oneForNoneEndpoint[BOTYPE] {
// 	ep := GetOne(
// 		// new (anonym) handler function here
// 		func(webCtx WebContext) (BOTYPE, hstatus.Code, string) {

// 			// boClass := GetClass[BOTYPE]()
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
// 			// boClass := GetClass[BOTYPE]()
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
