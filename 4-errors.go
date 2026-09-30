package goald

import (
	"fmt"
	"runtime/debug"

	"github.com/aldesgroup/goald/features/logging"
)

// ------------------------------------------------------------------------------------------------
// Easier error wrapping
// ------------------------------------------------------------------------------------------------

func Error(msg string, params ...any) error {
	return fmt.Errorf(msg, params...)
}

func ErrorC(cause error, msg string, params ...any) error {
	return fmt.Errorf(fmt.Sprintf(msg, params...)+" --==|| Cause: %w", cause)
}

func RecoverError(logger logging.ILogger, msg string, params ...any) {
	if err := recover(); err != nil {
		// TODO change
		logger.Error(false, fmt.Sprintf(msg+". Cause: %v. Stack: \n%s", append(params, err, string(debug.Stack()))...))
	}
}
