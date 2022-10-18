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
	Ignore    string
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

func (s *Subscription) KeywordsSlice() []string {
	return strings.Split(s.Keywords, ",")
}

func (s *Subscription) IgnoreSlice() []string {
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
