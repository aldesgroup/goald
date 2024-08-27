// ------------------------------------------------------------------------------------------------
// The code here is about how we describe endpoints
// ------------------------------------------------------------------------------------------------
package goald

import (
	"net/http"
	"strings"

	"github.com/aldesgroup/goald/features/hstatus"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Base structs & generic methods
// ------------------------------------------------------------------------------------------------
type iEndpoint interface {
	getMethod() string
	getResourceClass() className
	getIDProp() IField
	getFullPath() string
	getLabel() string
	getLoadingType() LoadingType
	isMultipleOutput() bool
	hasBodyOrParamsInput() bool
	isBodyInputRequired() bool
	isMultipleInput() bool
	getInputOrParamsClass() className
	isFromWebApp() bool
	returnOne(webCtx WebContext) (any, hstatus.Code, string)
	returnMany(webCtx WebContext) (any, hstatus.Code, string)
	returnOneForOne(webCtx WebContext, input any) (any, hstatus.Code, string)
	returnOneForMany(webCtx WebContext, inputs any) (any, hstatus.Code, string)
	returnManyForOne(webCtx WebContext, input any) (any, hstatus.Code, string)
	returnManyForMany(webCtx WebContext, inputs any) (any, hstatus.Code, string)
}

// an endpoint object is parametrized by the potential objects of type I,
// and the output objects of type O, i.e. the resource type
type endpoint[ResourceType IBusinessObject] struct {
	method             string      // get, post, put...
	resourceClass      className   // the class of the objects reached through this endpoint
	basePath           string      // the endpoint's base path, which is the lower-cased resource type name
	actionPath         string      // do we need an additional path for a non-CRUD action, like "reduce" in: "GET /document/reduce/:id"
	idProp             IField      // if a specific BO is targeted, this has to be through one of its properties
	fullPath           string      // resulting from the parameter type, action path and id property
	label              string      // short label to describe the endpoint
	multipleOutput     bool        // if true, then the endpoint delivers arrays of BOs, rather than a single one
	loadingType        LoadingType // how the returned resource(s) are loaded
	bodyInputRequired  bool        // if true, then we expect something in the request body
	multipleInput      bool        // if true, then we expect an array of BOs in the body, rather than a single one
	inputOrParamsClass className   // if bodyInputRequired = true, then this is the type of the input
	fromWebApp         bool        // if true then this endpoint can be called from the webapp, so the BOs involved might be synced through codegen
}

func (ep *endpoint[ResourceType]) getMethod() string {
	return ep.method
}

func (ep *endpoint[ResourceType]) getResourceClass() className {
	return ep.resourceClass
}

func (ep *endpoint[ResourceType]) getIDProp() IField {
	return ep.idProp
}

func (ep *endpoint[ResourceType]) getFullPath() string {
	if ep.fullPath == "" {
		ep.fullPath = "/" + ep.basePath
		if ep.actionPath != "" {
			if ep.actionPath[0:1] == "/" {
				ep.fullPath += ep.actionPath
			} else {
				ep.fullPath += "/" + ep.actionPath
			}
		}
		if ep.idProp != nil {
			ep.fullPath += "/:" + ep.idProp.getName()
		}
	}

	return ep.fullPath
}

func (ep *endpoint[ResourceType]) getLabel() string {
	return ep.label
}

func (ep *endpoint[ResourceType]) getLoadingType() LoadingType {
	return ep.loadingType
}

func (ep *endpoint[ResourceType]) isMultipleOutput() bool {
	return ep.multipleOutput
}

func (ep *endpoint[ResourceType]) hasBodyOrParamsInput() bool {
	return ep.inputOrParamsClass != ""
}

func (ep *endpoint[ResourceType]) isBodyInputRequired() bool {
	return ep.bodyInputRequired
}

func (ep *endpoint[ResourceType]) isMultipleInput() bool {
	return ep.multipleInput
}

func (ep *endpoint[ResourceType]) getInputOrParamsClass() className {
	return ep.inputOrParamsClass
}

func (ep *endpoint[ResourceType]) isFromWebApp() bool {
	return ep.fromWebApp
}

func (ep *endpoint[ResourceType]) returnOne(webCtx WebContext) (any, hstatus.Code, string) {
	panic("no generic implementation here")
}

func (ep *endpoint[ResourceType]) returnMany(webCtx WebContext) (any, hstatus.Code, string) {
	panic("no generic implementation here")
}

func (ep *endpoint[ResourceType]) returnOneForOne(webCtx WebContext, input any) (any, hstatus.Code, string) {
	panic("no generic implementation here")
}

func (ep *endpoint[ResourceType]) returnOneForMany(webCtx WebContext, inputs any) (any, hstatus.Code, string) {
	panic("no generic implementation here")
}

func (ep *endpoint[ResourceType]) returnManyForOne(webCtx WebContext, input any) (any, hstatus.Code, string) {
	panic("no generic implementation here")
}

func (ep *endpoint[ResourceType]) returnManyForMany(webCtx WebContext, inputs any) (any, hstatus.Code, string) {
	panic("no generic implementation here")
}

// ------------------------------------------------------------------------------------------------
// Endpoint declaration & building
// ------------------------------------------------------------------------------------------------

// common part of any endpoint
func newEndpoint[InputOrParamsType, ResourceType IBusinessObject](
	multipleOutput bool,
	method string,
	loadingType LoadingType,
	bodyInputRequired bool,
	multipleInput bool,
	withURLParams bool,
) *endpoint[ResourceType] {

	resourceClsName := className(utils.TypeNameOf((*new(ResourceType)), true))
	var inputOrParamsClsName className
	if bodyInputRequired || withURLParams {
		inputOrParamsClsName = className(utils.TypeNameOf((*new(InputOrParamsType)), true))
	}

	return &endpoint[ResourceType]{
		method:             method,
		resourceClass:      resourceClsName,
		basePath:           strings.ToLower(string(resourceClsName)),
		multipleOutput:     multipleOutput,
		loadingType:        loadingType,
		bodyInputRequired:  bodyInputRequired,
		multipleInput:      multipleInput,
		inputOrParamsClass: inputOrParamsClsName,
	}
}

// Allows to add an additional path to the base path: /basepath[/additionalpath][/:targetedPropValue]
func (thisEndpoint *endpoint[ResourceType]) At(actionPath string) *endpoint[ResourceType] {
	thisEndpoint.actionPath = actionPath

	return thisEndpoint
}

// Allows to add a property value to the base path: /basepath[/additionalpath][/:targetedPropValue]
func (thisEndpoint *endpoint[ResourceType]) TargetWith(idProp IField) *endpoint[ResourceType] {
	thisEndpoint.idProp = idProp

	return thisEndpoint
}

// Providing a short description for this endpoint
func (thisEndpoint *endpoint[ResourceType]) Label(label string) *endpoint[ResourceType] {
	thisEndpoint.label = label

	return thisEndpoint
}

// Indicating that this endpoint can be called from the associated web app (through Aldev),
// so that I/O code can be automatically generated within it
func (thisEndpoint *endpoint[ResourceType]) FromWebApp() *endpoint[ResourceType] {
	thisEndpoint.fromWebApp = true

	return thisEndpoint
}

// ------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------

// Declaring an endpoint to return 1 BO instance from a GET request
func GetOne[ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext) (ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *oneForNoneEndpoint[ResourceType] {

	return handleOne[ResourceType](http.MethodGet, handlerFunc, loadingType)
}

// Declaring an endpoint to delete 1 BO instance with a DELETE request
func DeleteOne[ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext) (ResourceType, hstatus.Code, string),
) *oneForNoneEndpoint[ResourceType] {

	return handleOne[ResourceType](http.MethodDelete, handlerFunc, "")
}

// Declaring an endpoint to return N BO instances from a GET request
func GetMany[ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext) ([]ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *manyForNoneEndpoint[ResourceType] {

	return handleMany[ResourceType](http.MethodGet, handlerFunc, loadingType)
}

// Declaring an endpoint to return 1 BO instance from 1 POSTed BO instance
func PostOneGetOne[InputType, ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext, input InputType) (ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *oneForOneEndpoint[InputType, ResourceType] {

	return handleOneForOne[InputType, ResourceType](http.MethodPost, handlerFunc, loadingType)
}

// Declaring an endpoint to return 1 BO instance from 1 PUT BO instance
func PutOne[InputType, ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext, input ResourceType) (ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *oneForOneEndpoint[ResourceType, ResourceType] {

	return handleOneForOne[ResourceType, ResourceType](http.MethodPut, handlerFunc, loadingType)
}

// Declaring an endpoint to return N BO instance from N POSTed BO instances
func PostManyGetMany[InputType, ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext, input []InputType) ([]ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *manyForManyEndpoint[InputType, ResourceType] {

	return handleManyForMany[InputType, ResourceType](http.MethodPost, handlerFunc, loadingType)
}

// Declaring an endpoint to return N BO instance from query parameters that are described with 1 URLQueryParams
func GetManyWithParams[ResourceType IBusinessObject, QueryParamsType IURLQueryParams](
	handlerFunc func(webCtx WebContext, queryParams QueryParamsType) ([]ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *manyForOneEndpoint[QueryParamsType, ResourceType] {

	return handleManyForOne[QueryParamsType, ResourceType](http.MethodGet, handlerFunc, loadingType, false)
}
