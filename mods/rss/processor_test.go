package rss

import (
	"testing"
)

func TestProcessor_Handle(t *testing.T) {
	sut := &Processor{
		feeds:  &InMemFeeds{},
		subs:   &InMemSubs{},
		parser: &NulledParser{},
	}
	feed := &Feed{Id: 1}
	sub := &Subscription{Channel: "#chat"}

	notes, err := sut.Handle()
	if err != nil {
		t.Error(err)
	}

	if len(notes) == 0 {
		t.Error("notes is empty")
	}

	n := notes[0]
	if n.Channel != sub.Channel {
		t.Error("expected notification not found")
	}
	if n.Users[0] != "adam" {
		t.Error("expected notification not found")
	}
	if n.Feed.Id != feed.Id {
		t.Error("expected notification not found")
	}
}
