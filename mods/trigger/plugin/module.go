// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"log"

	"walter"
	"walter/mods/trigger"
)

type TriggerMod struct {
	actions walter.Actions
	running bool
	repo    trigger.Repository
}

func New(repo trigger.Repository) *TriggerMod {
	return &TriggerMod{repo: repo}
}

func (mod *TriggerMod) Start(stream walter.Stream, actions walter.Actions) error {
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
