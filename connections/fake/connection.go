package fake

import (
	"github.com/miodzie/seras"
)

type Connection struct {
	messages []seras.Message
	stream   chan seras.Message
	Sender   *Messenger
}

func NewConnection() *Connection {
	con := &Connection{
		messages: []seras.Message{},
		stream:   make(chan seras.Message, 10),
	}

	return con
}

func (con *Connection) Server() string {
	return "fake"
}

func (con *Connection) Connect() (seras.Stream, error) {
	return con.stream, nil
}

func (con *Connection) Close() error {
	close(con.stream)
	return nil
}

type Messenger struct {
	stream chan<- seras.Message
}

func NewMessenger(con *Connection) *Messenger {
	return &Messenger{stream: con.stream}
}

func (messenger *Messenger) Send(msg seras.Message) error {
	messenger.stream <- msg
	return nil
}
