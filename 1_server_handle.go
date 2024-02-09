// ------------------------------------------------------------------------------------------------
// This is about handling HTTP requests
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/aldesgroup/goald/hstatus"
	r "github.com/julienschmidt/httprouter"
)

// ------------------------------------------------------------------------------------------------
// Serving the REST endpoints
// ------------------------------------------------------------------------------------------------

// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
// Search Params as JSON objects => maybe even some BOs ? Or another whole class of objects that can be described automatically

// TODO handle patching BOs with safeguards, like authorizing a limited list of fields (on the class for instance)

var reqCount int // to remove

func (thisServer *server) ServeEndpoint(ep iEndpoint, w http.ResponseWriter, req *http.Request, params r.Params) {
	// TODO requestHandler pool
	reqCtx := &httpRequestContext{
		server: thisServer,
	}

	reqCtx.serve(ep, w, req, params)

	// TODO requestHandler release
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

func errResp(status int, str string, params ...any) *response {
	return &response{}
}

// main HTTP SERVING function
func (thisReqCtx *httpRequestContext) serve(ep iEndpoint, w http.ResponseWriter, req *http.Request, params r.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// TODO remove
	reqCount++
	prefix := fmt.Sprintf("%06d", reqCount)
	log.Printf("[%s] Serving %s %s: %s", prefix, ep.getMethod(), ep.getFullPath(), ep.getLabel())

	// initialising the web context that's going to be passed to the applicative handler
	var targetRefOrID string
	if ep.getIDProp() != nil {
		targetRefOrID = params.ByName(ep.getIDProp().getName())
	}

	// prepping the context that's going to contain all the input data
	// + some of the current endpoint's config
	webCtx := newWebContext(thisReqCtx, ep.getLoadingType(), targetRefOrID)

	// prepping the response
	resp := &response{}

	// TODO check auth!

	// checking the input body
	var bodyInput any
	if ep.isRequiredInput() {
		var inputErr error
		if bodyInput, inputErr = retrieveInputData(req, webCtx, ep); inputErr != nil {
			resp.statusObj = hstatus.BadRequest
			resp.Message = fmt.Sprintf("Bad input in request body (%s)", inputErr)

			goto End
		}
	}

	// TODO do better - some "logging"
	log.Printf("Body: %s", string(webCtx.inputBody))

	// calling the endpoint's handler, which depends on its type
	if ep.isRequiredInput() {
		if ep.isMultipleOutput() {
			if ep.isMultipleInput() {
				resp.ObjectList, resp.statusObj, resp.Message = ep.returnManyForMany(webCtx, bodyInput)
			} else {
				resp.ObjectList, resp.statusObj, resp.Message = ep.returnManyForOne(webCtx, bodyInput)
			}
		} else {
			if ep.isMultipleInput() {
				resp.Object, resp.statusObj, resp.Message = ep.returnOneForMany(webCtx, bodyInput)
			} else {
				resp.Object, resp.statusObj, resp.Message = ep.returnOneForOne(webCtx, bodyInput)
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
		log.Printf("Error while writing out the JSON response: %s", errWrite)
	}
}

// parsing the request's body to
func retrieveInputData(request *http.Request, webContext *webContextImpl, ep iEndpoint) (any, error) {
	// Handling unreadable body
	inputBody, readErr := io.ReadAll(request.Body)
	if readErr != nil {
		return nil, ErrorC(readErr, "Could not read request body!")
	}

	// Handling empty body
	if len(inputBody) == 0 {
		return nil, Error("Request body is empty")
	}

	// keeping track of the raw body
	webContext.inputBody = inputBody

	if ep.isMultipleInput() {
		// Handling array of bObj input: []*package.BObj
		// explanation:      *           []              *     package.BObj         - removes the starting *
		bObjSlice := reflect.New(reflect.SliceOf(reflect.PtrTo(ep.getInputType()))).Elem()

		// Unmarshaling *[]*package.BObj as an interface - which is expected by the Unmarshal function
		if jsonErr := json.Unmarshal(inputBody, bObjSlice.Addr().Interface()); jsonErr != nil {
			return nil, ErrorC(jsonErr, "Could not unmarshall the JSON object array!")
		}

		// Not returning the reflect.Value, but the concrete instance associated with it
		return bObjSlice.Interface(), nil

	} else {
		// Handling single bObj input: *package.BObj
		// explanation: *   package.BObj      - needed by the unmarshaling
		bObj := reflect.New(ep.getInputType()).Interface()

		if jsonErr := json.Unmarshal(inputBody, bObj); jsonErr != nil {
			return nil, ErrorC(jsonErr, "Could not unmarshall the JSON object!")
		}

		return bObj, nil
	}
}
