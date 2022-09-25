package drivers

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/logger"
)

type MultiLogger struct {
	loggers []logger.Logger
}

func NewMultiLogger(args ...logger.Logger) *MultiLogger {
	return &MultiLogger{loggers: args}
}

func (l MultiLogger) Log(message seras.Message) error {
	var err error
	for _, l := range l.loggers {
		err = l.Log(message)
	}

	// Just return the last error for now. w/e
	return err
}
