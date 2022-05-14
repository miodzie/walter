package rss

import (
	"fmt"
	"testing"
)

func TestProcessor_Handle(t *testing.T) {
	parsed := &ParsedFeed{Title: "foo", Description: "bar"}
	sut := &Processor{
		feeds:  &InMemFeeds{},
		subs:   &InMemSubs{},
		parser: &NulledParser{parsed: parsed},
	}
	feed := &Feed{Id: 1}
	sut.feeds.Add(feed)
	adam := &Subscription{User: "adam", Channel: "#chat", Keywords: "foo", FeedId: feed.Id}
	sut.subs.Add(adam)
	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id}
	sut.subs.Add(alice)
	sut.subs.Add(&Subscription{User: "james", Channel: "#chat2", Keywords: "baz", FeedId: feed.Id})
	dakota := &Subscription{User: "dakota", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	sut.subs.Add(dakota)

	// Act
	notifs, _ := sut.Handle()

	// Assert
	if len(notifs) == 0 {
		t.Error("notes is empty")
	}

	// notifs[0] should have Users: adam and dakota
	fooNotif := notifs[0]
	for _, i := range notifs {
		fmt.Printf("%v\n", i)
	}
	if len(fooNotif.Users) != 2 {
		t.Errorf("fooNotif should have %s and %s", adam.User, dakota.User)
	}
	checkSub(t, fooNotif, adam, feed)
	checkSub(t, notifs[1], alice, feed)
	if len(notifs) > 2 {
		t.Error("there should only be a notification for fooSub and barSub")
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
