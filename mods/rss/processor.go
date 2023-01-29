package rss

import (
	"github.com/miodzie/walter/log"
)

// TODO: ThrottledAnnouncer decorator

// TODO: CachedFetcher decorator
// TODO: RealTime Processor option?
// TODO: AnnouncementFormatter?

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
	Address() string
	Content() string

	// TODO: Figure out how to update subscription on delivery.
	OnDelivery() func()
	Sub() Subscription
}

func (p *Processor) Process() (chan Notification, error) {
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

func (p *Processor) process(feeds []*UserFeed, matcher *Matcher, notes chan Notification) {
	for _, uf := range feeds {
		feed, err := p.fetcher.Fetch(uf.Url)
		if err != nil {
			log.Error(err) // TODO: retry?
			continue
		}
		matches := matcher.Match(feed.Items)
		for _, m := range matches {
			notes <- m
		}
	}
	close(notes)
}

func (p *Processor) createMatcher() (*Matcher, error) {
	subs, err := p.storage.Subs(SearchParams{})
	var litSubs []Subscription
	for _, s := range subs {
		litSubs = append(litSubs, *s)
	}
	return NewMatcher(litSubs), err
}
