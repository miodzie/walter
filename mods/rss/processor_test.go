package rss

import (
	"testing"
)

func TestProcessor_Handle_returns_the_expected_notifications(t *testing.T) {
	parsed := &ParsedFeed{Items: []*Item{{Title: "bar", GUID: "1"}}}
	sut := &Processor{parser: &NulledParser{parsed: parsed}, feeds: &InMemFeeds{}, subs: &InMemSubs{}}
	feed := &Feed{Id: 1}
	sut.feeds.Add(feed)

	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id}
	sut.subs.Add(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	sut.subs.Add(james)

	// Act
	results, _ := sut.Handle()

	// Assert
	if len(results) != 2 {
		t.Error("unexpected results")
	}

	checkNotif(t, results[0], alice, feed)
	checkNotif(t, results[1], james, feed)
}

func TestProcessor_Handle_returns_grouped_notifications_by_channel_and_item(t *testing.T) {
	parsed := &ParsedFeed{Items: []*Item{{Title: "bar", GUID: "1"}}}
	sut := &Processor{parser: &NulledParser{parsed: parsed}, feeds: &InMemFeeds{}, subs: &InMemSubs{}}
	feed := &Feed{Id: 1}
	sut.feeds.Add(feed)

	alice := &Subscription{User: "alice", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	sut.subs.Add(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	sut.subs.Add(james)

	// Act
	results, _ := sut.Handle()

	// Assert
	if len(results) != 1 {
		t.Error("unexpected results")
	}

	checkNotif(t, results[0], alice, feed)
	if len(results[0].Users) != 2 {
		t.Error("notification should have alice and james")
	}
}

func TestProcessor_Handle_returns_empty_when_no_keywords_found(t *testing.T) {
	p := &ParsedFeed{Items: []*Item{{Title: "foo"}}}
	sut := &Processor{parser: &NulledParser{parsed: p}, feeds: &InMemFeeds{}, subs: &InMemSubs{}}
	feed := &Feed{Id: 1}
	sut.feeds.Add(feed)
	sut.subs.Add(&Subscription{User: "james", Channel: "#chat", Keywords: "baz", FeedId: feed.Id})

	notifs, _ := sut.Handle()

	if len(notifs) != 0 {
		t.Fail()
	}
}

func checkNotif(t *testing.T, n *Notification, sub *Subscription, feed *Feed) {
	if n.Channel != sub.Channel {
		t.Error("unexpected notification.Channel")
	}
	if n.Users[0] != sub.User {
		t.Error("unexpected notification.Users")
	}
	if n.Feed.Id != feed.Id {
		t.Error("unexpected feed.Id")
	}
}
