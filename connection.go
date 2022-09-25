package seras

import (
	"time"
)

type Stream <-chan Message

type Connection interface {
	Connect() (Stream, error)
	Close() error
}

type Admin interface {
	// TODO: Refactor to accept just an Author struct.
	IsAdmin(userId string) bool
	TimeoutUser(channel string, user string, until time.Time) error
}
