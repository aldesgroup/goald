package goald

import "fmt"

// ------------------------------------------------------------------------------------------------
// Easier error wrapping
// ------------------------------------------------------------------------------------------------

func Error(msg string, params ...interface{}) error {
	return fmt.Errorf(msg, params...)
}

func ErrorC(cause error, msg string, params ...interface{}) error {
	return fmt.Errorf(fmt.Sprintf(msg, params...)+" --==|| Cause: %w", cause)
}
