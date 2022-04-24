package seras

import (
	"sync"
)

type Module interface {
	Loop(Stream, Actions) error
	Stop()
}

type Actions interface {
	Messenger
	Admin
}

type ModuleManager struct {
	modules []Module
	streams []chan Message
	actions Actions
}

func NewModManager(mods []Module, actions Actions) *ModuleManager {
	manager := &ModuleManager{
		modules: mods,
		actions: actions,
	}

	return manager
}

func (manager *ModuleManager) Run(stream Stream) {
	// Init mod streams, start them up.
	for _, mod := range manager.modules {
		modStream := make(chan Message)
		manager.streams = append(manager.streams, modStream)
		mod.Loop(modStream, manager.actions)
	}

	// Collect messages from stream, broadcast to mods.
	for msg := range stream {
		for _, ch := range manager.streams {
			ch <- msg
		}
	}
}

func (manager *ModuleManager) Stop() {
	for _, mod := range manager.modules {
		mod.Stop()
	}
	for _, stream := range manager.streams {
		close(stream)
	}
}

type BaseModule struct {
	Actions Actions
	Stream  Stream
	Running bool
	// See loopCheckExample()
	LoopCheck func()
	sync.Mutex
}

func loopCheckExample() {
	mod := &BaseModule{}
	mod.LoopCheck = func() {
		for mod.Running {
			msg := <-mod.Stream
			if msg.Content == "hello" {
				mod.Actions.Send(Message{Content: "Hi", Channel: msg.Channel})
			}
		}
	}
}

func (mod *BaseModule) Loop(stream Stream, actions Actions) error {
	mod.Lock()
	defer mod.Unlock()

	mod.Actions = actions
	mod.Stream = stream
	mod.Running = true
	go mod.LoopCheck()

	return nil
}

func (mod *BaseModule) Stop() {
	mod.Lock()
	defer mod.Unlock()

	mod.Running = false
}
