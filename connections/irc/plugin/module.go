package plugin

import (
	"fmt"
	"github.com/miodzie/seras"
	irc "github.com/thoj/go-ircevent"
)

type Mod struct {
	irc     *irc.Connection
	running bool
}

func New(irc *irc.Connection) *Mod {
	return &Mod{irc: irc}
}

func (mod *Mod) Name() string {
	return "irc"
}

func (mod *Mod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if msg.Code == "INVITE" {
			mod.irc.Join(msg.Arguments[1])
		}
		fmt.Printf("%+v\n", msg)
	}
	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}
