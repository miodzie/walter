package delivery

import (
	"github.com/miodzie/walter/log"
	"github.com/miodzie/walter/mods/rss"
)

type Processor struct {
	storage rss.Repository
	fetcher rss.Fetcher
}

func NewProcessor(f rss.Fetcher, r rss.Repository) *Processor {
	return &Processor{
		storage: r,
		fetcher: f,
	}
}

type Deliverable interface {
	// Where the Mail is being delivered.
	Address() string
	// The contents of the letter. You have to "open it",
	// e.g. call the func, to see what's inside.
	//
	// Returning nil means it was delivered successfully, therefor any post
	// delivery actions are OK to do. (e.g. update Subscription that item was seen/delivered)
	Deliver(deliver func(address, content string) error)
}

// TODO: Consider adding context.Context
func (p *Processor) Process() (chan Deliverable, error) {
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

func (p *Processor) process(feeds []*rss.UserFeed, deliveries chan Deliverable) {
	// TODO: Concurrent feed processing.
	for _, uf := range feeds {
		feed, err := p.fetcher.Fetch(uf.Url)
		if err != nil {
			log.Error(err) // TODO: retry?
			continue
		}
		subs, err := p.storage.Subs(rss.SearchParams{FeedId: uf.Id})
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

func (p *Processor) match(sub *rss.Subscription, items []rss.Item, matches chan Deliverable) {
	for _, item := range items {
		if sub.ShouldSee(item) {
			matches <- Notification{
				Channel: sub.Channel,
				Feed:    *sub.Feed,
				Item:    item,
				User:    sub.User,
				DeliveryHook: func() error {
					sub.Remember(item)
					return p.storage.UpdateSub(sub)
				},
			}
		}
	}
}
