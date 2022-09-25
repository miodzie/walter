// Package log is meant to be used by modules to log any important information
// or errors to be reviewed by the bot owner.
//
// This package should not be confused with the logger module,
// which is specifically meant to log all Connections seras.Messages.
//
// log is a global state package, all Connections use the same logger.
package log

var logger Logger = NullLogger{}

// Logger is a minimal abstraction for generic log levels.
// You can implement different third party log libraries that you prefer.
type Logger interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

func SetLogger(logg Logger) {
	logger = logg
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

type NullLogger struct {
}

func (n NullLogger) Trace(args ...interface{}) {
}

func (n NullLogger) Debug(args ...interface{}) {
}

func (n NullLogger) Info(args ...interface{}) {
}

func (n NullLogger) Warn(args ...interface{}) {
}

func (n NullLogger) Error(args ...interface{}) {
}
