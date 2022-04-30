package rss

import (
	"time"

	"github.com/miodzie/seras"
)

type RssMod struct {
	actions seras.Actions
	running bool
	feeds   FeedRepository
	subs    SubscriptionRepository
}

func New(feeds FeedRepository, subs SubscriptionRepository) *RssMod {
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
		msg.Command("add_feed", mod.addFeed)
		msg.Command("subscribe", mod.subscribe)
	}

	return nil
}

func (mod *RssMod) checkFeeds() {
	for mod.running {
		feeds, err := mod.feeds.All()
		if err != nil {
			panic(err)
		}
		for _, feed := range feeds {
			_, err := mod.subs.GetByFeedId(feed.Id)
			if err != nil {
				panic(err)
			}

		}
		time.Sleep(time.Minute * 30)
	}
}

func (mod *RssMod) Stop() {
	mod.running = false
}
