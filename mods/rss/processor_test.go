package rss

import (
	"testing"
)

func TestProcessor_Handle_returns_the_expected_notifications(t *testing.T) {
	item := &Item{Title: "bar", GUID: "1"}
	parsed := &ParsedFeed{Items: []*Item{item}}
	sut := &Processor{
		parser: &NulledParser{Parsed: parsed},
		repo:   NewInMemRepo(),
	}
	feed := &Feed{Id: 1}
	sut.repo.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat2", Keywords: "bar", FeedId: feed.Id}
	sut.repo.AddSub(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	sut.repo.AddSub(james)

	// Act
	results, _ := sut.Handle()

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

func TestProcessor_Handle_returns_grouped_notifications_by_channel_and_item(t *testing.T) {
	parsed := &ParsedFeed{Items: []*Item{{Title: "bar", GUID: "1"}}}
	sut := &Processor{
	    parser: &NulledParser{Parsed: parsed},
        repo: NewInMemRepo(),
	}
	feed := &Feed{Id: 1}
	sut.repo.AddFeed(feed)

	alice := &Subscription{User: "alice", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	sut.repo.AddSub(alice)
	james := &Subscription{User: "james", Channel: "#chat", Keywords: "bar", FeedId: feed.Id}
	sut.repo.AddSub(james)

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
	sut := &Processor{
	    parser: &NulledParser{Parsed: p},
        repo: NewInMemRepo(),
    }
	feed := &Feed{Id: 1}
	sut.repo.AddFeed(feed)
	sut.repo.AddSub(&Subscription{User: "james", Channel: "#chat", Keywords: "baz", FeedId: feed.Id})

	notifs, _ := sut.Handle()

	if len(notifs) != 0 {
		t.Fail()
	}
}

func TestProcessor_Handle_ignores_seen_items(t *testing.T) {
	item := &Item{Title: "foo"}
	p := &ParsedFeed{Items: []*Item{item}}
	sut := &Processor{parser: &NulledParser{Parsed: p}, repo: NewInMemRepo()}
	feed := &Feed{Id: 1}
	sut.repo.AddFeed(feed)
	sub := &Subscription{User: "james", Channel: "#chat", Keywords: "foo", FeedId: feed.Id}
	sub.See(*item)
	sut.repo.AddSub(sub)

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
