package rss

import "github.com/miodzie/walter/log"

// TODO: ThrottledAnnouncer decorator

// TODO: CachedFetcher decorator
// TODO: RealTime Processor option?
// TODO: AnnouncementFormatter?

type Processor struct {
	storage   Repository
	fetcher   Fetcher
	announcer Announcer
}

func NewProcessor(f Fetcher, r Repository, m Announcer) *Processor {
	return &Processor{
		storage:   r,
		fetcher:   f,
		announcer: m,
	}
}

func (p *Processor) Process() error {
	// TODO: Should only be active userFeeds that has subs.
	// Maybe at some point just have UserFeeds be actual user feeds.
	userFeeds, err := p.storage.Feeds()
	if err != nil {
		return err
	}
	matcher, err := p.createMatcher()
	if err != nil {
		return err
	}

	var notes []Notification
	for _, uf := range userFeeds {
		f, err := p.fetcher.Fetch(uf.Url)
		if err != nil {
			log.Error(err)
			continue
		}
		notes = append(notes, matcher.Match(f.Items)...)
	}

	// I can abstract this out into a pipeline that returns a channel of
	// Notifications, this enables more modularity.
	// I can then have that RealTimeProcessor,
	//that's constantly polling and sending new Notifications fresh off the channel.
	// While this aggregate into announcements can be a separate pipeline,
	//off the same base.

	organizer := AnnouncementOrganizer{}
	announcements := organizer.Organize(notes)

	// TODO: Add a "transaction" for subscriptions to fail on save if they're not
	// delivered?
	// e.g.
	if p.announcer.Announce(announcements) != nil {
		// rollback _all_ announcements?
		// what if some were delivered, some not?
	}

	return p.announcer.Announce(announcements)
}

func (p *Processor) createMatcher() (*Matcher, error) {
	subs, err := p.storage.Subs(SearchParams{})
	var litSubs []Subscription
	for _, s := range subs {
		litSubs = append(litSubs, *s)
	}
	return NewMatcher(litSubs), err
}
