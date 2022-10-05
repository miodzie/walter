package botnet

import (
	"github.com/miodzie/seras"
)

type SheepFactory struct {
}

func (r *SheepFactory) Create(a interface{}) (seras.Module, error) {
	return nil, nil
}

type Sheep struct {
	stream  seras.Stream
	running bool
}

func NewBot(stream seras.Stream) *Sheep {
	return &Sheep{stream: stream}
}

func (mod *Sheep) Name() string {
	return "botnet"
}

func (mod *Sheep) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-mod.stream
		actions.Send(msg)
	}
	return nil
}

func (mod *Sheep) Stop() {
	mod.running = false
}
