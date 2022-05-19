package rss

import (
	"testing"
)

func TestSubscription_Seen(t *testing.T) {
	sub := &Subscription{SeenItems: make(map[string]bool)}
	item := Item{GUID: "1234"}

	sub.See(item)

	if _, ok := sub.SeenItems[item.GUID]; !ok {
		t.Fail()
	}
}

func TestSubscription_HasSeen(t *testing.T) {
	sub := &Subscription{SeenItems: make(map[string]bool)}
	item := Item{GUID: "1234"}

	if sub.HasSeen(item) {
		t.Error("sub should not have seen item")
	}

	sub.See(item)

	if sub.HasSeen(item) == false {
		t.Error("sub should have seen item")
	}
}
