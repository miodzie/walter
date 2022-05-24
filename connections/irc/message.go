package irc

import (
	"github.com/miodzie/seras"
	irc "github.com/thoj/go-ircevent"
)

type Message struct {
	irc   seras.Connection
	event *irc.Event
}

func (msg *Message) ToMsg() seras.Message {
	return seras.Message{
		Content:   msg.Body(),
		Arguments: msg.Arguments(),
	}
}

func (msg *Message) Body() string {
	return msg.event.Message()
}

func (msg *Message) Arguments() []string {
	return msg.event.Arguments
}

func (msg *Message) Channel() string {
	return "implement me"
}

func (msg *Message) Server() *seras.Connection {
	return &msg.irc
}
