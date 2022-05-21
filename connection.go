package seras

import "time"

type Stream <-chan Message

var connections map[string]ConfigParser

type Connection interface {
	Connect() (Stream, error)
	Close() error
}

type Admin interface {
	TimeoutUser(channel string, user string, until time.Time) error
}

func RegisterConnection(conType string, parser ConfigParser) error {
	connections[conType] = parser
	return nil
}
