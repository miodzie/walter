package seras

import "time"

var token = "!"

func Token() string {
	return token
}

type Messenger interface {
	Send(Message) error
	Reply(Message, string) error
}

type Message struct {
	Content   string
	Arguments []string
	Target    string
	Author    Author
	Code      string
	// Name of the Connection it came from.
	ConnectionName string
	// JSON for Discord, Raw for IRC.
	Raw       string
	Timestamp time.Time
}

type Author struct {
	Id   string // Host in IRC, User ID in Discord.
	Nick string

	// Mention is starting to turn into an important thing for plugins,
	// Refactor the name.
	// Maybe UniqueMention?
	Mention string // TODO: Refactor?, this is quick fix to get mentions working in Discord.
}

type MessageFormatter interface {
	Bold(string) string
	Italicize(string) string
}

func (msg *Message) Command(command string, call func(Message)) {
	if msg.IsCommand(command) {
		call(*msg)
	}
}

func (msg *Message) IsCommand(command string) bool {
	return token+command == msg.Arguments[0]
}

type NullMessenger struct{}

func (messenger *NullMessenger) Send(msg Message) error {
	return nil
}
