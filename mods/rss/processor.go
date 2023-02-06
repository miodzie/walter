package rss

import (
	"github.com/miodzie/walter/log"
)

// TODO: CachedFetcher decorator
// TODO: RealTime processor option?
// TODO: AnnouncementFormatter? (DeliveryFormatter)

type processor struct {
	storage Repository
	fetcher Fetcher
}

func Processor(f Fetcher, r Repository) *processor {
	return &processor{
		storage: r,
		fetcher: f,
	}
}

type Deliverable interface {
	Address() string
	Deliver(deliver func(address, content string) error)
}

// TODO: Consider adding context.Context
// TODO: Return chan Deliverable
func (p *processor) Process() (chan Deliverable, error) {
	// TODO: Should only be active userFeeds that has subs.
	// Maybe at some point just have UserFeeds be actual user created feeds.
	userFeeds, err := p.storage.Feeds()
	if err != nil {
		return nil, err
	}

	deliveries := make(chan Deliverable)
	go p.process(userFeeds, deliveries)

	return deliveries, nil
}

func (p *processor) process(feeds []*UserFeed, deliveries chan Deliverable) {
	// TODO: Concurrent feed processing.
	for _, uf := range feeds {
		feed, err := p.fetcher.Fetch(uf.Url)
		if err != nil {
			log.Error(err) // TODO: retry?
			continue
		}
		subs, err := p.storage.Subs(SearchParams{FeedId: uf.Id})
		if err != nil {
			log.Error(err)
			continue
		}
		for _, sub := range subs {
			p.match(sub, feed.Items, deliveries)
		}
	}
	close(deliveries)
}

func (p *processor) match(sub *Subscription, items []Item, matches chan Deliverable) {
	for _, item := range items {
		if sub.ShouldSee(item) {
			matches <- Notification{
				Channel:      sub.Channel,
				Feed:         *sub.Feed,
				Item:         item,
				User:         sub.User,
				Subscription: *sub,
				DeliveryHook: func() {
					sub.Remember(item)
					_ = p.storage.UpdateSub(sub)
				},
			}
		}
	}
}
