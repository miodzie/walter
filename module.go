package seras

type Module interface {
	Name() string
	Start(Stream, Actions) error
	Stop()
}

type Actions interface {
	Messenger
	MessageFormatter
	Admin
}

type ModuleManager struct {
	modules []Module
	actions Actions
	streams map[string]chan Message
}

func NewModManager(mods []Module, actions Actions) (*ModuleManager, error) {
	manager := &ModuleManager{
		modules: mods,
		actions: actions,
		streams: make(map[string]chan Message),
	}

	return manager, nil
}

func (manager *ModuleManager) Run(stream Stream) error {
	for _, mod := range manager.modules {
		modStream := make(chan Message)
		manager.streams[mod.Name()] = modStream
		go mod.Start(modStream, manager.actions)
	}

	for msg := range stream {
		for _, ch := range manager.streams {
			ch <- msg
		}
	}

	return nil
}

func (manager *ModuleManager) Stop() {
	for _, mod := range manager.modules {
		mod.Stop()
	}
	for _, stream := range manager.streams {
		close(stream)
	}
}
