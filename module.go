package seras

import (
	"sync"
)

type Module interface {
	Loop(Stream, Messenger) error
	Stop()
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
