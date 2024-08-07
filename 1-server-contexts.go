// ------------------------------------------------------------------------------------------------
// Here are defined the contexts passed in the different layers of our multi-tier architecture:
// - (no context) for the business object layer (*__.go files)
// - BloContext for the Business LOgic code (used in *__blo.go files)
// - DaoContext for the Data Access Objects (used in *__dao.go files)
// - WebContext for the web endpoints code (used in *__web.go files)
// ------------------------------------------------------------------------------------------------
package goald

// ------------------------------------------------------------------------------------------------
// AppContext contains the minimal info set that should be accessible in all the layers of the app
// ------------------------------------------------------------------------------------------------
type AppContext interface {
	CustomConfig() ICustomConfig // returns the app's custom part of the config
}

type appContextImpl struct {
}

// ------------------------------------------------------------------------------------------------
// iRestContext is used in the context of handling with a REST resource (single or plural)
type iRestContext interface {
	AppContext
}

// ------------------------------------------------------------------------------------------------
// BloContext is a context that should provide the necessary info for Business LOgic code
type BloContext interface {
	iRestContext
	GetDaoContext() DaoContext
}

type baseBloContextImpl struct {
	*appContextImpl // common implem of AppContext
}

// default implementation for business logic context
type bloContextImpl struct {
	*baseBloContextImpl
	*httpRequestContext // wrapping one of the server's children handling 1 request
	daoContext          DaoContext
}

func (thisBloCtx *bloContextImpl) GetDaoContext() DaoContext {
	return thisBloCtx.daoContext // TODO instantiate
}

// the server is its own DAo context
func (thisServer *server) GetDaoContext() DaoContext {
	return thisServer
}

// ------------------------------------------------------------------------------------------------
// ServerContext is a particular Business Logic Context used at app startup
// Implemented by the `server` struct
type ServerContext interface {
	BloContext
	Start()
}

// ------------------------------------------------------------------------------------------------
// DaoContext should contain the necessary info for handling database access
type DaoContext interface {
	iRestContext
}

// ------------------------------------------------------------------------------------------------
// WebContext provides the necessary info to applicatively handle incoming HTTP requests
type WebContext interface {
	iRestContext
	GetBloContext() BloContext
	GetTargetRefOrID() string
	GetResourceClass() IBusinessObjectClass // the class of the resource being requested
	GetResourceLoadingType() LoadingType    // returns the loading type of the current main resources (BOs) being worked on
}

// default implementation for web context
type webContextImpl struct {
	*appContextImpl     // common implem of AppContext
	*httpRequestContext // wrapping one of the server's children handling 1 request
	ep                  iEndpoint
	resourceClass       IBusinessObjectClass
	targetRefOrID       string // the ID or ref, or whatever property value used to clearly identify a resource
	inputBodyBytes      []byte // keeping track of the incoming request body
	bloContext          BloContext
}

// type check
var _ WebContext = (*webContextImpl)(nil)

func newWebContext(reqCtx *httpRequestContext, ep iEndpoint, targetRefOrID string) *webContextImpl {
	return &webContextImpl{
		appContextImpl:     &appContextImpl{},
		httpRequestContext: reqCtx,
		ep:                 ep,
		targetRefOrID:      targetRefOrID,
	}
}

func (thisWebCtx *webContextImpl) GetBloContext() BloContext {
	// initialising it when first needed
	if thisWebCtx.bloContext == nil {
		thisWebCtx.bloContext = &bloContextImpl{
			httpRequestContext: thisWebCtx.httpRequestContext,
			baseBloContextImpl: &baseBloContextImpl{
				appContextImpl: thisWebCtx.appContextImpl,
			},
		}
	}

	return thisWebCtx.bloContext
}

func (thisWebCtx *webContextImpl) GetResourceClass() IBusinessObjectClass {
	if thisWebCtx.resourceClass == nil {
		thisWebCtx.resourceClass = classRegistry.classes[thisWebCtx.ep.getInputOrParamsClass()]
	}

	return thisWebCtx.resourceClass
}

func (thisWebCtx *webContextImpl) GetTargetRefOrID() string {
	return thisWebCtx.targetRefOrID
}

func (thisWebCtx *webContextImpl) GetResourceLoadingType() LoadingType {
	return thisWebCtx.ep.getLoadingType()
}
