package botnet

import (
	"errors"
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/log"
)

type SheepFactory struct {
}

func (r *SheepFactory) Create(a interface{}) (seras.Module, error) {
	if shepherd == nil {
		return nil, errors.New("help i'm lost")
	}
	sheep := &Sheep{
		instructions: newSheepStream(),
		running:      false,
	}

	return sheep, nil
}

type Sheep struct {
	instructions seras.Stream
	running      bool
}

func (mod *Sheep) Name() string {
	return "sheep"
}

func (mod *Sheep) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	go func() {
		for mod.running {
			log.Info("TALLY HO!")
			msg := <-mod.instructions
			log.Info("BAAAAAA!!!")
			log.Info(msg)
			actions.Send(msg)
		}
	}()
	// might as well drain the queue for memory leaks.
	for mod.running {
		_ = <-stream
	}
	return nil
}

func (mod *Sheep) Stop() {
	mod.running = false
}
