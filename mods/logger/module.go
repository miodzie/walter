package logger

import (
	"github.com/miodzie/seras"
)

type Mod struct {
	running bool
}

func New() *Mod {
	return &Mod{}
}

func (mod *Mod) Name() string {
	return "logger"
}

func (mod *Mod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		// TODO: impl me
		//_ := <-stream
	}

	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}
