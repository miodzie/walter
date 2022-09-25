package rss

import (
	"errors"
	"math/rand"
)

// Repository
// TODO: AddFeed defined errors here for Repositories to use.
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

func (r *InMemRepository) Subs(search SubSearchOpt) ([]*Subscription, error) {
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
			if search.FeedName != "" {
				feed, err := r.FeedByName(search.FeedName)
				if err != nil {
					return false
				}
				return sub.FeedId == feed.Id
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
	for _, sub := range r.subs {
		if matches(sub) {
			feeds, _ := r.Feeds()
			for _, feed := range feeds {
				if feed.Id == sub.FeedId {
					sub.Feed = feed
					break
				}
			}
			subs = append(subs, sub)
		}
	}

	return subs, nil
}

func NewInMemRepo() *InMemRepository {
	return &InMemRepository{subs: make(map[uint64]*Subscription)}
}

func (r *InMemRepository) Feeds() ([]*Feed, error) {
	return r.feeds, nil
}

func (r *InMemRepository) AddFeed(feed *Feed) error {
	r.feeds = append(r.feeds, feed)
	return nil
}

func (r *InMemRepository) FeedByName(name string) (*Feed, error) {
	for _, c := range r.feeds {
		if c.Name == name {
			return c, nil
		}
	}
	return &Feed{}, errors.New("feed not found")
}

func (r *InMemRepository) AddSub(s *Subscription) error {
	if s.Id == 0 {
		s.Id = rand.Uint64()
	}
	r.subs[s.Id] = s
	return nil
}

func (r *InMemRepository) UpdateSub(s *Subscription) error {
	r.subs[s.Id] = s
	return nil
}

func (r *InMemRepository) RemoveSub(subscription *Subscription) error {
	delete(r.subs, subscription.Id)
	return nil
}
