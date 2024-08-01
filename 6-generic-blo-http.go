// ------------------------------------------------------------------------------------------------
// Here is provided a general way of doing HTTP request to external data providers
// ------------------------------------------------------------------------------------------------
package goald

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aldesgroup/goald/features/hstatus"
)

// interface & general methods
type IExternalHttpRequest interface {
	WithData(dataObj any) IExternalHttpRequest
	WithTimeout(timeout time.Duration) IExternalHttpRequest
	WithBasicAuth(user, pass string) IExternalHttpRequest
	WithHeader(key, value string) IExternalHttpRequest
}

func newHttpReq[ResponseType any](method string, url string, responseObj *ResponseType) *externalHttpRequest[ResponseType] {
	return &externalHttpRequest[ResponseType]{method: method, url: url, responseObj: responseObj, headers: map[string]string{}}
}

func HttpPost[ResponseType any](url string, responseObj *ResponseType) *externalHttpRequest[ResponseType] {
	return newHttpReq(http.MethodPost, url, responseObj)
}

func HttpGet[ResponseType any](url string, responseObj *ResponseType) *externalHttpRequest[ResponseType] {
	return newHttpReq(http.MethodGet, url, responseObj)
}

func HttpPut[ResponseType any](url string, responseObj *ResponseType) *externalHttpRequest[ResponseType] {
	return newHttpReq(http.MethodPut, url, responseObj)
}

// implementation
type externalHttpRequest[ResponseType any] struct {
	method      string
	url         string
	dataObj     any
	timeout     time.Duration
	user        string
	pass        string
	headers     map[string]string
	responseObj *ResponseType
}

func (thisCtx *externalHttpRequest[ResponseType]) WithData(dataObj any) *externalHttpRequest[ResponseType] {
	thisCtx.dataObj = dataObj
	return thisCtx
}

func (thisCtx *externalHttpRequest[ResponseType]) WithTimeout(timeout time.Duration) *externalHttpRequest[ResponseType] {
	thisCtx.timeout = timeout
	return thisCtx
}

func (thisCtx *externalHttpRequest[ResponseType]) WithBasicAuth(user, pass string) *externalHttpRequest[ResponseType] {
	thisCtx.user = user
	thisCtx.pass = pass
	return thisCtx
}

func (thisCtx *externalHttpRequest[ResponseType]) WithHeader(key, value string) *externalHttpRequest[ResponseType] {
	thisCtx.headers[key] = value
	return thisCtx
}

// main external HTTP request execution method
func (thisCtx *externalHttpRequest[ResponseType]) Exec(failOnBadStatusCode bool) (*ResponseType, hstatus.Code, error) {
	// the object that's maybe being sent in the request body
	var dataBuffer *bytes.Buffer
	if thisCtx.dataObj != nil {
		dataBytes, errMarsh := json.Marshal(thisCtx.dataObj)
		if errMarsh != nil {
			return nil, hstatus.InternalServerError, ErrorC(errMarsh, "Could not marshall the Aldes cloud request")
		}
		dataBuffer = bytes.NewBuffer(dataBytes)

		// a bit of logging
		// TODO do better
		if true {
			prettyJson, _ := json.MarshalIndent(thisCtx.dataObj, "", "	")
			slog.Debug(fmt.Sprintf("Sending this data: %s", string(prettyJson)))
		}
	}

	// initialising an HTTP request to send to the service
	httpRequest, errReq := http.NewRequest(thisCtx.method, thisCtx.url, dataBuffer)
	if errReq != nil {
		return nil, hstatus.InternalServerError, ErrorC(errReq, "issue while initialising a request")
	}

	// initialising the connection
	client := &http.Client{Timeout: thisCtx.timeout}

	// adding basic authentication if required
	if thisCtx.user != "" {
		httpRequest.SetBasicAuth(thisCtx.user, thisCtx.pass)
	}

	// adding headers
	for key, value := range thisCtx.headers {
		httpRequest.Header.Add(key, value)
	}

	// processing the request by calling the remote URL using our client (and timing it)
	slog.Debug(fmt.Sprintf("HTTP call: %s %s", thisCtx.method, thisCtx.url))
	resp, errResponse := client.Do(httpRequest)
	if errResponse != nil {
		return nil, hstatus.InternalServerError, ErrorC(errResponse, "Error while HTTP calling")
	}

	// very important to prevent leaks !
	// cf. https://husobee.github.io/golang/memory/leak/2016/02/11/go-mem-leak.html
	// has to be done right before the read, since before that the body is nil
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			slog.Error(fmt.Sprintf("error while closing the response body: %s", errClose.Error()))
		}
	}()

	// reading the response body
	respBody, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil, hstatus.InternalServerError, ErrorC(errRead, "Error while reading the response body")
	}

	// TODO better logging
	if true {
		slog.Debug(fmt.Sprintf("Got this data: %s", string(respBody)))
	}

	// reading the response status - should we fail on a bad status code, we don't even need to read the body
	status := resp.StatusCode
	if resp.StatusCode >= 400 && failOnBadStatusCode {
		slog.Error(fmt.Sprintf("Response body: %s", string(respBody)))
		return nil, hstatus.For(status), Error("External service (%s) responded with a %d status code", thisCtx.url, status)
	}

	// unmarshalling the response
	if errJSON := json.Unmarshal(respBody, thisCtx.responseObj); errJSON != nil {
		return nil, hstatus.InternalServerError, ErrorC(errJSON, "Could not unmarshal response body")
	}

	// returning the response as a concrete object
	return thisCtx.responseObj, hstatus.For(status), nil
}
