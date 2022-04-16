package mods

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/miodzie/seras"
)

type BestBot struct {
	sender  seras.Messenger
	running bool
	stream  seras.Stream
	sync.Mutex
}

func NewBestBot() *BestBot {
	return &BestBot{running: false}
}

func (mod *BestBot) Loop(stream seras.Stream, sender seras.Messenger) error {
	fmt.Println("READY FOR ACTION!")
	mod.Lock()
	defer mod.Unlock()

	mod.sender = sender
	mod.stream = stream
	mod.running = true
	go mod.loop()

	return nil
}

func (mod *BestBot) Stop() {
	mod.Lock()
	defer mod.Unlock()

	mod.running = false
}

func (mod *BestBot) loop() {
	fmt.Println("BEE, BOO BOO, BOP")
	for mod.running {
		msg := <-mod.stream
		// fmt.Println("BestBot: MSG RECEIVED=" + msg.Content)
		r, _ := regexp.Compile(`cs(?:go)?\?`)
		if r.MatchString(msg.Content) {
			mod.sender.Send(seras.Message{Content: "https://tenor.com/view/diego-eric-csgo-csgo-players-counter-strike-gif-22766889", Channel: msg.Channel})
		}
		if msg.Content == "gentlemen" || msg.Content == "lenny" {
			mod.sender.Send(seras.Message{Content: "( ͡° ͜ʖ ͡° )", Channel: msg.Channel})
		}
		if msg.Content == "ladies" {
			mod.sender.Send(seras.Message{Content: "( ͡° ͜ʖ ͡° )", Channel: msg.Channel})
			mod.sender.Send(seras.Message{Content: "( ͡⊙ ͜ʖ ͡⊙ )", Channel: msg.Channel})
			mod.sender.Send(seras.Message{Content: "( ͡◉ ͜ʖ ͡◉ )", Channel: msg.Channel})
		}
	}
}
