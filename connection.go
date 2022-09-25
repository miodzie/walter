package seras

import (
	"time"
)

type Stream <-chan Message

type Connection interface {
	Name() string
	// SetName TODO: Remove, anything can change it.
	SetName(string)
	Connect() (Stream, error)
	Close() error
}

type Admin interface {
	// IsAdmin TODO: Refactor to accept just an Author struct.
	IsAdmin(userId string) bool
	TimeoutUser(channel string, user string, until time.Time) error
}
