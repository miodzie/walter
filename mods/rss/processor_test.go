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
	repository *InMemRepository
	fetcher    *StubFetcher

	userFeed *UserFeed
	item     Item
}

func (p *ProcessorSuite) PreTest(t *td.T, testName string) error {
	p.repository = NewInMemRepo()
	p.fetcher = NewStubFetcher()
	p.processor = NewProcessor(p.fetcher, p.repository)

	p.item = Item{Title: "The Go Blog", Link: "https://go.dev/blog", GUID: "1"}
	p.userFeed = &UserFeed{Id: 1, Url: "go.dev/blog"}
	feed := Feed{Items: []Item{p.item}}
	p.fetcher.Add(p.userFeed.Url, feed)

	return p.repository.AddFeed(p.userFeed)
}

func (p *ProcessorSuite) TestReturnsChannelOfNotifications(assert, require *td.T) {
	isaac := &Subscription{User: "isaac", Channel: "#general", FeedId: p.userFeed.Id}
	require.CmpNoError(p.repository.AddSub(isaac))
	jacob := &Subscription{User: "jacob", Channel: "#general", FeedId: p.userFeed.Id}
	require.CmpNoError(p.repository.AddSub(jacob))

	notes, err := p.processor.Process()
	require.CmpNoError(err)

	note := Notification{
		Channel: "#general",
		User:    "isaac",
		Item:    p.item,
		Feed:    *p.userFeed,
	}
	assert.Cmp(<-notes, note)
	note.User = "jacob"
	assert.Cmp(<-notes, note)
	assert.Cmp(<-notes, Notification{})
}

func TestRunProcessorSuite(t *testing.T) {
	tdsuite.Run(t, new(ProcessorSuite))
}
