package plugin

import (
	"fmt"
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/sed"
)

func New() *Mod {
	return &Mod{}
}

type Mod struct {
	log     map[string][]seras.Message
	running bool
}

func (mod *Mod) Name() string {
	return "sed"
}

func (mod *Mod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	mod.log = make(map[string][]seras.Message)
	for mod.running {
		msg := <-stream
		// TODO: holy parse user input better merciful lawd
		s := sed.ParseSed(msg.Content)
		if s.Command != ".s" {
			mod.logMsg(msg)
		}
		if s.Command == ".s" {
			for i := len(mod.log[msg.Target]) - 1; i >= 0; i-- {
				m := mod.log[msg.Target][i]
				if s.HasMatch(m.Content) {
					r := fmt.Sprintf("%s meant to say: %s", m.Author.Nick, s.Replace(m.Content))
					actions.Reply(msg, r)
					break
				}
			}
		}
	}

	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}

func (mod *Mod) logMsg(msg seras.Message) {
	mod.log[msg.Target] = append(mod.log[msg.Target], msg)
	if len(mod.log[msg.Target]) > 20 {
		mod.log[msg.Target] = mod.log[msg.Target][1:]
	}
}

type ModFactory struct {
}

func (m ModFactory) Create(config interface{}) (seras.Module, error) {
	return New(), nil
}
