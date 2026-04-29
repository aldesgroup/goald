// ------------------------------------------------------------------------------------------------
// This is about handling HTTP requests
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/hstatus"
	r "github.com/julienschmidt/httprouter"
)

// ------------------------------------------------------------------------------------------------
// Serving the REST endpoints
// ------------------------------------------------------------------------------------------------

// TODO handle patching BOs with safeguards, like authorizing a limited list of fields (on the class for instance)

var reqCount int // to remove

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
			slog.Error(fmt.Sprintf("Internal error n°%s = '%v', while calling '%s'%s", errorReference, err, req.RequestURI, reqBody))

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
	reqCtx = &httpRequestContext{
		server: thisServer,
	}

	reqCtx.serve(ep, w, req, params)

}

// ------------------------------------------------------------------------------------------------
// Serving the REST endpoints
// ------------------------------------------------------------------------------------------------

// the type of response returned by all our REST endpoints
type response struct {
	Object     any          `json:"Object,omitempty"`
	ObjectList any          `json:"ObjectList,omitempty"`
	statusObj  hstatus.Code `json:"-"`
	StatusCode int          `json:"StatusCode"`
	Status     string       `json:"Status"`
	Message    string       `json:"Message"`
}

func errResp(_ int, _ string, _ ...any) *response {
	return &response{}
}

// main HTTP SERVING function
func (thisReqCtx *httpRequestContext) serve(ep iEndpoint, w http.ResponseWriter, req *http.Request, params r.Params) {

	// TODO remove
	reqCount++
	prefix := fmt.Sprintf("%06d|%s", reqCount, thisReqCtx.instance) //
	slog.Info(fmt.Sprintf("[%s] Serving %s (%s)", prefix, ep.getPathAsString(), ep.getLabel()))

	// initialising the web context that's going to be passed to the applicative handler
	var targetRefOrID string
	if ep.getIDProp() != nil {
		targetRefOrID = params.ByName(ep.getIDProp().getName())
	}

	// prepping the context that's going to contain all the input data
	// + some of the current endpoint's config
	webCtx := newWebContext(thisReqCtx, ep, targetRefOrID)

	// prepping the response
	resp := &response{}

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
			slog.Debug(fmt.Sprintf("Body: %s [...]", string(webCtx.inputBodyBytes)[:trimTo]))
		} else {
			slog.Debug(fmt.Sprintf("Body: %s", string(webCtx.inputBodyBytes)))
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
		resp = errResp(http.StatusInternalServerError, "Could not unmarshal the response: %s", errMrsh)
		jsonBytes, _ = json.MarshalIndent(resp, "", "\t")
	}

	// writing the header before the body to avoid default HTTP code
	w.WriteHeader(resp.StatusCode)

	// actual writing out of the response
	if _, errWrite := w.Write(jsonBytes); errWrite != nil {
		// TODO change logging
		slog.Error(fmt.Sprintf("Error while writing out the JSON response: %s", errWrite))
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
		bObjClass := classRegistry.items[ep.getInputOrParamsClass()]
		if bObjClass == nil {
			return nil, Error("No '%s' class has been registered!", ep.getInputOrParamsClass())
		}
		bObjSlice := bObjClass.NewSlice()

		// Unmarshaling *[]*package.BObj as an interface - which is expected by the Unmarshal function
		if jsonErr := json.Unmarshal(inputBodyBytes, &bObjSlice); jsonErr != nil {
			return nil, ErrorC(jsonErr, "Could not unmarshall the JSON object array!")
		}

		// Not returning the reflect.Value, but the concrete instance associated with it
		return bObjSlice, nil

	} else {
		// Handling single bObj input: *package.BObj
		bObjClass := classRegistry.items[ep.getInputOrParamsClass()]
		if bObjClass == nil {
			return nil, Error("No '%s' class has been registered!", ep.getInputOrParamsClass())
		}
		bObj := bObjClass.NewObject()

		if jsonErr := json.Unmarshal(inputBodyBytes, bObj); jsonErr != nil {
			return nil, ErrorC(jsonErr, "Could not unmarshall the JSON object!")
		}

		return bObj, nil
	}
}

// parsing the request's URL to build the expected URLQueryParams object
func retrieveURLParams(request *http.Request, _ *webContextImpl, ep iEndpoint) (any, error) {
	// getting the right class utils
	classUtils := classRegistry.items[ep.getInputOrParamsClass()]

	// new URLQueryParams object
	urlParams := classUtils.NewObject().(IURLQueryParams)

	// transferring the URL param values from the URL to the object
	for _, field := range modelForName(ep.getInputOrParamsClass()).base().fields {
		valueToSet := request.URL.Query().Get(field.getName())
		if valueToSet == "" {
			valueToSet = field.getDefaultValue()
		}
		classUtils.SetValueAsString(urlParams, field.getName(), valueToSet)
	}

	return urlParams, nil
}
