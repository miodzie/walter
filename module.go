package seras

import (
	"sync"
)

type Module interface {
	Loop(Stream, Messenger) error
	Stop()
}

type ModuleManager struct {
  modules []Module
  streams []chan Message
  messenger Messenger
}

func NewModManager(mods []Module, messenger Messenger) *ModuleManager {
  manager := &ModuleManager{
    modules: mods,
    messenger: messenger,
  }

  return manager
}

func (manager *ModuleManager) Run(stream Stream) {
	// Init mod streams, start them up.
	for _, mod := range manager.modules {
		modStream := make(chan Message)
		manager.streams = append(manager.streams, modStream)
		mod.Loop(modStream, manager.messenger)
	}

	// Collect messages from stream, broadcast to mods.
	for msg := range stream {
		for _, ch := range manager.streams {
			ch <- msg
		}
	}
}

func (manager * ModuleManager) Stop() {
  for _, mod := range manager.modules {
    mod.Stop()
  }
  for _, stream := range manager.streams {
    close(stream)
  }
}

type BaseModule struct {
	Sender    Messenger
	Stream    Stream
	Running   bool
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
				mod.Sender.Send(Message{Content: "Hi", Channel: msg.Channel})
			}
		}
	}
}

func (mod *BaseModule) Loop(stream Stream, sender Messenger) error {
	mod.Lock()
	defer mod.Unlock()

	mod.Sender = sender
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
