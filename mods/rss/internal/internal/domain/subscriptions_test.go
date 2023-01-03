// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubscription_Seen(t *testing.T) {
	sub := &Subscription{SeenItems: make(map[string]bool)}
	item := Item{GUID: "1234"}

	sub.Remember(item)
	sub.Remember(Item{GUID: "1"})
	sub.Remember(Item{GUID: "1"})

	_, exists := sub.SeenItems[item.GUID]
	assert.True(t, exists)
	assert.Equal(t, "1234,1", sub.Seen)
}

func TestSubscription_HasSeen(t *testing.T) {
	sub := &Subscription{SeenItems: make(map[string]bool)}
	item := Item{GUID: "1234"}
	assert.False(t, sub.HasSeen(item))

	sub.Remember(item)

	assert.True(t, sub.HasSeen(item))
}
