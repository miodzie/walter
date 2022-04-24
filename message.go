package seras

type Messenger interface {
	Send(Message) error
}

type Message struct {
	Content   string
	Arguments []string
	Channel   string
  // TODO: Rename these?
	Author    string
	AuthorId  string
}

type NullMessenger struct{}

func (messenger *NullMessenger) Send(msg Message) error {
	return nil
}
