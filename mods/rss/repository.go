// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"errors"
	"math/rand"
)

// TODO: Add more Repository errors.

var FeedNotFoundError = errors.New("userFeed not found")

// Repository is staying here because I can.
type Repository interface {
	AddFeed(*UserFeed) error
	RemoveFeed(name string) error
	Feeds() ([]*UserFeed, error)
	FeedByName(name string) (*UserFeed, error)

	AddSub(*Subscription) error
	UpdateSub(*Subscription) error
	RemoveSub(*Subscription) error
	// Subs TODO: Refactor into broken down methods again.
	Subs(params SearchParams) ([]*Subscription, error)
}

type SearchParams struct {
	FeedId   uint64
	User     string
	FeedName string
	Channel  string
}

type InMemRepository struct {
	feeds            []*UserFeed
	subs             map[uint64]*Subscription
	delayForcedErrBy int
	ForcedErr        error
	tmpIgnoreErr     bool
}

func NewInMemRepo() *InMemRepository {
	return &InMemRepository{subs: make(map[uint64]*Subscription)}
}

// ForceError sets an error to be returned on the next called method.
// Used for forcing errors in testing.
func (r *InMemRepository) ForceError(err error, delay int) {
	r.ForcedErr = err
	r.delayForcedErrBy = delay
}

func (r *InMemRepository) ignoreErr() {
	r.tmpIgnoreErr = true
}

func (r *InMemRepository) popForcedErr() error {
	if r.tmpIgnoreErr {
		r.tmpIgnoreErr = false
		return nil
	}
	if r.delayForcedErrBy != 0 {
		r.delayForcedErrBy -= 1
		return nil
	}
	defer func() { r.ForcedErr = nil }()
	return r.ForcedErr
}

func (r *InMemRepository) Feeds() ([]*UserFeed, error) {
	return r.feeds, r.popForcedErr()
}

func (r *InMemRepository) AddFeed(feed *UserFeed) error {
	r.feeds = append(r.feeds, feed)
	return r.popForcedErr()
}

func (r *InMemRepository) RemoveFeed(name string) error {
	for i, f := range r.feeds {
		if f.Name == name {
			r.feeds[i] = r.feeds[len(r.feeds)-1]
			r.feeds = r.feeds[:len(r.feeds)-1]
			break
		}
	}
	return nil
}

func (r *InMemRepository) FeedByName(name string) (*UserFeed, error) {
	for _, c := range r.feeds {
		if c.Name == name {
			return c, r.popForcedErr()
		}
	}
	return &UserFeed{}, FeedNotFoundError
}

func (r *InMemRepository) AddSub(s *Subscription) error {
	if s.Id == 0 {
		s.Id = rand.Uint64()
	}
	r.subs[s.Id] = s
	return r.popForcedErr()
}

func (r *InMemRepository) UpdateSub(s *Subscription) error {
	r.subs[s.Id] = s
	return r.popForcedErr()
}

func (r *InMemRepository) RemoveSub(subscription *Subscription) error {
	delete(r.subs, subscription.Id)
	return r.popForcedErr()
}

func (r *InMemRepository) Subs(search SearchParams) ([]*Subscription, error) {
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
				r.ignoreErr()
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
			r.ignoreErr()
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

	return subs, r.popForcedErr()
}
