// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package domain

import (
	"strings"
)

type Subscriptions interface {
	Add(*Subscription) error
	Update(*Subscription) error
	ByFeedId(id uint64) ([]*Subscription, error)
}

type Subscription struct {
	Id        uint64
	FeedId    uint64
	User      string
	Keywords  string
	Ignore    string
	Channel   string
	Feed      *Feed
	Seen      string          // Item.GUID comma separated
	SeenItems map[string]bool // [guid]bool
}

func (s *Subscription) Remember(item Item) {
	s.makeSeenMap()
	if _, seen := s.SeenItems[item.GUID]; !seen {
		s.SeenItems[item.GUID] = true
		s.Seen = ""
		var keys []string
		for k := range s.SeenItems {
			keys = append(keys, k)
		}
		s.Seen = strings.Join(keys, ",")
	}
}

func (s *Subscription) HasSeen(item Item) bool {
	s.makeSeenMap()
	_, seen := s.SeenItems[item.GUID]
	return seen
}

func (s *Subscription) ShouldIgnore(item Item) bool {
	return s.HasSeen(item) ||
		(item.HasKeywords(s.IgnoreWords()) && s.Ignore != "")
}

func (s *Subscription) KeyWords() []string {
	return strings.Split(s.Keywords, ",")
}

func (s *Subscription) IgnoreWords() []string {
	return strings.Split(s.Ignore, ",")
}

func (s *Subscription) makeSeenMap() {
	if s.SeenItems == nil {
		s.SeenItems = make(map[string]bool)
		for _, i := range strings.Split(s.Seen, ",") {
			s.SeenItems[i] = true
		}
	}
}
