package log

var logger Logger

type any interface{}

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
