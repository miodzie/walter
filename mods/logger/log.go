package logger

import "github.com/miodzie/seras"

// Logger is meant for logging all incoming and outgoing seras.Messages.
// There is a default file option and ElasticSearch drivers available.
type Logger interface {
	Log(message seras.Message) error
}

type FileLogger struct {
	file string
}

func (f FileLogger) Log(message seras.Message) error {
	//TODO implement me
	panic("implement me")
}
