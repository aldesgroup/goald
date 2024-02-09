// ------------------------------------------------------------------------------------------------
// The code here is about how we describe endpoints
// ------------------------------------------------------------------------------------------------
package goald

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aldesgroup/goald/hstatus"
)

// ------------------------------------------------------------------------------------------------
// Base structs & generic methods
// ------------------------------------------------------------------------------------------------
type iEndpoint interface {
	getMethod() string
	getIDProp() IField
	getFullPath() string
	getLabel() string
	getLoadingType() LoadingType
	isMultipleOutput() bool
	isRequiredInput() bool
	isMultipleInput() bool
	getInputType() reflect.Type
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
	method         string       // get, post, put...
	basePath       string       // the endpoint's base path, which is the lower-cased resource type name
	actionPath     string       // do we need an additional path for a non-CRUD action, like "reduce" in: "GET /document/reduce/:id"
	idProp         IField       // if a specific BO is targeted, this has to be through one of its properties
	fullPath       string       // resulting from the parameter type, action path and id property
	label          string       // short label to describe the endpoint
	multipleOutput bool         // if true, then the endpoint delivers arrays of BOs, rather than a single one
	loadingType    LoadingType  // how the returned resource(s) are loaded
	inputExpected  bool         // if true, then we expect something in the request body
	multipleInput  bool         // if true, then we expect an array of BOs in the body, rather than a single one
	inputType      reflect.Type // if inputExpected = true, then this is the type of the input
}

func (ep *endpoint[ResourceType]) getMethod() string {
	return ep.method
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

func (ep *endpoint[ResourceType]) isRequiredInput() bool {
	return ep.inputExpected
}

func (ep *endpoint[ResourceType]) isMultipleInput() bool {
	return ep.multipleInput
}

func (ep *endpoint[ResourceType]) getInputType() reflect.Type {
	return ep.inputType
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
func newEndpoint[InputType, ResourceType IBusinessObject](
	multipleOutput bool,
	method string,
	loadingType LoadingType,
	inputExpected bool,
	multipleInput bool,
) *endpoint[ResourceType] {

	var inputType reflect.Type
	if inputExpected {
		// if the InputType is '*Project', then we're keeping  'Project' here (without the pointer)
		inputType = reflect.TypeOf(*new(InputType)).Elem()
	}

	return &endpoint[ResourceType]{
		method:         method,
		basePath:       strings.ToLower(reflect.TypeOf(*new(ResourceType)).Elem().Name()),
		multipleOutput: multipleOutput,
		loadingType:    loadingType,
		inputExpected:  inputExpected,
		multipleInput:  multipleInput,
		inputType:      inputType,
	}
}

// Allows to add an additional path to the base path: /basepath[/ad/di/tio/nal/path][/:targetedPropValue]
func (thisEndpoint *endpoint[ResourceType]) At(actionPath string) *endpoint[ResourceType] {
	thisEndpoint.actionPath = actionPath

	return thisEndpoint
}

// Allows to add a property value to the path: /basepath[/ad/di/tio/nal/path][/:targetedPropValue]
func (thisEndpoint *endpoint[ResourceType]) TargetWith(idProp IField) *endpoint[ResourceType] {
	thisEndpoint.idProp = idProp

	return thisEndpoint
}

// Providing a short description for this endpoint
func (thisEndpoint *endpoint[ResourceType]) Label(label string) *endpoint[ResourceType] {
	thisEndpoint.label = label

	return thisEndpoint
}

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
		endpoint:    newEndpoint[ResourceType, ResourceType](false, method, loadingType, false, false),
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
		endpoint:    newEndpoint[ResourceType, ResourceType](true, method, loadingType, false, false),
		handlerFunc: handlerFunc,
	}).(*manyForNoneEndpoint[ResourceType])
}

// adapting the parametrized function to a generic format that the main *httpRequestContext.serve() can call
func (ep *manyForNoneEndpoint[ResourceType]) returnMany(webCtx WebContext) (any, hstatus.Code, string) {
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
		endpoint:    newEndpoint[InputType, ResourceType](false, method, loadingType, true, false),
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

// // ------------------------------------------------------------------------------------------------
// // The different endpoint types: (5) = N resources for 1 input
// // ------------------------------------------------------------------------------------------------

// type manyForOneEndpoint[InputType, ResourceType IBusinessObject] struct {
// 	*endpoint[ResourceType]
// 	handlerFunc func(webCtx WebContext, input InputType) ([]ResourceType, hstatus.Code, string)
// }

// func handleManyForOne[InputType, ResourceType IBusinessObject](
// 	method string,
// 	handlerFunc func(webCtx WebContext, input InputType) ([]ResourceType, hstatus.Code, string),
// 	loadingType LoadingType,
// ) *manyForOneEndpoint[InputType, ResourceType] {

// 	return registerEndpoint(&manyForOneEndpoint[InputType, ResourceType]{
// 		endpoint:    newEndpoint[ResourceType](true, method, loadingType, true, false),
// 		handlerFunc: handlerFunc,
// 	}).(*manyForOneEndpoint[InputType, ResourceType])
// }

// // adapting the parametrized function to a generic format that the main *httpRequestContext.serve() can call
// func (ep *manyForOneEndpoint[InputType, ResourceType]) returnManyForOne(webCtx WebContext, input any) (any, hstatus.Code, string) {
// 	return ep.handlerFunc(webCtx, input.(InputType))
// }

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
		endpoint:    newEndpoint[InputType, ResourceType](true, method, loadingType, true, true),
		handlerFunc: handlerFunc,
	}).(*manyForManyEndpoint[InputType, ResourceType])
}

// adapting the parametrized function to a generic format that the main *httpRequestContext.serve() can call
func (ep *manyForManyEndpoint[InputType, ResourceType]) returnManyForMany(webCtx WebContext, inputs any) (any, hstatus.Code, string) {
	return ep.handlerFunc(webCtx, inputs.([]InputType))
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
