// ------------------------------------------------------------------------------------------------
// This package is used to type the HTTP status codes, and allowing using them,
// rather than writing bare values like 200, 401, etc.
// ------------------------------------------------------------------------------------------------
package hstatus

import "net/http"

type Code struct{ val int }

func wrap(val int) Code {
	return Code{val}
}

func (c Code) Val() int {
	return c.val
}

func (c Code) String() string {
	return http.StatusText(c.val)
}

var (
	Continue           = wrap(http.StatusContinue)
	SwitchingProtocols = wrap(http.StatusSwitchingProtocols)
	Processing         = wrap(http.StatusProcessing)
	EarlyHints         = wrap(http.StatusEarlyHints)

	OK                   = wrap(http.StatusOK)
	Created              = wrap(http.StatusCreated)
	Accepted             = wrap(http.StatusAccepted)
	NonAuthoritativeInfo = wrap(http.StatusNonAuthoritativeInfo)
	NoContent            = wrap(http.StatusNoContent)
	ResetContent         = wrap(http.StatusResetContent)
	PartialContent       = wrap(http.StatusPartialContent)
	MultiStatus          = wrap(http.StatusMultiStatus)
	AlreadyReported      = wrap(http.StatusAlreadyReported)
	IMUsed               = wrap(http.StatusIMUsed)

	MultipleChoices   = wrap(http.StatusMultipleChoices)
	MovedPermanently  = wrap(http.StatusMovedPermanently)
	Found             = wrap(http.StatusFound)
	SeeOther          = wrap(http.StatusSeeOther)
	NotModified       = wrap(http.StatusNotModified)
	UseProxy          = wrap(http.StatusUseProxy)
	TemporaryRedirect = wrap(http.StatusTemporaryRedirect)
	PermanentRedirect = wrap(http.StatusPermanentRedirect)

	BadRequest                   = wrap(http.StatusBadRequest)
	Unauthorized                 = wrap(http.StatusUnauthorized)
	PaymentRequired              = wrap(http.StatusPaymentRequired)
	Forbidden                    = wrap(http.StatusForbidden)
	NotFound                     = wrap(http.StatusNotFound)
	MethodNotAllowed             = wrap(http.StatusMethodNotAllowed)
	NotAcceptable                = wrap(http.StatusNotAcceptable)
	ProxyAuthRequired            = wrap(http.StatusProxyAuthRequired)
	RequestTimeout               = wrap(http.StatusRequestTimeout)
	Conflict                     = wrap(http.StatusConflict)
	Gone                         = wrap(http.StatusGone)
	LengthRequired               = wrap(http.StatusLengthRequired)
	PreconditionFailed           = wrap(http.StatusPreconditionFailed)
	RequestEntityTooLarge        = wrap(http.StatusRequestEntityTooLarge)
	RequestURITooLong            = wrap(http.StatusRequestURITooLong)
	UnsupportedMediaType         = wrap(http.StatusUnsupportedMediaType)
	RequestedRangeNotSatisfiable = wrap(http.StatusRequestedRangeNotSatisfiable)
	ExpectationFailed            = wrap(http.StatusExpectationFailed)
	Teapot                       = wrap(http.StatusTeapot)
	MisdirectedRequest           = wrap(http.StatusMisdirectedRequest)
	UnprocessableEntity          = wrap(http.StatusUnprocessableEntity)
	Locked                       = wrap(http.StatusLocked)
	FailedDependency             = wrap(http.StatusFailedDependency)
	TooEarly                     = wrap(http.StatusTooEarly)
	UpgradeRequired              = wrap(http.StatusUpgradeRequired)
	PreconditionRequired         = wrap(http.StatusPreconditionRequired)
	TooManyRequests              = wrap(http.StatusTooManyRequests)
	RequestHeaderFieldsTooLarge  = wrap(http.StatusRequestHeaderFieldsTooLarge)
	UnavailableForLegalReasons   = wrap(http.StatusUnavailableForLegalReasons)

	InternalServerError           = wrap(http.StatusInternalServerError)
	NotImplemented                = wrap(http.StatusNotImplemented)
	BadGateway                    = wrap(http.StatusBadGateway)
	ServiceUnavailable            = wrap(http.StatusServiceUnavailable)
	GatewayTimeout                = wrap(http.StatusGatewayTimeout)
	HTTPVersionNotSupported       = wrap(http.StatusHTTPVersionNotSupported)
	VariantAlsoNegotiates         = wrap(http.StatusVariantAlsoNegotiates)
	InsufficientStorage           = wrap(http.StatusInsufficientStorage)
	LoopDetected                  = wrap(http.StatusLoopDetected)
	NotExtended                   = wrap(http.StatusNotExtended)
	NetworkAuthenticationRequired = wrap(http.StatusNetworkAuthenticationRequired)
)
