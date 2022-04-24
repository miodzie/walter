package seras

import "time"

type Stream <-chan Message

type Admin interface {
  MuteUser(string, time.Duration)
}

type Messenger interface {
	Send(Message) error
}

type Module interface {
	Loop(Stream, Messenger) error
	Stop()
}

type Connection interface {
	Connect() (Stream, error)
	Close() error
}

type Message struct {
	Content   string
	Arguments []string
	Channel   string
	Author    string
}

type NullMessenger struct {
}

func (messenger *NullMessenger) Send(msg Message) error {
	return nil
}


