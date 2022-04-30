package seras

type Module interface {
	Name() string
	Start(Stream, Actions) error
	Stop()
}

type Actions interface {
	Messenger
	Admin
}
