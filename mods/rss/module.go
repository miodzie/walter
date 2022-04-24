package rss

import (
	"github.com/miodzie/seras"
	"github.com/mmcdole/gofeed"
)

const CRUNCHYROLL = "https://www.crunchyroll.com/rss/anime"

type RssMod struct {
	seras.BaseModule
}

func (mod *RssMod) Name() string {
  return "rss"
}

func NewMod() *RssMod {
	mod := &RssMod{}
	mod.LoopCheck = func() {
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
		fd := gofeed.NewParser()
		feed, _ := fd.ParseURL(CRUNCHYROLL)
	}
}
