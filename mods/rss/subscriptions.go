package rss

import (
	"strings"
)

type Subscription struct {
	Id       uint64
	FeedId   uint64
	User     string
	Keywords string
	Channel  string
	Feed     *Feed
	Seen     map[string]bool // [guid]item
}

func (s *Subscription) HasSeen(item Item) bool {
	return false
}

func (sub *Subscription) KeywordsSlice() []string {
	return strings.Split(sub.Keywords, ",")
}

type Subscriptions interface {
	Add(*Subscription) error
	ByFeedId(id uint64) ([]*Subscription, error)
}

type InMemSubs struct {
	items []*Subscription
}

func (m *InMemSubs) Add(s *Subscription) error {
	m.items = append(m.items, s)
	return nil
}

func (m *InMemSubs) ByFeedId(id uint64) ([]*Subscription, error) {
	var found []*Subscription

	for _, s := range m.items {
		if s.FeedId == id {
			found = append(found, s)
		}
	}

	return found, nil
}
