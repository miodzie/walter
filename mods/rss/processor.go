package rss

import "github.com/miodzie/walter/log"

// TODO: CachedFetcher decorator
// TODO: ThrottledAnnouncer decorator

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

	organizer := AnnouncementOrganizer{}
	announcements := organizer.Organize(notes)

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
