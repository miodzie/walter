package irc

import (
	"fmt"

	seras "github.com/miodzie/seras"
	irc "github.com/thoj/go-ircevent"
	"sync"
)

type Connection struct {
	irc    *irc.Connection
	config *Config
	mu     sync.Mutex
}

func New(conf Config) (*Connection, error) {
	con := &Connection{
		irc:    irc.IRC(conf.Nick, conf.Username),
		config: &conf,
	}

	return con, nil
}

func (con *Connection) Server() string {
	return con.config.Server
}

func (con *Connection) SendMessage(msg seras.Message) error {
	return nil
}

func (con *Connection) Connect() (seras.Stream, error) {
	con.mu.Lock()
	defer con.mu.Unlock()
	err := con.irc.Connect(con.config.Server)
	if err != nil {
		return nil, err
	}
	stream := make(chan seras.Message)

	con.irc.AddCallback("*", func(event *irc.Event) {
		// Convert to Message
		// Send to stream
		// seras.Log(event.Message())
		// fmt.Println(event.Message())
		stream <- (&Message{event: event, irc: con}).ToMsg()
	})

	go func() {
		con.irc.Loop()
	}()

	return stream, nil
}

func (con *Connection) Close() error {
	con.mu.Lock()
	defer con.mu.Unlock()
	con.irc.Disconnect()
	fmt.Println("why hang")
	con.irc.ClearCallback("*")

	return nil
}
