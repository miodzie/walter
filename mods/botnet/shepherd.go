package botnet

import (
	"github.com/miodzie/seras"
)

var shepherd *Shepherd

type ShepherdFactory struct {
}

func (b *ShepherdFactory) Create(a interface{}) (seras.Module, error) {
	return nil, nil
}

type Shepherd struct {
	id        int
	broadcast []chan seras.Message
	running   bool
}

func NewController(repeaters []chan seras.Message) *Shepherd {
	return &Shepherd{broadcast: repeaters}
}

func (mod *Shepherd) Name() string {
	return "shepherd"
}

func (mod *Shepherd) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if msg.IsCommand("gm") {
			msg.Content = "g'mrn frens"
			for _, b := range mod.broadcast {
				b <- msg
			}
		}
	}
	return nil
}

func (mod *Shepherd) Stop() {
	mod.running = false
}
