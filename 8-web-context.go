package goald

// ------------------------------------------------------------------------------------------------
// WebContext provides the necessary info to applicatively handle incoming HTTP requests
// ------------------------------------------------------------------------------------------------

type WebContext interface {
	restContext
	GetBloContext() BloContext
	GetResourceRefOrID() string
	GetResourceLoadingType() LoadingType // returns the loading type of the current main resources (BOs) being worked on
}

// default implementation for web context
type webContextImpl struct {
	*httpRequestContext // wrapping one of the server's children handling 1 request
	ep                  iEndpoint
	resourceModel       IBusinessObjectModel
	resourceRefOrID     string // the ID or ref, or whatever property value used to clearly identify a resource
	bloContext          BloContext
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

func (thisWebCtx *webContextImpl) GetResourceLoadingType() LoadingType {
	return thisWebCtx.ep.getLoadingType()
}
