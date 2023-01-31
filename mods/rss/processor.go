package rss

import (
	"github.com/miodzie/walter/log"
)

// TODO: ThrottledAnnouncer decorator (ThrottledDeliveries)

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
	Content() string

	OnDelivery()
}

// TODO: Consider adding context.Context
func (p *processor) Process() (chan Notification, error) {
	// TODO: Should only be active userFeeds that has subs.
	// Maybe at some point just have UserFeeds be actual user created feeds.
	userFeeds, err := p.storage.Feeds()
	if err != nil {
		return nil, err
	}
	matcher, err := p.createMatcher()
	if err != nil {
		return nil, err
	}

	notes := make(chan Notification)
	go p.process(userFeeds, matcher, notes)

	return notes, nil
}

func (p *processor) process(feeds []*UserFeed, matcher *Matcher, notes chan Notification) {
	for _, uf := range feeds {
		feed, err := p.fetcher.Fetch(uf.Url)
		if err != nil {
			log.Error(err) // TODO: retry?
			continue
		}
		matches := matcher.Match(feed.Items)
		for _, m := range matches {
			m.OnDeliveryHook = func() {
				m.Subscription.Remember(m.Item)
				_ = p.storage.UpdateSub(&m.Subscription)
			}
			notes <- m
		}
	}
	close(notes)
}

func (p *processor) createMatcher() (*Matcher, error) {
	subs, err := p.storage.Subs(SearchParams{})
	var litSubs []Subscription
	for _, s := range subs {
		litSubs = append(litSubs, *s)
	}
	return NewMatcher(litSubs), err
}
