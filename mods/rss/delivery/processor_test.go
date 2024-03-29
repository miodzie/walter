package delivery

import (
	"errors"
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"github.com/miodzie/walter/mods/rss"
	"testing"
)

// Fetch Feeds -> Notifications -> Announcements -> Announce

type ProcessorSuite struct {
	processor  *Processor
	repository *rss.InMemRepository
	fetcher    *StubFetcher

	userFeed *rss.UserFeed
	item     rss.Item
}

func (p *ProcessorSuite) PreTest(t *td.T, testName string) error {
	p.repository = rss.NewInMemRepo()
	p.fetcher = NewStubFetcher()
	p.processor = NewProcessor(p.fetcher, p.repository)

	p.item = rss.Item{Title: "The Go Blog", Link: "https://go.dev/blog", GUID: "1"}
	p.userFeed = &rss.UserFeed{Id: 1, Url: "go.dev/blog"}
	feed := rss.Feed{Items: []rss.Item{p.item}}
	p.fetcher.Add(p.userFeed.Url, feed)

	return p.repository.AddFeed(p.userFeed)
}

func (p *ProcessorSuite) TestReturnsChannelOfNotifications(assert, require *td.T) {
	isaac := &rss.Subscription{User: "isaac", Channel: "#general", FeedId: p.userFeed.Id}
	require.CmpNoError(p.repository.AddSub(isaac))
	jacob := &rss.Subscription{User: "jacob", Channel: "#general", FeedId: p.userFeed.Id}
	require.CmpNoError(p.repository.AddSub(jacob))

	notes, err := p.processor.Process()
	require.CmpNoError(err)

	assert.Cmp(getNote(notes), p.noteFromSub(isaac))
	assert.Cmp(getNote(notes), p.noteFromSub(jacob))
	assert.Nil(<-notes)
}

func (p *ProcessorSuite) TestIgnoresMatches(assert, require *td.T) {
	isaac := &rss.Subscription{User: "isaac", Channel: "#general", FeedId: p.userFeed.Id, Ignore: "Go"}
	require.CmpNoError(p.repository.AddSub(isaac))

	notes, err := p.processor.Process()
	require.CmpNoError(err)

	assert.Nil(<-notes)
}

func (p *ProcessorSuite) TestHidesSeenItems(assert, require *td.T) {
	isaac := &rss.Subscription{User: "isaac", Channel: "#general", FeedId: p.userFeed.Id}
	isaac.Remember(p.item)
	require.CmpNoError(p.repository.AddSub(isaac))

	notes, err := p.processor.Process()
	require.CmpNoError(err)

	assert.Nil(<-notes)
}

func (p *ProcessorSuite) TestItAddsSubscriptionRememberOnDeliveryHook(assert, require *td.T) {
	jacob := &rss.Subscription{User: "jacob", Channel: "#general", FeedId: p.userFeed.Id}
	require.CmpNoError(p.repository.AddSub(jacob))
	notes, err := p.processor.Process()
	require.CmpNoError(err)
	n := (<-notes).(Notification)
	p.repository.ForcedErr = errors.New("test")

	n.Deliver(func(address string, content string) error {
		return nil // it was delivered without error.
	})

	assert.Nil(p.repository.ForcedErr) // Confirms Repository call called
	assert.True(jacob.HasSeen(p.item))
}

func (p *ProcessorSuite) TestItDoesntMatchOtherFeedItems(assert, require *td.T) {
	isaac := &rss.Subscription{User: "isaac", Channel: "#general", FeedId: 2}
	require.CmpNoError(p.repository.AddFeed(&rss.UserFeed{Id: 2, Url: "go.dev/blog"}))
	require.CmpNoError(p.repository.AddSub(isaac))

	notes, err := p.processor.Process()
	require.CmpNoError(err)

	assert.Nil(<-notes)
}

func (p *ProcessorSuite) TestItOnlyNotifiesOnKeywords(assert, require *td.T) {
	isaac := &rss.Subscription{User: "isaac", Channel: "#general", FeedId: 1, Keywords: "potato"}
	require.CmpNoError(p.repository.AddSub(isaac))

	notes, err := p.processor.Process()
	require.CmpNoError(err)

	assert.Nil(<-notes)
}

func TestRunProcessorSuite(t *testing.T) {
	tdsuite.Run(t, new(ProcessorSuite))
}

func getNote(notes chan Deliverable) Notification {
	n2 := <-notes
	n := n2.(Notification)
	n.DeliveryHook = nil
	return n
}

func (p *ProcessorSuite) noteFromSub(sub *rss.Subscription) Notification {
	sub.HasSeen(rss.Item{})
	return Notification{
		Channel: sub.Channel,
		User:    sub.User,
		Item:    p.item,
		Feed:    *p.userFeed,
	}
}
