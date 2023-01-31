package rss

import (
	"errors"
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

// Fetch Feeds -> Notifications -> Announcements -> Announce

// TODO
// - ThrottledMessenger Decorator
// - Update Subscription seen items on successful delivery? Use hook?

type ProcessorSuite struct {
	processor  *processor
	repository *InMemRepository
	fetcher    *StubFetcher

	userFeed *UserFeed
	item     Item
}

func (p *ProcessorSuite) PreTest(t *td.T, testName string) error {
	p.repository = NewInMemRepo()
	p.fetcher = NewStubFetcher()
	p.processor = Processor(p.fetcher, p.repository)

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

	isaac.makeSeenMap()
	jacob.makeSeenMap()

	notes, err := p.processor.Process()
	require.CmpNoError(err)

	n1 := <-notes
	n1.OnDeliveryHook = nil
	assert.Cmp(n1, Notification{
		Channel:      "#general",
		User:         "isaac",
		Item:         p.item,
		Feed:         *p.userFeed,
		Subscription: *isaac,
	})
	n2 := <-notes
	n2.OnDeliveryHook = nil
	assert.Cmp(n2, Notification{
		Channel:      "#general",
		User:         "jacob",
		Item:         p.item,
		Feed:         *p.userFeed,
		Subscription: *jacob,
	})
	assert.Cmp(<-notes, Notification{})
}

func (p *ProcessorSuite) TestItAddsSubscriptionRememberOnDeliveryHook(assert, require *td.T) {
	isaac := &Subscription{User: "isaac", Channel: "#general", FeedId: p.userFeed.Id}
	require.CmpNoError(p.repository.AddSub(isaac))
	notes, err := p.processor.Process()
	require.CmpNoError(err)
	n := <-notes
	p.repository.forcedErr = errors.New("test")

	n.OnDelivery()

	assert.Nil(p.repository.forcedErr) // Confirms Repository call called
	assert.True(n.Subscription.HasSeen(p.item))
}

func TestRunProcessorSuite(t *testing.T) {
	tdsuite.Run(t, new(ProcessorSuite))
}
