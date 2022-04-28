package rss

import (
	"fmt"
	"time"

	"github.com/miodzie/seras"
	"github.com/mmcdole/gofeed"
)

const CRUNCHYROLL = "https://www.crunchyroll.com/rss/anime"

var listeners []*Listener

type RssMod struct {
	seras.BaseModule
}

func (mod *RssMod) Name() string {
	return "rss"
}

func New() *RssMod {
	mod := &RssMod{}
	mod.Run = func() {
		// Start Another routine to check RSS
		go mod.checkFeed()
		for mod.Running {
			msg := <-mod.Stream
			if msg.Arguments[0] == "add_rss" {
				// TODO:
				mod.Actions.Send(seras.Message{})
			}
		}
	}

	return mod
}

func (mod *RssMod) checkFeed() {
	for mod.Running {
		for _, listener := range listeners {
			msgs, err := listener.Process()
			if err != nil {
                // TODO: log.
			    fmt.Println(err)
				continue
			}
			for _, msg := range msgs {
                fmt.Println(msg)
				mod.Actions.Send(msg)
                break
			}
		}
		time.Sleep(time.Minute * 1)
	}
}
