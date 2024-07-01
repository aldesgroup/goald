// ------------------------------------------------------------------------------------------------
// Here we implement some generic request handlers involved in CRUD
// N.B. quick & dirty implems for now
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"

	"github.com/aldesgroup/goald/features/hstatus"
)

func GenericHandleCreate[BOTYPE IBusinessObject]() *oneForOneEndpoint[BOTYPE, BOTYPE] {
	ep := PostOneGetOne[BOTYPE](
		// new (anonym) handler function here
		func(webCtx WebContext, input BOTYPE) (BOTYPE, hstatus.Code, string) {
			if errCreate := CreateBO(webCtx.GetBloContext(), input); errCreate != nil {
				return *new(BOTYPE), hstatus.InternalServerError,
					fmt.Sprintf("Failed creating a new '%T' instance: %s", input, errCreate)
			}

			return input, hstatus.Created, fmt.Sprintf("Created a new '%T' instance", input)
		},
		// passing the loading type
		"")

	return ep
}

func GenericHandleRead[BOTYPE IBusinessObject](idProp IField, loadingType LoadingType) *oneForNoneEndpoint[BOTYPE] {
	ep := GetOne[BOTYPE](
		// new (anonym) handler function here
		func(webCtx WebContext) (BOTYPE, hstatus.Code, string) {
			// boClass := GetClass[BOTYPE]()
			output, errRead := ReadBO(webCtx.GetBloContext(), idProp, webCtx.GetTargetRefOrID(), loadingType)
			if errRead != nil {
				return *new(BOTYPE), hstatus.InternalServerError,
					fmt.Sprintf("Failed reading '%s' instance '%s': %s", idProp.ownerClass().base().name, webCtx.GetTargetRefOrID(), errRead)
			}

			return output.(BOTYPE), hstatus.OK, fmt.Sprintf("Found the targeted '%T' instance", output)
		},
		// passing the loading type
		loadingType)

	ep.TargetWith(idProp)

	return ep
}

func GenericHandleUpdate[BOTYPE IBusinessObject](loadingType LoadingType) *oneForOneEndpoint[BOTYPE, BOTYPE] {
	ep := PutOne[BOTYPE](
		// new (anonym) handler function here
		func(webCtx WebContext, input BOTYPE) (BOTYPE, hstatus.Code, string) {
			if errUpdate := UpdateBO(webCtx.GetBloContext(), input, loadingType); errUpdate != nil {
				return *new(BOTYPE), hstatus.InternalServerError,
					fmt.Sprintf("Failed creating a new '%T' instance: %s", input, errUpdate)
			}

			return input, hstatus.OK, fmt.Sprintf("Update the given '%T' instance", input)
		},
		// passing the loading type
		loadingType)

	return ep
}

func GenericHandleDelete[BOTYPE IBusinessObject](idProp IField) *oneForNoneEndpoint[BOTYPE] {
	ep := DeleteOne[BOTYPE](
		// new (anonym) handler function here
		func(webCtx WebContext) (BOTYPE, hstatus.Code, string) {
			// boClass := GetClass[BOTYPE]()
			output, errRead := DeleteBO(webCtx.GetBloContext(), idProp, webCtx.GetTargetRefOrID())
			if errRead != nil {
				return *new(BOTYPE), hstatus.InternalServerError,
					fmt.Sprintf("Failed reading '%s' instance '%s': %s", idProp.ownerClass().base().name, webCtx.GetTargetRefOrID(), errRead)
			}

			return output.(BOTYPE), hstatus.OK, fmt.Sprintf("Deleted the targeted '%T' instance", output)
		})

	ep.TargetWith(idProp)

	return ep
}
