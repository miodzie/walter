package plugin

import (
	"github.com/miodzie/seras/log"
	"time"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/rss"
)

type RssMod struct {
	actions seras.Actions
	running bool
	Context
}

type Context struct {
	rss.Repository
	rss.Parser
	rss.Formatter
}

func New(ctx Context) *RssMod {
	return &RssMod{Context: ctx}
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
		msg.Command("add_feed", mod.addFeed)
		msg.Command("subscribe", mod.subscribe)
		msg.Command("unsubscribe", mod.unsubscribe)
		msg.Command("subs", mod.subs)
	}

	return nil
}

func (mod *RssMod) checkFeeds() {
	time.Sleep(time.Minute * 1)
	p := rss.NewProcessor(mod.Repository, mod.Parser)
	for mod.running {
		log.Info("Processing feed subscriptions...")
		notifs, err := p.Process()
		if err != nil {
			log.Error(err)
		}
		log.Infof("%d notifications found\n", len(notifs))
		for _, notif := range notifs {
			msg := seras.Message{
				Target:  notif.Channel,
				Content: mod.Format(*notif),
			}
			log.Debug(notif)
			mod.actions.Send(msg)
		}
		time.Sleep(time.Minute * 30)
	}
}

func (mod *RssMod) Stop() {
	mod.running = false
}
