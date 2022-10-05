package repeater

import (
	"errors"
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/log"
)

var shepherd *Shepherd

func newSheepStream() seras.Stream {
	stream := make(chan seras.Message)
	shepherd.herd = append(shepherd.herd, stream)

	return stream
}

type ShepherdFactory struct {
}

func (b *ShepherdFactory) Create(a interface{}) (seras.Module, error) {
	if shepherd != nil {
		return nil, errors.New("there can only be one")
	}
	shepherd = &Shepherd{
		herd:    []chan seras.Message{},
		running: false,
	}

	return shepherd, nil
}

type Shepherd struct {
	herd    []chan seras.Message
	running bool
}

func (mod *Shepherd) Name() string {
	return "shepherd"
}

func (mod *Shepherd) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if msg.IsCommand("gm") {
			log.Info("GENTLEMEN! GOOD MORNING!")
			msg.Content = "g'mrn frens"
			for i, b := range mod.herd {
				log.Info("Sending her off!")
				b <- msg
				log.Info(i)
			}
			log.Info("Done!!")
		}
	}
	return nil
}

func (mod *Shepherd) Stop() {
	mod.running = false
}
