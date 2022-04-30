package policing

import (
	"fmt"
	"time"

	"github.com/miodzie/seras"
)

type PolicingMod struct {
	running bool
}

func New() *PolicingMod {
	return &PolicingMod{}
}

func (mod *PolicingMod) Name() string {
	return "police"
}

func (mod *PolicingMod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if IsSpam(msg) {
			err := actions.TimeoutUser(msg.Channel, msg.AuthorId, time.Now().Add(time.Minute*1))
			if err != nil {
				fmt.Printf("Failed to TimeoutUser: \"%s\"\n", err)
			}
		}
	}
	return nil
}

func (mod *PolicingMod) Stop() {
	mod.running = false
}
