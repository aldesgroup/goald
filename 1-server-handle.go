// ------------------------------------------------------------------------------------------------
// This is about handling HTTP requests
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/hstatus"
	r "github.com/julienschmidt/httprouter"
)

// ------------------------------------------------------------------------------------------------
// Serving the REST endpoints
// ------------------------------------------------------------------------------------------------

// TODO handle patching BOs with safeguards, like authorizing a limited list of fields (on the model for instance)

func (thisServer *server) ServeEndpoint(ep iEndpoint, w http.ResponseWriter, req *http.Request, params r.Params) {

	var reqCtx *httpRequestContext

	// Protecting against panics
	// recovering from an error happening in THIS routine while calling the operation;
	// THIS DOES NOT RECOVER what can happen in any sub-routine
	defer func() {
		// TODO limit this
		if err := recover(); err != nil {
			// unique ref
			errorReference := core.RandomString(8)

			// logging some details
			reqBody := ""
			if len(reqCtx.inputBodyBytes) > 0 {
				reqBody = " with body: " + string(reqCtx.inputBodyBytes)
			}
			thisServer.Error(false, fmt.Sprintf("Internal error n°%s = '%v', while calling '%s'%s. Stack: %s", errorReference, err, req.RequestURI, reqBody, string(debug.Stack())))

			// responding to the client
			reqCtx.write(&response{
				statusObj: hstatus.InternalServerError,
				Message:   fmt.Sprintf("Internal error n°%s", errorReference),
			}, w)
		}
	}()

	// TODO requestHandler pool
	// TODO defer : requestHandler release

	// TODO sync.Pool
	reqNum := thisServer.reqCount.Add(1)
	reqCtx = &httpRequestContext{
		server: thisServer,
		reqNum: reqNum,
		ILogger: thisServer.ILogger.WithPrefix(
			fmt.Sprintf("%s|%06d", thisServer.instance, reqNum),
		),
		start: time.Now(),
	}

	reqCtx.serve(ep, w, req, params)
}

// ------------------------------------------------------------------------------------------------
// Serving the REST endpoints
// ------------------------------------------------------------------------------------------------

// the type of response returned by all our REST endpoints
type response struct {
	Object     any          `json:"object,omitempty"`
	ObjectList any          `json:"objectList,omitempty"`
	statusObj  hstatus.Code `json:"-"`
	StatusCode int          `json:"statusCode"`
	Status     string       `json:"status"`
	Message    string       `json:"message"`
	Version    string       `json:"version"`
}

func errResp(status hstatus.Code, msg string, args ...any) *response {
	return &response{
		statusObj:  status,
		StatusCode: status.Val(),
		Status:     status.String(),
		Message:    fmt.Sprintf(msg, args...),
	}
}

// main HTTP SERVING functiont.De
func (thisReqCtx *httpRequestContext) serve(ep iEndpoint, w http.ResponseWriter, req *http.Request, params r.Params) {
	// logging
	thisReqCtx.Info(fmt.Sprintf("Serving %s (%s)", ep.getPathAsString(), ep.getLabel())) // TODO change

	// initialising the web context that's going to be passed to the applicative handler
	var targetRefOrID string
	if ep.getIDProp() != nil {
		targetRefOrID = params.ByName(ep.getIDProp().GetName())
	}

	// prepping the context that's going to contain all the input data
	// + some of the current endpoint's config
	webCtx := newWebContext(thisReqCtx, ep, targetRefOrID)

	// prepping the response
	resp := &response{Version: thisReqCtx.server.config.base().Version}

	// TODO check auth!

	// checking the input
	var input any
	if ep.hasBodyOrParamsInput() {
		var inputErr error
		if ep.isBodyInputRequired() {
			if input, inputErr = retrieveInputData(req, webCtx, ep); inputErr != nil {
				resp.statusObj = hstatus.BadRequest
				resp.Message = fmt.Sprintf("Bad input in request body (%s)", inputErr)

				goto End
			}
		} else {
			if input, inputErr = retrieveURLParams(req, webCtx, ep); inputErr != nil {
				resp.statusObj = hstatus.BadRequest
				resp.Message = fmt.Sprintf("Bad URL params (%s)", inputErr)

				goto End
			}
		}
	}

	// TODO do better - some "logging"
	// TODO only do this in verbose mode!
	if len(webCtx.inputBodyBytes) > 0 {
		if trimTo := ep.trimBodyLoggingTo(); trimTo > 0 && len(webCtx.inputBodyBytes) > trimTo {
			thisReqCtx.Debug(fmt.Sprintf("Body: %s [...]", string(webCtx.inputBodyBytes)[:trimTo]))
		} else {
			thisReqCtx.Debug(fmt.Sprintf("Body: %s", string(webCtx.inputBodyBytes)))
		}
	}

	// calling the endpoint's handler, which depends on its type
	if ep.hasBodyOrParamsInput() {
		if ep.isMultipleOutput() {
			if ep.isMultipleInput() {
				resp.ObjectList, resp.statusObj, resp.Message = ep.returnManyForMany(webCtx, input)
			} else {
				resp.ObjectList, resp.statusObj, resp.Message = ep.returnManyForOne(webCtx, input)
			}
		} else {
			if ep.isMultipleInput() {
				resp.Object, resp.statusObj, resp.Message = ep.returnOneForMany(webCtx, input)
			} else {
				resp.Object, resp.statusObj, resp.Message = ep.returnOneForOne(webCtx, input)
			}
		}
	} else {
		if ep.isMultipleOutput() {
			resp.ObjectList, resp.statusObj, resp.Message = ep.returnMany(webCtx)
		} else {
			resp.Object, resp.statusObj, resp.Message = ep.returnOne(webCtx)
		}
	}

End:
	// writing out the response
	thisReqCtx.write(resp, w)
}

// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

// writing out any response as JSON
func (thisReqCtx *httpRequestContext) write(resp *response, w http.ResponseWriter) {
	// setting headers must be the first thing done before writing anything
	// otherwise it is not taken into account.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// serialising the status
	resp.StatusCode = resp.statusObj.Val()
	resp.Status = resp.statusObj.String()

	// JSON-marshaling of the response
	jsonBytes, errMrsh := json.MarshalIndent(resp, "", "\t")
	if errMrsh != nil {
		resp = errResp(hstatus.InternalServerError, "Could not unmarshal the response: %s", errMrsh)
		jsonBytes, _ = json.MarshalIndent(resp, "", "\t")
	}

	// bit of logging
	if resp.statusObj.Val() > hstatus.BadRequest.Val() {
		thisReqCtx.Error(false, strconv.Itoa(resp.statusObj.Val())+": "+resp.Message+", in "+time.Since(thisReqCtx.start).String())
	} else {
		thisReqCtx.Info(strconv.Itoa(resp.statusObj.Val()) + ": " + resp.Message + ", in " + time.Since(thisReqCtx.start).String())
	}

	// writing the header before the body to avoid default HTTP code
	w.WriteHeader(resp.StatusCode)

	// actual writing out of the response
	if _, errWrite := w.Write(jsonBytes); errWrite != nil {
		// TODO change logging, make he request context a logger
		thisReqCtx.server.Error(true, fmt.Sprintf("Error while writing out the JSON response: %s", errWrite))
	}
}

// parsing the request's body to return the business object - or list of BOs - expected as input
func retrieveInputData(request *http.Request, webContext *webContextImpl, ep iEndpoint) (any, error) {
	// Handling unreadable body
	inputBodyBytes, readErr := io.ReadAll(request.Body)
	if readErr != nil {
		return nil, ErrorC(readErr, "Could not read request body!")
	}

	// Handling empty body
	if len(inputBodyBytes) == 0 {
		return nil, Error("Request body is empty")
	}

	// keeping track of the raw body
	webContext.inputBodyBytes = inputBodyBytes

	if ep.isMultipleInput() {
		// Handling array of bObj input: []*package.BObj
		bObjSlice := ep.getInputOrParamsModel().NewSlice()

		// Unmarshaling *[]*package.BObj as an interface - which is expected by the Unmarshal function
		if jsonErr := json.Unmarshal(inputBodyBytes, &bObjSlice); jsonErr != nil {
			return nil, ErrorC(jsonErr, "Could not unmarshall the JSON object array!")
		}

		// Not returning the reflect.Value, but the concrete instance associated with it
		return bObjSlice, nil

	} else {
		// Handling single bObj input: *package.BObj
		bObj := ep.getInputOrParamsModel().NewObject()

		if jsonErr := unmarshalBObj(inputBodyBytes, bObj); jsonErr != nil {
			return nil, ErrorC(jsonErr, "Could not unmarshall the JSON object!")
		}

		return bObj, nil
	}
}

// parsing the request's URL to build the expected URLQueryParams object
func retrieveURLParams(request *http.Request, _ *webContextImpl, ep iEndpoint) (any, error) {
	// new URLQueryParams object
	urlParams := ep.getInputOrParamsModel().NewObject().(IURLQueryParams)

	// transferring the URL param values from the URL to the object
	for _, field := range urlParams.getModel(urlParams).getFields() {
		valueToSet := request.URL.Query().Get(field.GetName())
		if valueToSet == "" {
			valueToSet = field.getDefaultValue()
		}
		urlParams.SetValueAsString(field.GetName(), valueToSet)
	}

	return urlParams, nil
}
