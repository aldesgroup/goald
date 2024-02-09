package goald

import (
	"fmt"
)

// ------------------------------------------------------------------------------------------------
// Panics
// ------------------------------------------------------------------------------------------------

func panicf(str string, params ...any) {
	panic(fmt.Sprintf(str, params...))
}

func panicErrf(err error, str string, params ...any) {
	if err != nil {
		panic(fmt.Sprintf(str, params...) + fmt.Sprintf("; cause: %s", err))
	}
}

// ------------------------------------------------------------------------------------------------
// Easier error wrapping
// ------------------------------------------------------------------------------------------------

func Error(msg string, params ...interface{}) error {
	return fmt.Errorf(msg, params...)
}

func ErrorC(cause error, msg string, params ...interface{}) error {
	return fmt.Errorf(fmt.Sprintf(msg, params...)+" --==|| Cause: %w", cause)
}
