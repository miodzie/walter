package logger

import "github.com/miodzie/seras"

// Logger is meant for logging all incoming and outgoing Messages on a per-connection
// basis.
// There is a default null, file option, and ElasticSearch drivers available.
type Logger interface {
	Log(message seras.Message) error
}

type NullLogger struct {
}

func (n NullLogger) Log(message seras.Message) error {
	return nil
}

type FileLogger struct {
	file string
}

func (f FileLogger) Log(message seras.Message) error {
	//TODO implement me
	panic("implement me")
}
