package rss

import (
	"errors"
	"time"
)

// Feeds are the allowed and available web feeds that users can subscribe
// to.
type Feed struct {
	Id            uint64
	Name          string
	Url           string
	LastPublished time.Time
}

type Feeds interface {
	All() ([]*Feed, error)
	Add(*Feed) error
	ByName(name string) (*Feed, error)
}

type InMemFeeds struct {
	f []*Feed
}

func (f *InMemFeeds) All() ([]*Feed, error) {
	return f.f, nil
}

func (f *InMemFeeds) Add(feed *Feed) error {
	f.f = append(f.f, feed)
	return nil
}

func (f *InMemFeeds) ByName(name string) (*Feed, error) {
	for _, c := range f.f {
		if c.Name == name {
			return c, nil
		}
	}
	return &Feed{}, errors.New("feed not found")
}
