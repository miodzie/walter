package rss

import (
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

// Fetch Feeds -> Notifications -> Announcements -> Announce

// TODO
// - ThrottledMessenger Decorator
// - Update Subscription seen items on successful delivery? Use hook?

type ProcessorSuite struct {
	processor  *Processor
	announcer  *StubAnnouncer
	repository *InMemRepository
	fetcher    *StubFetcher

	userFeed *UserFeed
	item     Item
}

func (p *ProcessorSuite) PreTest(t *td.T, testName string) error {
	p.repository = NewInMemRepo()
	p.fetcher = NewStubFetcher()
	p.announcer = &StubAnnouncer{}
	p.processor = NewProcessor(p.fetcher, p.repository, p.announcer)

	p.item = Item{Title: "The Go Blog", Link: "https://go.dev/blog", GUID: "1"}
	p.userFeed = &UserFeed{Id: 1, Url: "go.dev/blog"}
	feed := Feed{Items: []Item{p.item}}
	p.fetcher.Add(p.userFeed.Url, feed)

	return p.repository.AddFeed(p.userFeed)
}

func (p *ProcessorSuite) TestDeliversNotification(assert, require *td.T) {
	isaac := &Subscription{User: "isaac", Channel: "#general", FeedId: p.userFeed.Id}
	jacob := &Subscription{User: "jacob", Channel: "#general", FeedId: p.userFeed.Id}
	require.CmpNoError(p.repository.AddSub(jacob))
	require.CmpNoError(p.repository.AddSub(isaac))

	err := p.processor.Process()
	require.CmpNoError(err)

	if assert.Len(p.announcer.delivered, 1) {
		a := p.announcer.delivered[0]
		assert.Cmp(a.Room, "#general")
		assert.Cmp(a.Message, "The Go Blog - https://go.dev/blog : jacob,isaac")
	}
}

func TestRunProcessorSuite(t *testing.T) {
	tdsuite.Run(t, new(ProcessorSuite))
}

///////////////////////////////////////////////////////

type StubAnnouncer struct {
	delivered []Announcement
}

func (m *StubAnnouncer) Announce(announcements []Announcement) error {
	m.delivered = announcements
	return nil
}
