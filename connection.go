package seras

type Stream <-chan Message

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
	Channel string
}

type ModuleManager struct {
	modules []Module
}

func (manager *ModuleManager) Run() {

}

type NullMessenger struct {
}

func (messenger *NullMessenger) Send(msg Message) error {
	return nil
}
