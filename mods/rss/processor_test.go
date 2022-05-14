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
	sut.feeds.Add(&Feed{})

	notes, err := sut.Handle()
	if err != nil {
		t.Error(err)
	}

	if len(notes) == 0 {
		t.Error("notes is empty")
	}
}
