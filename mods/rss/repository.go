package rss

import (
	"errors"
	"fmt"
	"math/rand"
)

type Repository interface {
	Feeds() ([]*Feed, error)
	AddFeed(*Feed) error
	FeedByName(name string) (*Feed, error)

	AddSub(*Subscription) error
	UpdateSub(*Subscription) error
	RemoveSub(*Subscription) error
	SubByUserFeedNameChannel(user, feedName, channel string) (*Subscription, error)
	// TODO: Refactor to use Feed Name, Id is a database concept.
	// Separate the all mighty database from the domain!
	// e.g. SubsByFeed(*Feed) .. Let the repository accept the whole feed,
	//and let them figure it out, or,
	//just use the Feed.Name since it should be unique anyways.
	SubsByFeedId(id uint64) ([]*Subscription, error)
}

type InMemRepository struct {
	feeds []*Feed
	subs  map[uint64]*Subscription
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

func (m *InMemRepository) SubByUserFeedNameChannel(user, feedName, channel string) (*Subscription, error) {
	var feed *Feed
	for _, f := range m.feeds {
		if f.Name == feedName {
			feed = f
		}
	}
	if feed == nil {
		return nil, errors.New(
			fmt.Sprintf("Could not locate Feed with name: %s",
				feedName,
			))
	}

	for _, sub := range m.subs {
		if sub.User == user && sub.FeedId == feed.Id && sub.Channel == channel {
			return sub, nil
		}
	}
	return nil, errors.New("subscription not found")
}

func (m *InMemRepository) SubsByFeedId(id uint64) ([]*Subscription, error) {
	var found []*Subscription

	for _, s := range m.subs {
		if s.FeedId == id {
			found = append(found, s)
		}
	}

	return found, nil
}
