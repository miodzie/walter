package bestbot

import (
	"regexp"

	"github.com/miodzie/seras"
)

type Mod struct {
	running bool
}

func New() *Mod {
	return &Mod{}
}

func (mod *Mod) Name() string {
	return "best_bot"
}

func (mod *Mod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		r, _ := regexp.Compile(`(?i)\bcs(?:go)?\?`)
		if r.MatchString(msg.Content) {
			actions.Reply(msg, "https://tenor.com/view/diego-eric-csgo-csgo-players-counter-strike-gif-22766889")
		}
		r, _ = regexp.Compile(`(?i)\bruck\b`)
		if r.MatchString(msg.Content) {
			actions.Reply(msg, "yes")
		}
		if msg.IsCommand("zoop") {
			actions.Reply(msg, "👉😎👉")
		}
		if msg.IsCommand("gentlemen") || msg.IsCommand("lenny") {
			actions.Reply(msg, "( ͡° ͜ʖ ͡°)")
		}
		if msg.IsCommand("lads") {
			actions.Reply(msg, "     🎩٩")
			actions.Reply(msg, "( ͡👁 ͜ʖ ͡👁)")
		}
	}
	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}
