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
	thisLogger.Logger = makeJsonLogger(level, thisLogger.contextString, thisLogger.withFuncNames)

	return thisLogger
}

func makeJsonLogger(level slog.Level, contextString string, showFunctionNames bool) *slog.Logger {
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

	return slog.New(slog.NewJSONHandler(os.Stdout, loggerOpts)).With("_ctx", contextString)
}

func newJsonLogger(level slog.Level, contextString string, showFunctionNames bool) ILogger {
	return &jsonLogger{
		Logger:        makeJsonLogger(level, contextString, showFunctionNames),
		level:         level,
		contextString: contextString,
		withFuncNames: showFunctionNames,
	}
}
