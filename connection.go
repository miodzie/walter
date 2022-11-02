// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package walter

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

// BaseConnection is a base template for Connections to use.
type BaseConnection struct {
	NAME   string `toml:"name"`
	Type   string
	Admins []string
	Mods   map[string]any
}

func (c *BaseConnection) Name() string {
	return c.NAME
}

func (c *BaseConnection) SetName(name string) {
	c.NAME = name
}

type Admin interface {
	// IsAdmin TODO: Refactor to accept just an Author struct.
	IsAdmin(userId string) bool
	TimeoutUser(channel string, user string, until time.Time) error
}
