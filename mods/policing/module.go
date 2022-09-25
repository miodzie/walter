package policing

import (
	"fmt"
	"time"

	"github.com/miodzie/seras"
)

type Mod struct {
	running bool
}

func New() *Mod {
	return &Mod{}
}

func (mod *Mod) Name() string {
	return "police"
}

func (mod *Mod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if IsSpam(msg) {
			err := actions.TimeoutUser(msg.Target, msg.Author.Id, time.Now().Add(time.Minute*1))
			if err != nil {
				fmt.Printf("Failed to TimeoutUser: \"%s\"\n", err)
			}
		}
	}
	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}
