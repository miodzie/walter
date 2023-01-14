package usecases

import (
	"github.com/miodzie/walter/mods/rss"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RemoveFeedSuite struct {
	suite.Suite
	removeFeed *RemoveFeed
	repository rss.Repository
}

func (t *RemoveFeedSuite) SetupTest() {
	t.repository = rss.NewInMemRepo()
	t.removeFeed = NewRemoveFeed(t.repository)
}

func (t *RemoveFeedSuite) TestItRemovesAllSubscriptionsToThatFeed() {
	feed := &rss.UserFeed{Id: 42, Name: "some_feed"}
	t.NoError(t.repository.AddFeed(feed))
	sub := &rss.Subscription{User: "isaac", FeedId: feed.Id}
	sub2 := &rss.Subscription{User: "jacob", FeedId: feed.Id}
	t.NoError(t.repository.AddSub(sub))
	t.NoError(t.repository.AddSub(sub2))

	err := t.removeFeed.Remove(feed.Name)

	if t.NoError(err) {
		subs, _ := t.repository.Subs(rss.SearchParams{FeedId: feed.Id})
		t.Len(subs, 0)
	}
}

func (t *RemoveFeedSuite) TestItRemovesAFeed() {
	feed := &rss.UserFeed{Id: 42, Name: "some_feed"}
	t.NoError(t.repository.AddFeed(feed))

	err := t.removeFeed.Remove(feed.Name)

	if t.NoError(err) {
		f, _ := t.repository.Feeds()
		t.Len(f, 0)
	}
}

func TestRemoveFeedSuite(t *testing.T) {
	suite.Run(t, new(RemoveFeedSuite))
}
