package rss

import (
	"errors"
	"math/rand"
)

type Repository interface {
	AllFeeds() ([]*Feed, error)
	AddFeed(*Feed) error
	FeedByName(name string) (*Feed, error)

	AddSub(*Subscription) error
	UpdateSub(*Subscription) error
	SubByFeedId(id uint64) ([]*Subscription, error)
}

type InMemRepository struct {
	feeds []*Feed
	subs  map[uint64]*Subscription
}

func NewInMemRepo() *InMemRepository {
	return &InMemRepository{subs: make(map[uint64]*Subscription)}
}

func (f *InMemRepository) AllFeeds() ([]*Feed, error) {
	return f.feeds, nil
}

func (f *InMemRepository) AddFeed(feed *Feed) error {
	f.feeds = append(f.feeds, feed)
	return nil
}

func (f *InMemRepository) FeedByName(name string) (*Feed, error) {
	for _, c := range f.feeds {
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

func (m *InMemRepository) SubByFeedId(id uint64) ([]*Subscription, error) {
	var found []*Subscription

	for _, s := range m.subs {
		if s.FeedId == id {
			found = append(found, s)
		}
	}

	return found, nil
}
