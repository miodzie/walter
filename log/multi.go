package log

func NewMultiLogger(args ...Logger) *MultiLogger {
	return &MultiLogger{loggers: args}
}

type MultiLogger struct {
	loggers []Logger
}

func (m MultiLogger) Trace(args ...interface{}) {
	for _, l := range m.loggers {
		l.Trace(args)
	}
}

func (m MultiLogger) Debug(args ...interface{}) {
	for _, l := range m.loggers {
		l.Debug(args)
	}
}

func (m MultiLogger) Info(args ...interface{}) {
	for _, l := range m.loggers {
		l.Info(args)
	}
}

func (m MultiLogger) Warn(args ...interface{}) {
	for _, l := range m.loggers {
		l.Warn(args)
	}
}

func (m MultiLogger) Error(args ...interface{}) {
	for _, l := range m.loggers {
		l.Error(args)
	}
}
