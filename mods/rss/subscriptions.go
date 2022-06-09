package rss

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
	Channel   string
	Feed      *Feed
	Seen      string          // Item.GUID comma separated
	SeenItems map[string]bool // [guid]bool
}

func (s *Subscription) See(item Item) {
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

func (sub *Subscription) KeywordsSlice() []string {
	return strings.Split(sub.Keywords, ",")
}

func (s *Subscription) makeSeenMap() {
	if s.SeenItems == nil {
		s.SeenItems = make(map[string]bool)
		for _, i := range strings.Split(s.Seen, ",") {
			s.SeenItems[i] = true
		}
	}
}
