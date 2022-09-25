// Package log is meant to be used by modules to log any important information
// or errors to be reviewed by the bot owner.
//
// This package should not be confused with the logger module,
// which is specifically meant to log all Connections seras.Messages.
//
// log is a global state package, all Connections use the same logger.
// Likely each Connection will get their own Logger in the future.
package log

var logger Logger = ConsoleLogger{}

// Will log anything that is the level or above (warn, error, fatal, panic).
var level = TraceLevel

type Level int

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger is a minimal abstraction for generic log levels.
// You can implement different third party log libraries that you prefer.
type Logger interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func SetLogger(logg Logger) {
	logger = logg
}

func SetLevel(l Level) {
	level = l
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
func Info(args ...interface{}) {
	logger.Info(args...)
}
func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func Error(args ...interface{}) {
	logger.Error(args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
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

func (n NullLogger) Tracef(format string, args ...interface{}) {
}

func (n NullLogger) Debugf(format string, args ...interface{}) {
}

func (n NullLogger) Infof(format string, args ...interface{}) {
}

func (n NullLogger) Warnf(format string, args ...interface{}) {
}

func (n NullLogger) Errorf(format string, args ...interface{}) {
}
