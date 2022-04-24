package policing

import (
	"time"
	"github.com/miodzie/seras"
)

type PolicingMod struct {
	seras.BaseModule
}

func NewPolicingMod() *PolicingMod {
	mod := &PolicingMod{}
	mod.LoopCheck = func() {
		for mod.Running {
			msg := <-mod.Stream
			if IsSpam(msg) {
				mod.Actions.Send(seras.Message{Content: "bruh, shut up", Channel: msg.Channel})
				mod.Actions.TimeoutUser(msg.Channel, msg.Author, time.Now().Add(time.Minute * 1))
			}
		}
	}

	return mod
}
