// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type FeedProcessorSuite struct {
	processor  *FeedProcessor
	repository Repository
	feed       *UserFeed
	item       *Item
	suite.Suite
}

func (s *FeedProcessorSuite) SetupTest() {
	s.item = &Item{Title: "bar", GUID: "1"}
	parsed := &Feed{Items: []*Item{s.item}}
	s.repository = NewInMemRepo()
	s.processor = NewProcessor(s.repository, &StubParser{Parsed: parsed})
	s.feed = &UserFeed{Id: 1}
	_ = s.processor.repository.AddFeed(s.feed)
}

// TODO: This does tests nothing.
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
	item := &Item{Title: "bar", GUID: "1"}
	parsed := &Feed{Items: []*Item{item}}
	processor := NewProcessor(NewInMemRepo(), &StubParser{Parsed: parsed})
	feed := &UserFeed{Id: 1}
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
	parsed := &Feed{Items: []*Item{{Title: "bar", GUID: "1"}}}
	processor := NewProcessor(NewInMemRepo(), &StubParser{Parsed: parsed})
	feed := &UserFeed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(james)

	results, _ := processor.Process()

	assert.Len(t, results, 1)

	assertNotificationCorrect(t, results[0], alice, feed)
	assert.Len(t, results[0].User, 2, "notification should have alice and james")
}

func TestProcessor_Process_returns_empty_when_no_keywords_found(t *testing.T) {
	p := &Feed{Items: []*Item{{Title: "foo"}}}
	processor := NewProcessor(NewInMemRepo(), &StubParser{Parsed: p})
	feed := &UserFeed{Id: 1}
	processor.repository.AddFeed(feed)
	processor.repository.AddSub(&Subscription{User: "james", Channel: "#chat", Keywords: "baz", FeedId: feed.Id})

	notes, _ := processor.Process()

	assert.Empty(t, notes)
}

func TestProcessor_Process_ignores_seen_items(t *testing.T) {
	item := &Item{Title: "foo"}
	p := &Feed{Items: []*Item{item}}
	processor := NewProcessor(NewInMemRepo(), &StubParser{Parsed: p})
	feed := &UserFeed{Id: 1}
	processor.repository.AddFeed(feed)
	sub := &Subscription{User: "james", Channel: "#chat", Keywords: "foo", FeedId: feed.Id}
	sub.Remember(*item)
	processor.repository.AddSub(sub)

	notes, _ := processor.Process()

	assert.Empty(t, notes)
}

func TestProcessor_Process_rate_limits_notifications_per_channel(t *testing.T) {
	item := &Item{Title: "bar", GUID: "1"}
	parsed := &Feed{Items: []*Item{
		item,
		{Title: "bar", GUID: "2"},
		{Title: "bar", GUID: "3"},
		{Title: "bar", GUID: "4"},
	}}
	processor := NewProcessor(NewInMemRepo(), &StubParser{Parsed: parsed})
	feed := &UserFeed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(alice)

	results, _ := processor.Process()

	assert.Len(t, results, 3, "limiter should have only allowed 3 notifications")
}

func TestProcessor_Process_returns_empty_when_keywords_found_but_has_ignore_words(t *testing.T) {
	p := &Feed{Items: []*Item{{Title: "foo bar", GUID: "1"}}}
	processor := NewProcessor(NewInMemRepo(), &StubParser{Parsed: p})
	feed := &UserFeed{Id: 1}
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
	processor := NewProcessor(repo, &StubParser{})

	notes, err := processor.Process()
	assert.Nil(t, notes)
	assert.Error(t, err)
}

func assertNotificationCorrect(t *testing.T, n *Notification, sub *Subscription, feed *UserFeed) {
	assert.Equal(t, n.Channel, sub.Channel)
	assert.Contains(t, n.User, sub.User)
	assert.Equal(t, n.Feed.Id, feed.Id)
}
