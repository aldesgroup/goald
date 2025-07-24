package core

import (
	"time"
)

// ------------------------------------------------------------------------------------------------
// Dates
// ------------------------------------------------------------------------------------------------

// custom time format, on top of the most used RFC3339
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

// Now returns the current time
func Now() *time.Time {
	t := time.Now()
	return &t
}
