package plugin

import (
	"fmt"
	"time"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/rss"
	"github.com/miodzie/seras/mods/rss/usecases"
)

type RssMod struct {
	actions seras.Actions
	running bool
	feeds   rss.Feeds
	subs    rss.Subscriptions
}

func New(feeds rss.Feeds, subs rss.Subscriptions) *RssMod {
	return &RssMod{feeds: feeds, subs: subs}
}
func (mod *RssMod) Name() string {
	return "rss"
}

func (mod *RssMod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	mod.actions = actions
	go mod.checkFeeds()
	for mod.running {
		msg := <-stream
		msg.Command("feeds", mod.showFeeds)
		// msg.Command("add_feed", mod.addFeed)
		msg.Command("subscribe", mod.subscribe)
	}

	return nil
}

func (mod *RssMod) checkFeeds() {
    // TODO: Replace parser.
	checkFeeds := usecases.NewCheckFeeds(mod.feeds, mod.subs, &rss.StubParser{})
	for mod.running {
		resp := checkFeeds.Handle()
		if resp.Error != nil {
			fmt.Println(resp.Error)
		}
        // TODO: Send Messages
		time.Sleep(time.Minute * 30)
	}
}

func (mod *RssMod) Stop() {
	mod.running = false
}
