package util

import (
	"io"
	"log"
)

// Logger wrapping 4 level loggers
type Logger struct {
	traceLogger, infoLogger, warnLogger, errorLogger *log.Logger
}

// Trace level log
func (l Logger) Trace(v ...interface{}) {
	l.traceLogger.Println(v...)
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

// InitLogger with level loggers
func InitLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warnHandle io.Writer,
	errorHandle io.Writer) *Logger {

	return &Logger{
		traceLogger: log.New(traceHandle, "[Trace]: ", log.Ldate|log.Ltime),
		infoLogger:  log.New(infoHandle, "[Info]: ", log.Ldate|log.Ltime),
		warnLogger:  log.New(warnHandle, "[Warn]: ", log.Ldate|log.Ltime),
		errorLogger: log.New(errorHandle, "[Error]: ", log.Ldate|log.Ltime),
	}
}
