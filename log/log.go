// Package log is meant to be used by modules to log any important information
// or errors to be reviewed by the bot owner.
//
// This package should not be confused with the logger module,
// which is specifically meant to log all Connections seras.Messages.
//
// log is a global state package, all Connections use the same logger.
package log

var logger Logger

type any interface{}

// Logger is an abstraction for generic log levels.
// You can implement different third party log libraries that you prefer.
type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Error(err error, args ...any)
	Warn(msg string, args ...any)
}

func SetLogger(logg Logger) {
	logger = logg
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args)
}

func Error(err error, args ...any) {
	logger.Error(err, args)
}

type NullLogger struct {
}

func (n NullLogger) Info(msg string, args ...any) {
}

func (n NullLogger) Debug(msg string, args ...any) {
}

func (n NullLogger) Error(err error, args ...any) {
}

func (n NullLogger) Warn(msg string, args ...any) {
}
