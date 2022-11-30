// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessor_Process_returns_the_expected_notifications(t *testing.T) {
	item := &Item{Title: "bar", GUID: "1"}
	parsed := &ParsedFeed{Items: []*Item{item}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: parsed})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id, Ignore: "dub"}
	processor.repository.AddSub(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(james)

	// Act
	results, _ := processor.Process()

	// Assert
	assert.Len(t, results, 2)

	assertNotificationCorrect(t, results[0], alice, feed)
	assert.True(t, alice.HasSeen(*item))
	assertNotificationCorrect(t, results[1], james, feed)
	assert.True(t, james.HasSeen(*item))
}

func TestProcessor_Process_returns_grouped_notifications_by_channel_and_item(t *testing.T) {
	parsed := &ParsedFeed{Items: []*Item{{Title: "bar", GUID: "1"}}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: parsed})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(james)

	// Act
	results, _ := processor.Process()

	// Assert
	if len(results) != 1 {
		t.Errorf("expected len(results)=1, got %d instead", len(results))
	}

	assertNotificationCorrect(t, results[0], alice, feed)
	assert.Len(t, results[0].Users, 2, "notification should have alice and james")
}

func TestProcessor_Process_returns_empty_when_no_keywords_found(t *testing.T) {
	p := &ParsedFeed{Items: []*Item{{Title: "foo"}}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: p})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)
	processor.repository.AddSub(&Subscription{User: "james", Channel: "#chat", Keywords: "baz", FeedId: feed.Id})

	notes, _ := processor.Process()

	assert.Empty(t, notes)
}

func TestProcessor_Process_ignores_seen_items(t *testing.T) {
	item := &Item{Title: "foo"}
	p := &ParsedFeed{Items: []*Item{item}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: p})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)
	sub := &Subscription{User: "james", Channel: "#chat", Keywords: "foo", FeedId: feed.Id}
	sub.Remember(*item)
	processor.repository.AddSub(sub)

	// Act
	notes, _ := processor.Process()

	// Assert
	assert.Empty(t, notes)
}

func TestProcessor_Process_rate_limits_notifications_per_channel(t *testing.T) {
	item := &Item{Title: "bar", GUID: "1"}
	parsed := &ParsedFeed{Items: []*Item{
		item,
		{Title: "bar", GUID: "2"},
		{Title: "bar", GUID: "3"},
		{Title: "bar", GUID: "4"},
	}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: parsed})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(alice)

	// Act
	results, _ := processor.Process()

	// Assert
	assert.Len(t, results, 3, "limiter should have only allowed 3 notifications")
}

func TestProcessor_Process_returns_empty_when_keywords_found_but_has_ignore_words(t *testing.T) {
	p := &ParsedFeed{Items: []*Item{{Title: "foo bar", GUID: "1"}}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: p})
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

func assertNotificationCorrect(t *testing.T, n *Notification, sub *Subscription, feed *Feed) {
	assert.Equal(t, n.Channel, sub.Channel)
	assert.Equal(t, n.Users[0], sub.User)
	assert.Equal(t, n.Feed.Id, feed.Id)
}
