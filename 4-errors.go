package goald

import (
	"fmt"
	"log/slog"
	"runtime/debug"
)

// ------------------------------------------------------------------------------------------------
// Easier error wrapping
// ------------------------------------------------------------------------------------------------

func Error(msg string, params ...interface{}) error {
	return fmt.Errorf(msg, params...)
}

func ErrorC(cause error, msg string, params ...interface{}) error {
	return fmt.Errorf(fmt.Sprintf(msg, params...)+" --==|| Cause: %w", cause)
}

func RecoverError(msg string, params ...interface{}) {
	if err := recover(); err != nil {
		// TODO change
		slog.Error(fmt.Sprintf(msg+". Cause: %v. Stack: \n%s", append(params, err, string(debug.Stack()))...))
	}
}
