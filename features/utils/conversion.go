package utils

import (
	"strconv"
	"time"
)

// ------------------------------------------------------------------------------------------------
// Conversion functions to harmonize the way we de/serialize data from/to Go values
// ------------------------------------------------------------------------------------------------

func StringToBool(valueAsString, context string) bool {
	value, errConv := strconv.ParseBool(valueAsString)
	PanicErrf(errConv, "'%s' is not a right bool value for context '%s'", valueAsString, context)
	return value
}

func StringToDate(valueAsString, context string) *time.Time {
	value, errConv := time.Parse(RFC3339Milli, valueAsString)
	PanicErrf(errConv, "'%s' is not a right date value for context '%s'", valueAsString, context)
	return &value
}

func StringToInt64(valueAsString, context string) int64 {
	value, errConv := strconv.ParseInt(valueAsString, 10, 64)
	PanicErrf(errConv, "'%s' is not a right int64 value for context '%s'", valueAsString, context)
	return value
}

func StringToInt(valueAsString, context string) int {
	value, errConv := strconv.Atoi(valueAsString)
	PanicErrf(errConv, "'%s' is not a right int value for context '%s'", valueAsString, context)
	return value
}

func StringToFloat32(valueAsString, context string) float32 {
	value, errConv := strconv.ParseFloat(valueAsString, 32)
	PanicErrf(errConv, "'%s' is not a right float32 value for context '%s'", valueAsString, context)
	return float32(value)
}

func StringToFloat64(valueAsString, context string) float64 {
	value, errConv := strconv.ParseFloat(valueAsString, 64)
	PanicErrf(errConv, "'%s' is not a right float64 value for context '%s'", valueAsString, context)
	return value
}

func BoolToString(value bool) string {
	return strconv.FormatBool(value)
}

func DateToString(value *time.Time) string {
	return value.Format(RFC3339Milli)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func Float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func Float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
