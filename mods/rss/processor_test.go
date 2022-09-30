package rss

import (
	"testing"
)

func TestProcessor_Process_returns_the_expected_notifications(t *testing.T) {
	item := &Item{Title: "bar", GUID: "1"}
	parsed := &ParsedFeed{Items: []*Item{item}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: parsed})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	processor.repository.AddSub(james)

	// Act
	results, _ := processor.Process()

	// Assert
	if len(results) != 2 {
		t.Error("unexpected results")
	}

	checkNotif(t, results[0], alice, feed)
	if !alice.HasSeen(*item) {
		t.Fail()
	}
	checkNotif(t, results[1], james, feed)
	if !james.HasSeen(*item) {
		t.Fail()
	}
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
		t.Error("unexpected results")
	}

	checkNotif(t, results[0], alice, feed)
	if len(results[0].Users) != 2 {
		t.Error("notification should have alice and james")
	}
}

func TestProcessor_Process_returns_empty_when_no_keywords_found(t *testing.T) {
	p := &ParsedFeed{Items: []*Item{{Title: "foo"}}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: p})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)
	processor.repository.AddSub(&Subscription{User: "james", Channel: "#chat", Keywords: "baz", FeedId: feed.Id})

	notifs, _ := processor.Process()

	if len(notifs) != 0 {
		t.Fail()
	}
}

func TestProcessor_Process_ignores_seen_items(t *testing.T) {
	item := &Item{Title: "foo"}
	p := &ParsedFeed{Items: []*Item{item}}
	processor := NewProcessor(NewInMemRepo(), &NullParser{Parsed: p})
	feed := &Feed{Id: 1}
	processor.repository.AddFeed(feed)
	sub := &Subscription{User: "james", Channel: "#chat", Keywords: "foo", FeedId: feed.Id}
	sub.See(*item)
	processor.repository.AddSub(sub)

	// Act
	notifs, _ := processor.Process()

	// Assert
	if len(notifs) != 0 {
		t.Fail()
	}
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
	if len(results) != 3 {
		t.Logf("len(results)=%d, expected 3", len(results))
		t.Error("limiter should have only allowed 3 notifications")
	}
}

func checkNotif(t *testing.T, n *Notification, sub *Subscription, feed *Feed) {
	if n.Channel != sub.Channel {
		t.Error("unexpected notification.Target")
	}
	if n.Users[0] != sub.User {
		t.Error("unexpected notification.Users")
	}
	if n.Feed.Id != feed.Id {
		t.Error("unexpected feed.Id")
	}
}
