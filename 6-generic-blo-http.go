// ------------------------------------------------------------------------------------------------
// Here is provided a general way of doing HTTP request to external data providers
// ------------------------------------------------------------------------------------------------
package goald

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aldesgroup/goald/features/hstatus"
	"github.com/aldesgroup/goald/features/logging"
)

// interface & general methods
type IOutgoingHttpRequest interface {
	WithData(dataObj any) IOutgoingHttpRequest
	WithTimeout(timeout time.Duration) IOutgoingHttpRequest
	WithBasicAuth(user, pass string) IOutgoingHttpRequest
	WithHeader(key, value string) IOutgoingHttpRequest
}

func newHttpReq(method string, url string, failOnBadStatusCode bool) *outgoingHttpRequest {
	return &outgoingHttpRequest{method: method, url: url, headers: map[string]string{}, failOnBadStatusCode: failOnBadStatusCode}
}

func HttpPost(url string, failOnBadStatusCode bool) *outgoingHttpRequest {
	return newHttpReq(http.MethodPost, url, failOnBadStatusCode)
}

func HttpGet(url string, failOnBadStatusCode bool) *outgoingHttpRequest {
	return newHttpReq(http.MethodGet, url, failOnBadStatusCode)
}

func HttpPut(url string, failOnBadStatusCode bool) *outgoingHttpRequest {
	return newHttpReq(http.MethodPut, url, failOnBadStatusCode)
}

// implementation
type outgoingHttpRequest struct {
	method              string
	url                 string
	dataObj             any
	timeout             time.Duration
	user                string
	pass                string
	headers             map[string]string
	failOnBadStatusCode bool
}

func (thisReq *outgoingHttpRequest) WithData(dataObj any) *outgoingHttpRequest {
	thisReq.dataObj = dataObj
	return thisReq
}

func (thisReq *outgoingHttpRequest) WithTimeout(timeout time.Duration) *outgoingHttpRequest {
	thisReq.timeout = timeout
	return thisReq
}

func (thisReq *outgoingHttpRequest) WithBasicAuth(user, pass string) *outgoingHttpRequest {
	thisReq.user = user
	thisReq.pass = pass
	return thisReq
}

func (thisReq *outgoingHttpRequest) WithHeader(key, value string) *outgoingHttpRequest {
	thisReq.headers[key] = value
	return thisReq
}

// Runs a prepared Rest HTTP request to an external service, and fills the given object with the response body
func HttpGetResponseAs[ResponseType any](logger logging.ILogger, outReq *outgoingHttpRequest, responseObj *ResponseType) (*ResponseType, hstatus.Code, error) {
	// the object that's maybe being sent in the request body
	var dataBuffer *bytes.Buffer
	if outReq.dataObj != nil {
		dataBytes, errMarsh := json.Marshal(outReq.dataObj)
		if errMarsh != nil {
			return nil, hstatus.InternalServerError, ErrorC(errMarsh, "Could not marshall the Aldes cloud request")
		}
		dataBuffer = bytes.NewBuffer(dataBytes)

		// a bit of logging
		// TODO do better
		if true {
			prettyJson, _ := json.MarshalIndent(outReq.dataObj, "", "	")
			logger.Debug(fmt.Sprintf("Sending this data: %s", string(prettyJson)))
		}
	}

	// initialising an HTTP request to send to the service
	httpRequest, errReq := http.NewRequest(outReq.method, outReq.url, dataBuffer)
	if errReq != nil {
		return nil, hstatus.InternalServerError, ErrorC(errReq, "issue while initialising a request")
	}

	// initialising the connection
	client := &http.Client{Timeout: outReq.timeout}

	// adding basic authentication if required
	if outReq.user != "" {
		httpRequest.SetBasicAuth(outReq.user, outReq.pass)
	}

	// adding headers
	for key, value := range outReq.headers {
		httpRequest.Header.Add(key, value)
	}

	// processing the request by calling the remote URL using our client (and timing it)
	logger.Debug(fmt.Sprintf("HTTP call: %s %s", outReq.method, outReq.url))
	resp, errResponse := client.Do(httpRequest)
	if errResponse != nil {
		return nil, hstatus.InternalServerError, ErrorC(errResponse, "Error while HTTP calling")
	}

	// very important to prevent leaks !
	// cf. https://husobee.github.io/golang/memory/leak/2016/02/11/go-mem-leak.html
	// has to be done right before the read, since before that the body is nil
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.Error(false, fmt.Sprintf("error while closing the response body: %s", errClose.Error()))
		}
	}()

	// reading the response body
	respBody, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil, hstatus.InternalServerError, ErrorC(errRead, "Error while reading the response body")
	}

	// TODO better logging
	if true {
		logger.Debug(fmt.Sprintf("Got this data: %s", string(respBody)))
	}

	// reading the response status - should we fail on a bad status code, we don't even need to read the body
	status := resp.StatusCode
	if resp.StatusCode >= 400 && outReq.failOnBadStatusCode {
		return nil, hstatus.For(status), Error("External service (%s) responded with a %d status code", outReq.url, status)
	}

	// unmarshalling the response
	if errJSON := json.Unmarshal(respBody, responseObj); errJSON != nil {
		return nil, hstatus.InternalServerError, ErrorC(errJSON, "Could not unmarshal response body")
	}

	// returning the response as a concrete object
	return responseObj, hstatus.For(status), nil
}

// Runs a prepared Rest HTTP request to an external service, and retrieves the Goald business object from the response
func HttpGetBObject[ResponseType IBusinessObject](logger logging.ILogger, outReq *outgoingHttpRequest, responseObj ResponseType) (ResponseType, hstatus.Code, error) {
	resp, status, err := HttpGetResponseAs(logger, outReq, &response{Object: responseObj})
	if err != nil || resp == nil || resp.Object == nil {
		return responseObj, status, err
	}
	return resp.Object.(ResponseType), status, err
}

// Runs a prepared Rest HTTP request to an external service, and retrieves the list of Goald business objects from the response
func HttpGetBObjList[ResponseType IBusinessObject](logger logging.ILogger, outReq *outgoingHttpRequest, responseList []ResponseType) ([]ResponseType, hstatus.Code, error) {
	resp, status, err := HttpGetResponseAs(logger, outReq, &response{ObjectList: responseList})
	if err != nil || resp == nil || resp.ObjectList == nil {
		return responseList, status, err
	}
	return resp.ObjectList.([]ResponseType), status, err
}
