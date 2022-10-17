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

// ConnectionConfig is a base template for Connections to use.
type ConnectionConfig struct {
	NAME   string `toml:"name"`
	Type   string
	Admins []string
	Mods   map[string]any
}

func (c *ConnectionConfig) Name() string {
	return c.NAME
}

func (c *ConnectionConfig) SetName(name string) {
	c.NAME = name
}

type Admin interface {
	// IsAdmin TODO: Refactor to accept just an Author struct.
	IsAdmin(userId string) bool
	TimeoutUser(channel string, user string, until time.Time) error
}
