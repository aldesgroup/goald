package core

import (
	"strconv"
	"time"
)

// ------------------------------------------------------------------------------------------------
// Conversion functions to harmonize the way we de/serialize data from/to Go values
// ------------------------------------------------------------------------------------------------

// StringToBool converts a string to a boolean
func StringToBool(valueAsString, context string) bool {
	value, errConv := strconv.ParseBool(valueAsString)
	PanicMsgIfErr(errConv, "'%s' is not a right bool value for context '%s'", valueAsString, context)
	return value
}

// StringToDate converts a string to a date
func StringToDate(valueAsString, context string) *time.Time {
	value, errConv := time.Parse(RFC3339Milli, valueAsString)
	PanicMsgIfErr(errConv, "'%s' is not a right date value for context '%s'", valueAsString, context)
	return &value
}

// StringToInt64 converts a string to a int64
func StringToInt64(valueAsString, context string) int64 {
	value, errConv := strconv.ParseInt(valueAsString, 10, 64)
	PanicMsgIfErr(errConv, "'%s' is not a right int64 value for context '%s'", valueAsString, context)
	return value
}

// StringToInt converts a string to a int
func StringToInt(valueAsString, context string) int {
	value, errConv := strconv.Atoi(valueAsString)
	PanicMsgIfErr(errConv, "'%s' is not a right int value for context '%s'", valueAsString, context)
	return value
}

// StringToFloat32 converts a string to a float32
func StringToFloat32(valueAsString, context string) float32 {
	value, errConv := strconv.ParseFloat(valueAsString, 32)
	PanicMsgIfErr(errConv, "'%s' is not a right float32 value for context '%s'", valueAsString, context)
	return float32(value)
}

// StringToFloat64 converts a string to a float64
func StringToFloat64(valueAsString, context string) float64 {
	value, errConv := strconv.ParseFloat(valueAsString, 64)
	PanicMsgIfErr(errConv, "'%s' is not a right float64 value for context '%s'", valueAsString, context)
	return value
}

// BoolToString converts a boolean to a string
func BoolToString(value bool) string {
	return strconv.FormatBool(value)
}

// DateToString converts a date to a string
func DateToString(value *time.Time) string {
	return value.Format(RFC3339Milli)
}

// Int64ToString converts a int64 to a string
func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

// IntToString converts a int to a string
func IntToString(value int) string {
	return strconv.Itoa(value)
}

// Float32ToString converts a float32 to a string
func Float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

// Float64ToString converts a float64 to a string
func Float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
