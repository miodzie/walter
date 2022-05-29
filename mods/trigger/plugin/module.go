package plugin

import "github.com/miodzie/seras"

type TriggerMod struct {
	actions seras.Actions
	running bool
}

func New() *TriggerMod {
	return &TriggerMod{}
}

func (mod *TriggerMod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	mod.actions = actions
	for mod.running {
		msg := <-stream
		msg.Command("feeds", mod.addTrig)
	}

	return nil
}
