package goald

// ------------------------------------------------------------------------------------------------
// DaoContext should contain the necessary info for handling database access
// ------------------------------------------------------------------------------------------------
type DaoContext_ interface {
	AppContext // ability to log and get access to the app config

	ExecInsertQuery(bObj IBusinessObject) (int64, error) // executes the insert query for a given business object and returns its new ID
}

// ------------------------------------------------------------------------------------------------
// Implementation
// ------------------------------------------------------------------------------------------------

type daoContextImpl struct {
	AppContext
}

// // type check
// var _ DaoContext = (*daoContextImpl)(nil)

// // constructors
// func newDaoContextFromHttpBloCtx(httpBloCtx *httpBloContextImpl, bObj IBusinessObject) DaoContext {
// 	return &daoContextImpl{
// 		AppContext:         httpBloCtx.restContext,
// 		iBusinessObjectDAO: getDaoFor(bObj.getClassName()),
// 		iDbProvider:        bObj,
// 	}
// }

// this won't work dammit...
