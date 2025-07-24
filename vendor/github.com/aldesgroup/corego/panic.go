package core

import (
	"fmt"
)

// ------------------------------------------------------------------------------------------------
// Panics
// ------------------------------------------------------------------------------------------------

// PanicMsg panics with the given message
func PanicMsg(str string, params ...any) {
	panic(fmt.Sprintf(str, params...))
}

// PanicMsgIf panics with the given message if the condition is true
func PanicMsgIf(cond bool, str string, params ...any) {
	if cond {
		panic(fmt.Sprintf(str, params...))
	}
}

// PanicMsgIfErr panics with the given message if the error is not nil
func PanicMsgIfErr(err error, str string, params ...any) {
	if err != nil {
		panic(fmt.Sprintf(str, params...) + fmt.Sprintf("; cause: %s", err))
	}
}

// PanicIfErr panics if the error is not nil
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
