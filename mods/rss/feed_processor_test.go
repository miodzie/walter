package rss

import (
	"testing"
)

func TestProcessor_Handle(t *testing.T) {
	sut := &CheckFeeds{
		feeds:  &InMemFeeds{},
		subs:   &InMemSubs{},
		parser: &NulledParser{},
	}
	sut.feeds.Add(&Feed{})

	sut.Handle()
}
