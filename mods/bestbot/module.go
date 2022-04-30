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
			actions.Send(seras.Message{Content: "https://tenor.com/view/diego-eric-csgo-csgo-players-counter-strike-gif-22766889", Channel: msg.Channel})
		}
		r, _ = regexp.Compile(`(?i)\bruck\b`)
		if r.MatchString(msg.Content) {
			actions.Send(seras.Message{Content: "yes", Channel: msg.Channel})
		}
		if msg.Content == "gentlemen" || msg.Content == "lenny" {
			actions.Send(seras.Message{Content: "( ͡° ͜ʖ ͡° )", Channel: msg.Channel})
		}
		if msg.Content == "ladies" {
			actions.Send(seras.Message{Content: "( ͡° ͜ʖ ͡° )", Channel: msg.Channel})
			actions.Send(seras.Message{Content: "( ͡⊙ ͜ʖ ͡⊙ )", Channel: msg.Channel})
			actions.Send(seras.Message{Content: "( ͡◉ ͜ʖ ͡◉ )", Channel: msg.Channel})
		}
	}

	return nil
}

func (mod *BestBotMod) Stop() {
	mod.running = false
}
