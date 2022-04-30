package rss

import (
	"time"
)

// Feeds are the allowed and available RSS feeds that users can subscribe
// to.
type Feed struct {
	Id            uint64
	Name          string
	Url           string
	LastPublished time.Time
}

type FeedRepository interface {
	All() ([]Feed, error)
	Save(*Feed) error
	GetByName(name string) (Feed, error)
}
