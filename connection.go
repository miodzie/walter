package seras

import (
	"fmt"
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

var Connects map[string]Connection

func ParseConnects(cfg *Config) error {
	Connects = make(map[string]Connection)
	for name, con := range cfg.Connections {
		parser, ok := connectors[name]
		if !ok {
			return fmt.Errorf("unable to parse connector: %s", name)
		}
		var err error
		Connects[name], err = parser.Parse(con)
		if err != nil {
			return err
		}
	}

	return nil
}
