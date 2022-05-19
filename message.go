package seras

var token string = "!"

func Token() string {
	return token
}

type Messenger interface {
	Send(Message) error
	Reply(Message, string) error
}

type Message struct {
	Content    string
	Arguments  []string
	Channel    string
	AuthorId   string
	AuthorNick string
	// TODO: Refactor, this is quick fix.
	AuthorMention string
}

type Author struct {
	Id      string
	Nick    string
	Mention string
}

type MessageFormatter interface {
	Bold(string) string
	Italicize(string) string
	StrikeThrough(string) string
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
