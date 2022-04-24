package mods

import (
	"fmt"
	"regexp"

	"github.com/miodzie/seras"
)

type BestBotMod struct {
  seras.BaseModule
}

func NewBestBotMod() *BestBotMod {
	mod := &BestBotMod{}
	mod.LoopCheck = func() {
		fmt.Println("BEE, BOO BOO, BOP")
		for mod.Running {
			msg := <-mod.Stream
			// fmt.Println("BestBot: MSG RECEIVED=" + msg.Content)
			r, _ := regexp.Compile(`[Cc]s(?:go)?\?`)
			if r.MatchString(msg.Content) {
				mod.Actions.Send(seras.Message{Content: "https://tenor.com/view/diego-eric-csgo-csgo-players-counter-strike-gif-22766889", Channel: msg.Channel})
			}
			if msg.Content == "gentlemen" || msg.Content == "lenny" {
				mod.Actions.Send(seras.Message{Content: "( ͡° ͜ʖ ͡° )", Channel: msg.Channel})
			}
			if msg.Content == "ladies" {
				mod.Actions.Send(seras.Message{Content: "( ͡° ͜ʖ ͡° )", Channel: msg.Channel})
				mod.Actions.Send(seras.Message{Content: "( ͡⊙ ͜ʖ ͡⊙ )", Channel: msg.Channel})
				mod.Actions.Send(seras.Message{Content: "( ͡◉ ͜ʖ ͡◉ )", Channel: msg.Channel})
			}
		}
	}

	return mod
}
