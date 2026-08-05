// ------------------------------------------------------------------------------------------------
// The code here is about how we describe endpoints
// ------------------------------------------------------------------------------------------------
package goald

import (
	"net/http"
	"strings"

	"github.com/aldesgroup/goald/features/hstatus"
	"github.com/aldesgroup/goald/features/reflection"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Endpoint groups
// ------------------------------------------------------------------------------------------------

type EndpointGroup struct {
	Name        string
	Description string
}

var allEndpointGroups = make(map[string]*EndpointGroup)

func NewEndpointGroup(name string, description string) *EndpointGroup {
	// a group cannot already exist with the same name
	if _, alreadyExists := allEndpointGroups[name]; alreadyExists {
		panic("an endpoint group already exists with this name: " + name)
	}

	allEndpointGroups[name] = &EndpointGroup{
		Name:        name,
		Description: description,
	}

	return allEndpointGroups[name]
}

// ------------------------------------------------------------------------------------------------
// Base structs & generic methods
// ------------------------------------------------------------------------------------------------

type iEndpoint interface {
	getMethod() string
	getResourceModel() IBusinessObjectModel
	getIDProp() IField
	getOperationPath(includeIdOrRefName bool) string // e.g. /rest/document/reduce/:id if includeIdOrRefName = true, /rest/document/reduce otherwise
	getPathAsString() string                         // e.g. GET /rest/document/reduce/:id
	getLabel() string
	getDescription() string
	getLoadingType() LoadingType
	isMultipleOutput() bool
	hasBodyOrParamsInput() bool
	isBodyInputRequired() bool
	isMultipleInput() bool
	getInputOrParamsModel() IBusinessObjectModel
	isCalledFromWebApp() bool
	isCalledFromNativeApp() bool
	trimBodyLoggingTo() int
	getGroup() *EndpointGroup
	getOutputListLen(list any) int // the number of elements in a returned object list, whatever its concrete []ResourceType is
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
	method                 string               // get, post, put...
	resourceModelName      utils.ModelName      // the name of the objects reached through this endpoint
	resourceModel          IBusinessObjectModel // the model of the objects reached through this endpoint
	basePath               string               // the endpoint's base path, which is the lower-cased resource type name
	actionPath             string               // do we need an additional path for a non-CRUD action, like "reduce" in: "GET /document/reduce/:id"
	idProp                 IField               // if a specific BO is targeted, this has to be through one of its properties
	label                  string               // short label to describe the endpoint
	description            string               // a bit longer text to describe the endpoint
	multipleOutput         bool                 // if true, then the endpoint delivers arrays of BOs, rather than a single one
	loadingType            LoadingType          // how the returned resource(s) are loaded
	bodyInputRequired      bool                 // if true, then we expect something in the request body
	multipleInput          bool                 // if true, then we expect an array of BOs in the body, rather than a single one
	inputOrParamsModelName utils.ModelName      // if bodyInputRequired = true, then this is the type of the input
	inputOrParamsModel     IBusinessObjectModel // if bodyInputRequired = true, then this is the type of the input
	calledFromWebApp       bool                 // if true then this endpoint can be called from the webapp, so the BOs involved might be synced through codegen
	calledFromNativeApp    bool                 // if true then this endpoint can be called from the native app, so the BOs involved might be synced through codegen
	logBodyFirstCharsNb    int                  // if > 0, we're logging only the n-th first chars of the body, not it's entirety
	group                  *EndpointGroup       // if not empty, this is the name of the group this endpoint belongs to, which can be used for documentation or other purposes
}

func (ep *endpoint[ResourceType]) getMethod() string {
	return ep.method
}

// getOutputListLen returns the number of elements held by a returned object list. Since
// isMultipleOutput() handlers actually return a concrete []ResourceType (not a []any), this
// generic method - resolved at compile time for each concrete ResourceType, no runtime
// reflection needed - is the only safe way to read that slice's length back in server code that
// only deals with iEndpoint and "any" values.
func (ep *endpoint[ResourceType]) getOutputListLen(list any) int {
	if list == nil {
		return 0
	}
	return len(list.([]ResourceType))
}

func (ep *endpoint[ResourceType]) getResourceModel() IBusinessObjectModel {
	if ep.resourceModelName != "" && ep.resourceModel == nil {
		ep.resourceModel = modelFor(ep.resourceModelName, true)
	}
	return ep.resourceModel
}

func (ep *endpoint[ResourceType]) getIDProp() IField {
	return ep.idProp
}

func (ep *endpoint[ResourceType]) getOperationPath(includeIdOrRefName bool) string {
	operationPath := apiPath + "/" + ep.basePath
	if ep.actionPath != "" {
		if ep.actionPath[0:1] == "/" {
			operationPath += ep.actionPath
		} else {
			operationPath += "/" + ep.actionPath
		}
	}
	if ep.idProp != nil && includeIdOrRefName {
		operationPath += "/:" + ep.idProp.GetName()
	}

	return operationPath
}

func (ep *endpoint[ResourceType]) getPathAsString() string {
	return ep.getMethod() + " " + ep.getOperationPath(true)
}

func (ep *endpoint[ResourceType]) getLabel() string {
	return ep.label
}

func (ep *endpoint[ResourceType]) getDescription() string {
	return ep.description
}

func (ep *endpoint[ResourceType]) getLoadingType() LoadingType {
	return ep.loadingType
}

func (ep *endpoint[ResourceType]) isMultipleOutput() bool {
	return ep.multipleOutput
}

func (ep *endpoint[ResourceType]) hasBodyOrParamsInput() bool {
	return ep.inputOrParamsModelName != ""
}

func (ep *endpoint[ResourceType]) isBodyInputRequired() bool {
	return ep.bodyInputRequired
}

func (ep *endpoint[ResourceType]) isMultipleInput() bool {
	return ep.multipleInput
}

func (ep *endpoint[ResourceType]) getInputOrParamsModel() IBusinessObjectModel {
	if ep.inputOrParamsModelName != "" && ep.inputOrParamsModel == nil {
		ep.inputOrParamsModel = modelFor(ep.inputOrParamsModelName, true)
	}
	return ep.inputOrParamsModel
}

func (ep *endpoint[ResourceType]) isCalledFromWebApp() bool {
	return ep.calledFromWebApp
}

func (ep *endpoint[ResourceType]) isCalledFromNativeApp() bool {
	return ep.calledFromNativeApp
}

func (ep *endpoint[ResourceType]) trimBodyLoggingTo() int {
	return ep.logBodyFirstCharsNb
}

func (ep *endpoint[ResourceType]) getGroup() *EndpointGroup {
	return ep.group
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

	resourceModelName := utils.ModelName(reflection.TypeNameOf((*new(ResourceType)), true))
	var inputOrParamsModelName utils.ModelName
	if bodyInputRequired || withURLParams {
		inputOrParamsModelName = utils.ModelName(reflection.TypeNameOf((*new(InputOrParamsType)), true))
	}

	return &endpoint[ResourceType]{
		method:                 method,
		resourceModelName:      resourceModelName,
		basePath:               strings.ToLower(string(resourceModelName)),
		multipleOutput:         multipleOutput,
		loadingType:            loadingType,
		bodyInputRequired:      bodyInputRequired,
		multipleInput:          multipleInput,
		inputOrParamsModelName: inputOrParamsModelName,
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

// Providing a longer description for this endpoint
func (thisEndpoint *endpoint[ResourceType]) Description(description string) *endpoint[ResourceType] {
	thisEndpoint.description = description

	return thisEndpoint
}

// Indicating that this endpoint can be called from the associated web app (through Aldev),
// so that I/O code can be automatically generated within it
func (thisEndpoint *endpoint[ResourceType]) SetCalledFromWebApp() *endpoint[ResourceType] {
	thisEndpoint.calledFromWebApp = true

	// taking the opportunity to enrich the model that are impacted here
	modelFor(thisEndpoint.resourceModelName).setUsedInWebApp()
	if thisEndpoint.inputOrParamsModelName != "" {
		modelFor(thisEndpoint.inputOrParamsModelName).setUsedInWebApp()
	}

	return thisEndpoint
}

// Indicating that this endpoint can be called from the associated native app (through Aldev),
// so that I/O code can be automatically generated within it
func (thisEndpoint *endpoint[ResourceType]) SetCalledFromNativeApp() *endpoint[ResourceType] {
	thisEndpoint.calledFromNativeApp = true

	// taking the opportunity to enrich the model that are impacted here
	modelFor(thisEndpoint.resourceModelName).setUsedInNativeApp()
	if thisEndpoint.inputOrParamsModelName != "" {
		modelFor(thisEndpoint.inputOrParamsModelName).setUsedInNativeApp()
	}

	return thisEndpoint
}

// Avoid too many logs
func (thisEndpoint *endpoint[ResourceType]) TrimBodyLogging(trimTo int) *endpoint[ResourceType] {
	thisEndpoint.logBodyFirstCharsNb = trimTo

	return thisEndpoint
}

// Putting the endpoint in a group, which can be used for documentation or other purposes
func (thisEndpoint *endpoint[ResourceType]) InGroup(group *EndpointGroup) *endpoint[ResourceType] {
	thisEndpoint.group = group

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

	return handleOne(http.MethodGet, handlerFunc, loadingType)
}

// Declaring an endpoint to delete 1 BO instance with a DELETE request
func DeleteOne[ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext) (ResourceType, hstatus.Code, string),
) *oneForNoneEndpoint[ResourceType] {

	return handleOne(http.MethodDelete, handlerFunc, "")
}

// // Declaring an endpoint to return N BO instances from a GET request
// func GetMany[ResourceType IBusinessObject](
// 	handlerFunc func(webCtx WebContext) ([]ResourceType, hstatus.Code, string),
// 	loadingType LoadingType,
// ) *manyForNoneEndpoint[ResourceType] {

// 	return handleMany(http.MethodGet, handlerFunc, loadingType)
// }

// Declaring an endpoint to return 1 BO instance from 1 POSTed BO instance
func PostOneGetOne[InputType, ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext, input InputType) (ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *oneForOneEndpoint[InputType, ResourceType] {

	return handleOneForOne(http.MethodPost, handlerFunc, loadingType)
}

// Declaring an endpoint to return 1 BO instance from 1 PUT BO instance
func PutOne[InputType, ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext, input ResourceType) (ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *oneForOneEndpoint[ResourceType, ResourceType] {

	return handleOneForOne(http.MethodPut, handlerFunc, loadingType)
}

// Declaring an endpoint to return N BO instance from N POSTed BO instances
func PostManyGetMany[InputType, ResourceType IBusinessObject](
	handlerFunc func(webCtx WebContext, input []InputType) ([]ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *manyForManyEndpoint[InputType, ResourceType] {

	return handleManyForMany(http.MethodPost, handlerFunc, loadingType)
}

// Declaring an endpoint to return N BO instance from query parameters that are described with 1 SearchParamValues
func GetManyWithParams[ResourceType IBusinessObject, SearchParamValuesType ISearchParamValues](
	handlerFunc func(webCtx WebContext, SearchParamValues SearchParamValuesType) ([]ResourceType, hstatus.Code, string),
	loadingType LoadingType,
) *manyForOneEndpoint[SearchParamValuesType, ResourceType] {

	return handleManyForOne(http.MethodGet, handlerFunc, loadingType, false)
}
