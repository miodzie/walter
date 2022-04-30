package rss

import (
	"fmt"
	"time"

	"github.com/miodzie/seras"
	//"github.com/mmcdole/gofeed"
)


const CRUNCHYROLL = "https://www.crunchyroll.com/rss/anime"

var listeners []*Listener

type RssMod struct {
	actions seras.Actions
	running bool
}

func New() *RssMod {
	return &RssMod{}
}
func (mod *RssMod) Name() string {
	return "rss"
}

func (mod *RssMod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	mod.actions = actions
	// Start Another routine to check RSS
	go mod.checkFeed()
	for mod.running {
		msg := <-stream
		if msg.Arguments[0] == "add_rss" {
			// TODO:
			actions.Send(seras.Message{})
		}
	}

	return nil
}

func (mod *RssMod) checkFeed() {
	for mod.running {
		for _, listener := range listeners {
			msgs, err := listener.Process()
			if err != nil {
				// TODO: log.
				fmt.Println(err)
				continue
			}
			for _, msg := range msgs {
				fmt.Println(msg)
				mod.actions.Send(msg)
				break
			}
		}
		time.Sleep(time.Minute * 30)
	}
}

func (mod *RssMod) Stop() {
	mod.running = false
}
