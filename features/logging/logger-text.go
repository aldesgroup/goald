package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// ----------------------------------------------------------------------------
// Text logging
// ----------------------------------------------------------------------------

type textLogger struct {
	*slog.Logger
	prefix        string
	level         slog.Level
	withFuncNames bool
}

var _ ILogger = (*textLogger)(nil)

func (thisLogger *textLogger) IsVerbose() bool {
	return thisLogger.level < slog.LevelInfo
}

func (thisLogger *textLogger) Debug(msg string, args ...any) {
	thisLogger.Logger.Debug(thisLogger.makePrefix()+msg, args...)
}

func (thisLogger *textLogger) Info(msg string, args ...any) {
	thisLogger.Logger.Info(thisLogger.makePrefix()+msg, args...)
}

func (thisLogger *textLogger) Warn(msg string, args ...any) {
	thisLogger.Logger.Warn(thisLogger.makePrefix()+msg, args...)
	slog.Info("")
}

func (thisLogger *textLogger) Error(isFatal bool, msg string, args ...any) {
	thisLogger.Logger.Error(thisLogger.makePrefix()+msg, args...)
	if isFatal {
		thisLogger.Logger.Error("Stopping here!")
		os.Exit(1)
	}
}

func (thisLogger *textLogger) WithLevel(level slog.Level) ILogger {
	// since we're using the default logger, we use the default way of setting the level
	slog.SetLogLoggerLevel(level)

	return thisLogger
}

func (thisLogger *textLogger) makePrefix() string {
	// function name init
	functionName := ""

	// Getting the file and proc from the runtime
	programCounter, file, line, ok := runtime.Caller(2) // we're adding an offset to skip the internal calls
	file = filepath.Base(file)                          // we only need the file name, not the directory

	// Getting the calling function information from this
	if thisLogger.withFuncNames {
		if fn := runtime.FuncForPC(programCounter); ok && fn != nil {
			name := fn.Name()
			functionName = name[strings.LastIndexByte(name, '.')+1:] + ":"
		}
	}

	// let's build the prefix now
	return "[" + thisLogger.prefix + "][" + file + ":" + functionName + strconv.Itoa(line) + "] "
}

func newTextLogger(level slog.Level, prefix string, showFunctionNames bool) ILogger {
	return (&textLogger{
		Logger:        slog.Default(),
		prefix:        prefix,
		level:         level,
		withFuncNames: showFunctionNames,
	}).WithLevel(level)
}
