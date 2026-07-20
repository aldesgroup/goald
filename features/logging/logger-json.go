package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

// ----------------------------------------------------------------------------
// JSON logging
// ----------------------------------------------------------------------------

type jsonLogger struct {
	*slog.Logger
	baseLogger    *slog.Logger // the handler-backed logger, before the '_ctx' attribute gets added - kept around so we can cheaply derive loggers with a different prefix
	level         slog.Level
	contextString string
	withFuncNames bool
}

var _ ILogger = (*jsonLogger)(nil)

func (thisLogger *jsonLogger) IsVerbose() bool {
	return thisLogger.level < slog.LevelInfo

}

func (thisLogger *jsonLogger) Error(isFatal bool, msg string, args ...any) {
	thisLogger.Logger.Error(msg, args...)
	if isFatal {
		thisLogger.Logger.Error("Stopping here!")
		os.Exit(1)
	}
}

func (thisLogger *jsonLogger) WithLevel(level slog.Level) ILogger {
	thisLogger.level = level
	thisLogger.baseLogger = makeJsonLogger(level, thisLogger.withFuncNames)
	thisLogger.Logger = thisLogger.baseLogger.With("_ctx", thisLogger.contextString)

	return thisLogger
}

func (thisLogger *jsonLogger) WithPrefix(prefix string) ILogger {
	// reusing the same handler-backed logger, and just deriving a new one with a different '_ctx' attribute,
	// instead of rebuilding a whole new JSON handler + options from scratch
	return &jsonLogger{
		Logger:        thisLogger.baseLogger.With("_ctx", prefix),
		baseLogger:    thisLogger.baseLogger,
		level:         thisLogger.level,
		contextString: prefix,
		withFuncNames: thisLogger.withFuncNames,
	}
}

func makeJsonLogger(level slog.Level, showFunctionNames bool) *slog.Logger {
	loggerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: showFunctionNames,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.SourceKey {
				return a
			}

			source, ok := a.Value.Any().(*slog.Source)
			if !ok || source == nil {
				return a
			}

			a.Value = slog.StringValue(filepath.Base(source.File) + ":" + strconv.Itoa(source.Line))
			return a
		},
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, loggerOpts))
}

func newJsonLogger(level slog.Level, contextString string, showFunctionNames bool) ILogger {
	baseLogger := makeJsonLogger(level, showFunctionNames)

	return &jsonLogger{
		Logger:        baseLogger.With("_ctx", contextString),
		baseLogger:    baseLogger,
		level:         level,
		contextString: contextString,
		withFuncNames: showFunctionNames,
	}
}
