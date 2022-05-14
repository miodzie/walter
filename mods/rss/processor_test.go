package rss

import (
	"testing"
)

func TestProcessor_Handle(t *testing.T) {
	parsed := &ParsedFeed{Title: "foo"}
	sut := &Processor{
		feeds:  &InMemFeeds{},
		subs:   &InMemSubs{},
		parser: &NulledParser{parsed: parsed},
	}
	feed := &Feed{Id: 1}
	sut.feeds.Add(feed)
	fooSub := &Subscription{
		User:     "adam",
		Channel:  "#chat",
		Keywords: "foo",
		FeedId:   feed.Id,
	}
	sut.subs.Add(fooSub)
	barSub := &Subscription{
		User:     "alice",
		Channel:  "#chat2",
		Keywords: "bar",
		FeedId:   feed.Id,
	}
	sut.subs.Add(barSub)

	// Act
	notifs, _ := sut.Handle()

	// Assert
	if len(notifs) == 0 {
		t.Error("notes is empty")
	}

	fooNotif := notifs[0]
	checkSub(t, fooNotif, fooSub, feed)

	if len(notifs) > 1 {
		t.Error("there should only be a notification for fooSub")
	}
}

func checkSub(t *testing.T, n *Notification, sub *Subscription, feed *Feed) {
	if n.Channel != sub.Channel {
		t.Error("expected notification not found")
	}
	if n.Users[0] != sub.User {
		t.Error("expected notification not found")
	}
	if n.Feed.Id != feed.Id {
		t.Error("expected notification not found")
	}
}
