package log

import (
	"fmt"
	"log"
	"os"
)

// Logger represents a simple logger with different log levels
type Logger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

// NewLogger creates a new logger
func NewLogger() *Logger {
	return &Logger{
		debug: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		warn:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime),
		error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatal: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		msg = fmt.Sprintf("%s %v", msg, fields)
	}
	l.debug.Println(msg)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		msg = fmt.Sprintf("%s %v", msg, fields)
	}
	l.info.Println(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		msg = fmt.Sprintf("%s %v", msg, fields)
	}
	l.warn.Println(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		msg = fmt.Sprintf("%s %v", msg, fields)
	}
	l.error.Println(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		msg = fmt.Sprintf("%s %v", msg, fields)
	}
	l.fatal.Fatalln(msg)
}

// Sync is a no-op function to maintain API compatibility
func (l *Logger) Sync() error {
	return nil
}
