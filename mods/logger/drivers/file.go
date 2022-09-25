package drivers

import (
	"github.com/miodzie/seras"
)

type FileLogger struct {
}

func (l FileLogger) Log(message seras.Message) error {
	//TODO implement me
	panic("implement me")
	return nil
}
