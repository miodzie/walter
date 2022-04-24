package seras

import "time"

type Stream <-chan Message

type Connection interface {
	Connect() (Stream, error)
	Close() error
}

type Admin interface {
	MuteUser(string, time.Duration)
}

