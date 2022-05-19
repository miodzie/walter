package plugin

import (
	"fmt"
	"time"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/rss"
)

type RssMod struct {
	actions seras.Actions
	running bool
	feeds   rss.Feeds
	subs    rss.Subscriptions
	parser rss.Parser
}

func New(feeds rss.Feeds, subs rss.Subscriptions, parser rss.Parser) *RssMod {
	return &RssMod{feeds: feeds, subs: subs, parser: parser}
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
		// Disabled until I have admin users.
		// msg.Command("add_feed", mod.addFeed)
		msg.Command("subscribe", mod.subscribe)
	}

	return nil
}

func (mod *RssMod) checkFeeds() {
	p := rss.NewProcessor(mod.feeds, mod.subs, mod.parser)
	for mod.running {
		notifs, err := p.Handle()
		if err != nil {
			fmt.Println(err)
		}
		for _, notif := range notifs {
			msg := seras.Message{
				Channel: notif.Channel,
				Content: notif.String(),
			}
			mod.actions.Send(msg)
		}
		time.Sleep(time.Minute * 30)
	}
}

func (mod *RssMod) Stop() {
	mod.running = false
}
