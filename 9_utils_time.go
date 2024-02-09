package goald

import (
	"time"
)

// ------------------------------------------------------------------------------------------------
// Dates
// ------------------------------------------------------------------------------------------------

const dateFormatSECONDS = "2006-01-02 15:04:05-07"
const oneSECOND = 1 * time.Second

func Now() *time.Time {
	t := time.Now()
	return &t
}
