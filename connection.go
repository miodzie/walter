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
	IsAdmin(userId string) bool
	TimeoutUser(channel string, user string, until time.Time) error
}
