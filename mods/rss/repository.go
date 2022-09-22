package rss

import (
	"errors"
	"math/rand"
)

// Repository
// TODO: Add defined errors here for Repositories to use.
type Repository interface {
	Feeds() ([]*Feed, error)
	AddFeed(*Feed) error
	FeedByName(name string) (*Feed, error)

	AddSub(*Subscription) error
	UpdateSub(*Subscription) error
	RemoveSub(*Subscription) error
	Subs(search SubSearchOpt) ([]*Subscription, error)
}

type SubSearchOpt struct {
	FeedId   uint64
	User     string
	FeedName string
	Channel  string
}

type InMemRepository struct {
	feeds []*Feed
	subs  map[uint64]*Subscription
}

func (m *InMemRepository) Subs(search SubSearchOpt) ([]*Subscription, error) {
	var subs []*Subscription
	var constraints = []func(sub *Subscription) bool{
		func(sub *Subscription) bool {
			if search.User != "" {
				return sub.User == search.User
			}
			return true
		},
		func(sub *Subscription) bool {
			if search.Channel != "" {
				return sub.Channel == search.Channel
			}
			return true
		},
		func(sub *Subscription) bool {
			if search.FeedId != 0 {
				return sub.FeedId == search.FeedId
			}
			return true
		},
		func(sub *Subscription) bool {
			if search.FeedId != 0 {
				feed, err := m.FeedByName(search.FeedName)
				if err != nil {
					return false
				}
				return feed.Name == search.FeedName
			}
			return true
		},
	}

	matches := func(sub *Subscription) bool {
		for _, check := range constraints {
			if !check(sub) {
				return false
			}
		}
		return true
	}
	for _, sub := range m.subs {
		if matches(sub) {
			subs = append(subs, sub)
		}
	}

	return subs, nil
}

func NewInMemRepo() *InMemRepository {
	return &InMemRepository{subs: make(map[uint64]*Subscription)}
}

func (m *InMemRepository) Feeds() ([]*Feed, error) {
	return m.feeds, nil
}

func (m *InMemRepository) AddFeed(feed *Feed) error {
	m.feeds = append(m.feeds, feed)
	return nil
}

func (m *InMemRepository) FeedByName(name string) (*Feed, error) {
	for _, c := range m.feeds {
		if c.Name == name {
			return c, nil
		}
	}
	return &Feed{}, errors.New("feed not found")
}

func (m *InMemRepository) AddSub(s *Subscription) error {
	if s.Id == 0 {
		s.Id = rand.Uint64()
	}
	m.subs[s.Id] = s
	return nil
}

func (m *InMemRepository) UpdateSub(s *Subscription) error {
	m.subs[s.Id] = s
	return nil
}

func (m *InMemRepository) RemoveSub(subscription *Subscription) error {
	delete(m.subs, subscription.Id)
	return nil
}
