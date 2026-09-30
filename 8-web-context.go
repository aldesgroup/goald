package goald

import (
	"fmt"

	"github.com/aldesgroup/goald/features/auth"
)

// ------------------------------------------------------------------------------------------------
// WebContext provides the necessary info to applicatively handle incoming HTTP requests
// ------------------------------------------------------------------------------------------------

type WebContext interface {
	restContext
	GetBloContext() BloContext
	GetResourceRefOrID() string
}

// default implementation for web context
type webContextImpl struct {
	*httpRequestContext // wrapping one of the server's children handling 1 request
	ep                  iEndpoint
	resourceModel       IBusinessObjectModel
	resourceRefOrID     string // the ID or ref, or whatever property value used to clearly identify a resource
	bloContext          BloContext
	authClaims          *auth.Claims // the caller's claims, set once authenticateRequest() succeeds; nil for public endpoints or when auth isn't configured
	currentUser         IUser        // the caller's IUser, lazily resolved from authClaims via the registered UserResolverFunc
	currentUserResolved bool         // whether currentUser has already been (attempted to be) resolved, so we only try once per request
}

// type check
var _ WebContext = (*webContextImpl)(nil)

// constructors
func newWebContext(reqCtx *httpRequestContext, ep iEndpoint, targetRefOrID string) *webContextImpl {
	return &webContextImpl{
		httpRequestContext: reqCtx,
		ep:                 ep,
		resourceModel:      ep.getResourceModel(),
		resourceRefOrID:    targetRefOrID,
	}
}

func (thisWebCtx *webContextImpl) GetBloContext() BloContext {
	if thisWebCtx.bloContext == nil {
		thisWebCtx.bloContext = newHttpBloContextFromWebCtx(thisWebCtx)
	}

	return thisWebCtx.bloContext
}

func (thisWebCtx *webContextImpl) getResourceModel() IBusinessObjectModel {
	if thisWebCtx.resourceModel == nil {
		thisWebCtx.resourceModel = thisWebCtx.ep.getResourceModel()
	}

	return thisWebCtx.resourceModel
}

func (thisWebCtx *webContextImpl) GetResourceRefOrID() string {
	return thisWebCtx.resourceRefOrID
}

// GetCurrentUser implements [BloContext] & [restContext], lazily resolving the caller's IUser (if
// any) from its authenticated claims, via the application-registered UserResolverFunc.
func (thisWebCtx *webContextImpl) GetCurrentUser() IUser {
	if thisWebCtx.authClaims == nil {
		return nil
	}

	if !thisWebCtx.currentUserResolved {
		thisWebCtx.currentUserResolved = true

		resolver := getUserResolver()
		if resolver == nil {
			thisWebCtx.Warn("An authenticated request came in, but no user resolver is registered (see RegisterUserResolver)")
			return nil
		}

		user, err := resolver(thisWebCtx.GetBloContext(), thisWebCtx.authClaims)
		if err != nil {
			thisWebCtx.Error(false, fmt.Sprintf("Could not resolve the current user: %s", err))
			return nil
		}

		thisWebCtx.currentUser = user
	}

	return thisWebCtx.currentUser
}
