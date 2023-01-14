package rss

type Processor struct {
	storage   Repository
	fetcher   Fetcher
	messenger Messenger
}

func NewProcessor(f Fetcher, r Repository, m Messenger) *Processor {
	return &Processor{
		storage:   r,
		fetcher:   f,
		messenger: m,
	}
}

func (p *Processor) Process() error {
	// TODO: Should only be active userFeeds that has subs.
	// Maybe at some point just have UserFeeds be actual user feeds.
	userFeeds, _ := p.storage.Feeds()
	matcher, _ := p.createMatcher()

	var notes []Notification
	for _, uf := range userFeeds {
		f, _ := p.fetcher.Fetch(uf.Url)
		notes = append(notes, matcher.Match(f.Items)...)
	}

	organizer := AnnouncementOrganizer{}
	announcements := organizer.Organize(notes)

	return p.messenger.Deliver(announcements)
}

func (p *Processor) createMatcher() (*Matcher, error) {
	subs, err := p.storage.Subs(SearchParams{})
	var litSubs []Subscription
	for _, s := range subs {
		litSubs = append(litSubs, *s)
	}
	return NewMatcher(litSubs), err
}
