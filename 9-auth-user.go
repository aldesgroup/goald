package goald

type IUser interface {
	IBusinessObject
	GetUsername() string
}
