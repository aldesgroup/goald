package logging

import "log/slog"

// ----------------------------------------------------------------------------
// The logger interface, for all our logging needs
// ----------------------------------------------------------------------------

type ILogger interface {
	IsVerbose() bool
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(isFatal bool, msg string, args ...interface{})
	WithLevel(level slog.Level) ILogger
}

type LoggingType string

const (
	LoggingTypeTEXT LoggingType = "text"
	LoggingTypeJSON LoggingType = "json"
)

// ----------------------------------------------------------------------------
// Creating different types of loggers
// ----------------------------------------------------------------------------

func NewLogger(level slog.Level, contextString string, loggerType LoggingType, showFunctionNames bool) ILogger {
	switch loggerType {
	case LoggingTypeTEXT:
		return newTextLogger(level, contextString, showFunctionNames)
	case LoggingTypeJSON:
		return newJsonLogger(level, contextString, showFunctionNames)
	default:
		panic("Unknown logger type: " + string(loggerType))
	}
}
