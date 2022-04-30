package seras

var token string = "!"

func Token() string {
    return token
}

type Messenger interface {
	Send(Message) error
}

type Message struct {
	Content   string
	Arguments []string
	Channel   string
	// TODO: Rename these?
	Author   string
	AuthorId string
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
