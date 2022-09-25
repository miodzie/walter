package logger

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/log"
)

type Mod struct {
	running bool
	logger  Logger
}

func New(logger Logger) *Mod {
	return &Mod{logger: logger}
}

func (mod *Mod) Name() string {
	return "logger"
}

func (mod *Mod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		err := mod.logger.Log(msg)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}
