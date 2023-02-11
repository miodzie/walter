package rss

import (
	"github.com/miodzie/walter/log"
)

type Processor struct {
	storage Repository
	fetcher Fetcher
}

func NewProcessor(f Fetcher, r Repository) *Processor {
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

func (p *Processor) process(feeds []*UserFeed, deliveries chan Deliverable) {
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

func (p *Processor) match(sub *Subscription, items []Item, matches chan Deliverable) {
	for _, item := range items {
		if sub.ShouldSee(item) {
			matches <- Notification{
				Channel: sub.Channel,
				Feed:    *sub.Feed,
				Item:    item,
				User:    sub.User,
				// TODO: I still don't like having this here, could maybe have a Subscription method?
				DeliveryHook: func() {
					log.Debug(item.GUID, item)
					sub.Remember(item)
					_ = p.storage.UpdateSub(sub)
				},
			}
		}
	}
}
