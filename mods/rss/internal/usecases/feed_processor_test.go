// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"errors"
	. "github.com/miodzie/walter/mods/rss/internal/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type FeedProcessorSuite struct {
	processor  *FeedProcessor
	repository Repository
	feed       *Feed
	item       *ParsedItem
	suite.Suite
}

func (s *FeedProcessorSuite) SetupTest() {
	s.item = &ParsedItem{Title: "bar", GUID: "1"}
	parsed := &ParsedFeed{Items: []*ParsedItem{s.item}}
	s.repository = NewInMemRepo()
	s.processor = NewFeedProcessor(s.repository, &StubParser{Parsed: parsed})
	s.feed = &Feed{Id: 1}
	_ = s.processor.repository.AddFeed(s.feed)
}

// TODO: This test does nothing. lol
func (s *FeedProcessorSuite) TestSubscriptionsUpdated() {
	alice := Subscription{User: "alice", Channel: "#chat2", FeedId: s.feed.Id}
	aliceCopy := alice
	s.Nil(s.repository.AddSub(&aliceCopy))

	_, _ = s.processor.Process()

	r, _ := s.repository.Subs(SearchParams{User: "alice"})
	if s.Len(r, 1) {
		s.True(r[0].SeenItems[s.item.GUID])
	}
}

func (s *FeedProcessorSuite) TestSubscribeNoKeywords() {
	alice := &Subscription{User: "alice", Channel: "#chat2", FeedId: s.feed.Id}
	s.Nil(s.repository.AddSub(alice))

	results, _ := s.processor.Process()

	s.Len(results, 1)
	assertNotificationCorrect(s.T(), results[0], alice, s.feed)
	s.True(alice.HasSeen(*s.item))
}

func TestRunProcessorSuite(t *testing.T) {
	suite.Run(t, new(FeedProcessorSuite))
}

func TestProcessor_Process_returns_the_expected_notifications(t *testing.T) {
	item := &ParsedItem{Title: "bar", GUID: "1"}
	parsed := &ParsedFeed{Items: []*ParsedItem{item}}
	processor := NewFeedProcessor(NewInMemRepo(), &StubParser{Parsed: parsed})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id, Ignore: "dub"}
	processor.repository.AddSub(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(james)

	results, _ := processor.Process()

	assert.Len(t, results, 2)

	// TODO: lol fix race condition
	assertNotificationCorrect(t, results[0], alice, feed)
	assert.True(t, alice.HasSeen(*item))
	assertNotificationCorrect(t, results[1], james, feed)
	assert.True(t, james.HasSeen(*item))
}

func TestProcessor_Process_returns_grouped_notifications_by_channel_and_item(t *testing.T) {
	parsed := &ParsedFeed{Items: []*ParsedItem{{Title: "bar", GUID: "1"}}}
	processor := NewFeedProcessor(NewInMemRepo(), &StubParser{Parsed: parsed})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(james)

	results, _ := processor.Process()

	assert.Len(t, results, 1)

	assertNotificationCorrect(t, results[0], alice, feed)
	assert.Len(t, results[0].Users, 2, "notification should have alice and james")
}

func TestProcessor_Process_returns_empty_when_no_keywords_found(t *testing.T) {
	p := &ParsedFeed{Items: []*ParsedItem{{Title: "foo"}}}
	processor := NewFeedProcessor(NewInMemRepo(), &StubParser{Parsed: p})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)
	processor.repository.AddSub(&Subscription{User: "james", Channel: "#chat", Keywords: "baz", FeedId: feed.Id})

	notes, _ := processor.Process()

	assert.Empty(t, notes)
}

func TestProcessor_Process_ignores_seen_items(t *testing.T) {
	item := &ParsedItem{Title: "foo"}
	p := &ParsedFeed{Items: []*ParsedItem{item}}
	processor := NewFeedProcessor(NewInMemRepo(), &StubParser{Parsed: p})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)
	sub := &Subscription{User: "james", Channel: "#chat", Keywords: "foo", FeedId: feed.Id}
	sub.Remember(*item)
	processor.repository.AddSub(sub)

	notes, _ := processor.Process()

	assert.Empty(t, notes)
}

func TestProcessor_Process_rate_limits_notifications_per_channel(t *testing.T) {
	item := &ParsedItem{Title: "bar", GUID: "1"}
	parsed := &ParsedFeed{Items: []*ParsedItem{
		item,
		{Title: "bar", GUID: "2"},
		{Title: "bar", GUID: "3"},
		{Title: "bar", GUID: "4"},
	}}
	processor := NewFeedProcessor(NewInMemRepo(), &StubParser{Parsed: parsed})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(alice)

	results, _ := processor.Process()

	assert.Len(t, results, 3, "limiter should have only allowed 3 notifications")
}

func TestProcessor_Process_returns_empty_when_keywords_found_but_has_ignore_words(t *testing.T) {
	p := &ParsedFeed{Items: []*ParsedItem{{Title: "foo bar", GUID: "1"}}}
	processor := NewFeedProcessor(NewInMemRepo(), &StubParser{Parsed: p})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)
	processor.repository.AddSub(&Subscription{
		User: "james", Channel: "#chat",
		Keywords: "foo",
		Ignore:   "bar",
		FeedId:   feed.Id})

	notes, _ := processor.Process()

	assert.Empty(t, notes)
}

func TestProcessReturnsRepositoryError(t *testing.T) {
	repo := NewInMemRepo()
	repo.ForceError(errors.New("forced"), 0)
	processor := NewFeedProcessor(repo, &StubParser{})

	notes, err := processor.Process()
	assert.Nil(t, notes)
	assert.Error(t, err)
}

func assertNotificationCorrect(t *testing.T, n *Notification, sub *Subscription, feed *Feed) {
	assert.Equal(t, n.Channel, sub.Channel)
	assert.Contains(t, n.Users, sub.User)
	assert.Equal(t, n.Feed.Id, feed.Id)
}
