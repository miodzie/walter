package policing

import (
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
				mod.Sender.Send(seras.Message{Content: "bruh, shut up", Channel: msg.Channel})
			}
		}
	}

	return mod
}
