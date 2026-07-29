package goald

// ------------------------------------------------------------------------------------------------
// WebContext provides the necessary info to applicatively handle incoming HTTP requests
// ------------------------------------------------------------------------------------------------

type WebContext interface {
	restContext
	GetBloContext() BloContext
	GetTargetRefOrID() string
	GetResource() IBusinessObjectModel   // the class of the resource being requested
	GetResourceLoadingType() LoadingType // returns the loading type of the current main resources (BOs) being worked on
}

// default implementation for web context
type webContextImpl struct {
	*httpRequestContext // wrapping one of the server's children handling 1 request
	ep                  iEndpoint
	resource            IBusinessObjectModel
	bloContext          BloContext
	targetClass         IClass // the class of the resource being requested
	targetRefOrID       string // the ID or ref, or whatever property value used to clearly identify a resource
}

// type check
var _ WebContext = (*webContextImpl)(nil)

// constructors
func newWebContext(reqCtx *httpRequestContext, ep iEndpoint, targetRefOrID string) *webContextImpl {
	return &webContextImpl{
		httpRequestContext: reqCtx,
		ep:                 ep,
		targetClass:        ep.getResourceClass(),
		targetRefOrID:      targetRefOrID,
	}
}

func (thisWebCtx *webContextImpl) GetBloContext() BloContext {
	if thisWebCtx.bloContext == nil {
		thisWebCtx.bloContext = newHttpBloContextFromWebCtx(thisWebCtx)
	}

	return thisWebCtx.bloContext
}

func (thisWebCtx *webContextImpl) GetResource() IBusinessObjectModel {
	if thisWebCtx.resource == nil {
		thisWebCtx.resource = thisWebCtx.ep.getResourceClass().getModel()
	}

	return thisWebCtx.resource
}

func (thisWebCtx *webContextImpl) getTargetResourceClass() IClass {
	return thisWebCtx.targetClass
}

func (thisWebCtx *webContextImpl) GetTargetRefOrID() string {
	return thisWebCtx.targetRefOrID
}

func (thisWebCtx *webContextImpl) GetResourceLoadingType() LoadingType {
	return thisWebCtx.ep.getLoadingType()
}
