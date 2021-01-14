package util

import (
	"io"
	"log"
)

// Logger wrapping 4 level loggers
type Logger struct {
	VerboseEnabled bool
	verboseLogger  *log.Logger
	infoLogger     *log.Logger
	warnLogger     *log.Logger
	errorLogger    *log.Logger
}

// Verbose log
func (l Logger) Verbose(v ...interface{}) {
	if l.VerboseEnabled {
		l.verboseLogger.Println(v...)
	}
}

// Info level log
func (l Logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

// Warn level log
func (l Logger) Warn(v ...interface{}) {
	l.warnLogger.Println(v...)
}

// Error level log
func (l Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}

// Fatal error level log and os.Exit(1)
func (l Logger) Fatal(v ...interface{}) {
	l.errorLogger.Fatal(v...)
}

// InitLogger with level loggers
func InitLogger(
	verboseHandle io.Writer,
	infoHandle io.Writer,
	warnHandle io.Writer,
	errorHandle io.Writer) *Logger {

	return &Logger{
		verboseLogger: log.New(verboseHandle, "[Verbose]: ", 0),
		infoLogger:    log.New(infoHandle, "[Info]: ", 0),
		warnLogger:    log.New(warnHandle, "[Warn]: ", 0),
		errorLogger:   log.New(errorHandle, "[Error]: ", 0),
	}
}
