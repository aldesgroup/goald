// ------------------------------------------------------------------------------------------------
// The code here contains functions that help declaring endpoints
// ------------------------------------------------------------------------------------------------

package goald

import "github.com/aldesgroup/goald/features/hstatus"

// ------------------------------------------------------------------------------------------------
// The different endpoint types: (1) = 1 resource for no input
// ------------------------------------------------------------------------------------------------

type oneForNoneEndpoint[ResourceType IBusinessObject] struct {
	*endpoint[ResourceType]
	handlerFunc func(webCtx WebContext) (ResourceType, hstatus.Code, string)
}

func handleOne[ResourceType IBusinessObject](
	method string,
	handlerFunc func(webCtx WebContext) (ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *oneForNoneEndpoint[ResourceType] {

	return registerEndpoint(&oneForNoneEndpoint[ResourceType]{
		endpoint: newEndpoint[ResourceType, ResourceType](
			false,
			method,
			loadingType,
			false,
			false,
			false),
		handlerFunc: handlerFunc,
	}).(*oneForNoneEndpoint[ResourceType])
}

func (ep *oneForNoneEndpoint[ResourceType]) returnOne(webCtx WebContext) (any, hstatus.Code, string) {
	return ep.handlerFunc(webCtx)
}

// ------------------------------------------------------------------------------------------------
// The different endpoint types: (2) = N resources for no input
// ------------------------------------------------------------------------------------------------

type manyForNoneEndpoint[ResourceType IBusinessObject] struct {
	*endpoint[ResourceType]
	handlerFunc func(webCtx WebContext) ([]ResourceType, hstatus.Code, string)
}

func handleMany[ResourceType IBusinessObject](
	method string,
	handlerFunc func(webCtx WebContext) ([]ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *manyForNoneEndpoint[ResourceType] {

	return registerEndpoint(&manyForNoneEndpoint[ResourceType]{
		endpoint: newEndpoint[ResourceType, ResourceType](
			true,
			method,
			loadingType,
			false,
			false,
			false),
		handlerFunc: handlerFunc,
	}).(*manyForNoneEndpoint[ResourceType])
}

// adapting the parametrized function to a generic format that the main *httpRequestContext.serve() can call
func (ep *manyForNoneEndpoint[ResourceType]) returnMany(webCtx WebContext) (any,
	hstatus.Code, string) {
	return ep.handlerFunc(webCtx)
}

// ------------------------------------------------------------------------------------------------
// The different endpoint types: (3) = 1 resource for 1 input
// ------------------------------------------------------------------------------------------------

type oneForOneEndpoint[InputType, ResourceType IBusinessObject] struct {
	*endpoint[ResourceType]
	handlerFunc func(webCtx WebContext, input InputType) (ResourceType, hstatus.Code, string)
}

func handleOneForOne[InputType, ResourceType IBusinessObject](
	method string,
	handlerFunc func(webCtx WebContext, input InputType) (ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *oneForOneEndpoint[InputType, ResourceType] {

	return registerEndpoint(&oneForOneEndpoint[InputType, ResourceType]{
		endpoint: newEndpoint[InputType, ResourceType](
			false,
			method,
			loadingType,
			true,
			false,
			false),
		handlerFunc: handlerFunc,
	}).(*oneForOneEndpoint[InputType, ResourceType])
}

// adapting the parametrized function to a generic format that the main *httpRequestContext.serve() can call
func (ep *oneForOneEndpoint[InputType, ResourceType]) returnOneForOne(webCtx WebContext, input any) (any, hstatus.Code, string) {
	return ep.handlerFunc(webCtx, input.(InputType))
}

// // ------------------------------------------------------------------------------------------------
// // The different endpoint types: (4) = 1 resource for N inputs
// // ------------------------------------------------------------------------------------------------

// type oneForManyEndpoint[InputType, ResourceType IBusinessObject] struct {
// 	*endpoint[ResourceType]
// 	handlerFunc func(webCtx WebContext, input []InputType) (ResourceType, hstatus.Code, string)
// }

// func handleOneForMany[InputType, ResourceType IBusinessObject](
// 	method string,
// 	handlerFunc func(webCtx WebContext, input []InputType) (ResourceType, hstatus.Code, string),
// 	loadingType LoadingType,
// ) *oneForManyEndpoint[InputType, ResourceType] {

// 	return registerEndpoint(&oneForManyEndpoint[InputType, ResourceType]{
// 		endpoint:    newEndpoint[ResourceType](false, method, loadingType, true, true),
// 		handlerFunc: handlerFunc,
// 	}).(*oneForManyEndpoint[InputType, ResourceType])
// }

// // adapting the parametrized function to a generic format that the main *httpRequestContext.serve() can call
// func (ep *oneForManyEndpoint[InputType, ResourceType]) returnOneForMany(webCtx WebContext, inputs any) (any, hstatus.Code, string) {
// 	return ep.handlerFunc(webCtx, inputs.([]InputType))
// }

// ------------------------------------------------------------------------------------------------
// The different endpoint types: (5) = N resources for 1 input
// ------------------------------------------------------------------------------------------------

type manyForOneEndpoint[InputOrParamsType, ResourceType IBusinessObject] struct {
	*endpoint[ResourceType]
	handlerFunc func(webCtx WebContext, input InputOrParamsType) ([]ResourceType, hstatus.Code, string)
}

func handleManyForOne[InputOrParamsType, ResourceType IBusinessObject](
	method string,
	handlerFunc func(webCtx WebContext, input InputOrParamsType) ([]ResourceType, hstatus.Code, string),
	loadingType LoadingType,
	withBodyInput bool,
) *manyForOneEndpoint[InputOrParamsType, ResourceType] {

	return registerEndpoint(&manyForOneEndpoint[InputOrParamsType, ResourceType]{
		endpoint: newEndpoint[InputOrParamsType, ResourceType](
			true,
			method,
			loadingType,
			withBodyInput,
			false,
			!withBodyInput),
		handlerFunc: handlerFunc,
	}).(*manyForOneEndpoint[InputOrParamsType, ResourceType])
}

// adapting the parametrized function to a generic format that the main *httpRequestContext.serve() can call
func (ep *manyForOneEndpoint[InputOrParamsType, ResourceType]) returnManyForOne(webCtx WebContext, input any) (any, hstatus.Code, string) {
	return ep.handlerFunc(webCtx, input.(InputOrParamsType))
}

// ------------------------------------------------------------------------------------------------
// The different endpoint types: (6) = N resources for N inputs
// ------------------------------------------------------------------------------------------------

type manyForManyEndpoint[InputType, ResourceType IBusinessObject] struct {
	*endpoint[ResourceType]
	handlerFunc func(webCtx WebContext, input []InputType) ([]ResourceType, hstatus.Code, string)
}

func handleManyForMany[InputType, ResourceType IBusinessObject](
	method string,
	handlerFunc func(webCtx WebContext, input []InputType) ([]ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *manyForManyEndpoint[InputType, ResourceType] {

	return registerEndpoint(&manyForManyEndpoint[InputType, ResourceType]{
		endpoint: newEndpoint[InputType, ResourceType](
			true,
			method,
			loadingType,
			true,
			true,
			false),
		handlerFunc: handlerFunc,
	}).(*manyForManyEndpoint[InputType, ResourceType])
}

// adapting the parametrized function to a generic format that the main *httpRequestContext.serve() can call
func (ep *manyForManyEndpoint[InputType, ResourceType]) returnManyForMany(webCtx WebContext, inputs any) (any, hstatus.Code, string) {
	return ep.handlerFunc(webCtx, inputs.([]InputType))
}

// ------------------------------------------------------------------------------------------------
// Querying for BOs through URLs
// ------------------------------------------------------------------------------------------------

// particular business object class
type IURLQueryParamsClass interface {
	IBusinessObjectClass
}

// particular business object
type IURLQueryParams interface {
	IBusinessObject
}

// particular business object class implem
type urlQueryParamsClass struct {
	businessObjectClass
}

// particular business object implem
type URLQueryParams struct {
	BusinessObject
}

func NewURLQueryParamsClass() IURLQueryParamsClass {
	class := &urlQueryParamsClass{
		businessObjectClass: businessObjectClass{
			fields: map[string]IField{},
			inNoDB: true,
		},
	}

	return class
}
