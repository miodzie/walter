package test

import (
	"github.com/miodzie/seras"
)

type Connection struct {
	messages []seras.Message
	stream   chan seras.Message
}

func NewConnection() *Connection {
	con := &Connection{
		messages: []seras.Message{},
		stream:   make(chan seras.Message, 10),
	}

	return con
}

func (con *Connection) Name() string {
	return "test"
}

func (con *Connection) SetName(s string) {
}

func (con *Connection) Server() string {
	return "test"
}

func (con *Connection) Connect() (seras.Stream, error) {
	return con.stream, nil
}

func (con *Connection) Close() error {
	close(con.stream)
	return nil
}

func (con *Connection) Send(msg seras.Message) error {
	con.stream <- msg
	return nil
}
