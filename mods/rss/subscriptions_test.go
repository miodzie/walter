// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"testing"
)

func TestSubscription_Seen(t *testing.T) {
	sub := &Subscription{SeenItems: make(map[string]bool)}
	item := Item{GUID: "1234"}

	sub.See(item)
	sub.See(Item{GUID: "1"})
	sub.See(Item{GUID: "1"})

	if _, ok := sub.SeenItems[item.GUID]; !ok {
		t.Fail()
	}
	if sub.Seen != "1234,1" {
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
