package plugin

import (
	"log"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/trigger"
)

type TriggerMod struct {
	actions seras.Actions
	running bool
	repo    trigger.Repository
}

func New(repo trigger.Repository) *TriggerMod {
	return &TriggerMod{repo: repo}
}

func (mod *TriggerMod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	mod.actions = actions
	for mod.running {
		msg := <-stream
		msg.Command("add_trig", mod.addTrig)
		trigs, err := mod.repo.All()
		if err != nil {
			log.Printf("err: failed retrieving triggers to check message: %s\n", err)
			continue
		}

		for _, t := range trigs {
			if t.Check(msg.Content) {
				mod.actions.Reply(msg, t.Reply)
			}
		}
	}

	return nil
}
