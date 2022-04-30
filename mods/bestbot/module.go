package bestbot

import (
	"fmt"
	"regexp"

	"github.com/miodzie/seras"
)

type BestBotMod struct {
	running bool
}

func New() *BestBotMod {
	return &BestBotMod{}
}

func (mod *BestBotMod) Name() string {
	return "best_bot"
}

func (mod *BestBotMod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	fmt.Println("BEE, BOO BOO, BOP")
	for mod.running {
		msg := <-stream
		r, _ := regexp.Compile(`(?i)cs(?:go)?\?`)
		if r.MatchString(msg.Content) {
			actions.Reply(msg, "https://tenor.com/view/diego-eric-csgo-csgo-players-counter-strike-gif-22766889")
		}
		r, _ = regexp.Compile(`(?i)\bruck\b`)
		if r.MatchString(msg.Content) {
			actions.Reply(msg, "yes")
		}
		if msg.IsCommand("gentlemen") || msg.IsCommand("lenny") {
			actions.Reply(msg, "( ͡° ͜ʖ ͡° )")
		}
		if msg.IsCommand("ladies") {
			actions.Reply(msg, "( ͡° ͜ʖ ͡° )")
			actions.Reply(msg, "( ͡⊙ ͜ʖ ͡⊙ )")
			actions.Reply(msg, "( ͡◉ ͜ʖ ͡◉ )")
		}
	}

	return nil
}

func (mod *BestBotMod) Stop() {
	mod.running = false
}
