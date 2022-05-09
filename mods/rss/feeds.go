package rss

import (
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

type FeedRepository interface {
	All() ([]Feed, error)
	Add(*Feed) error
	ByName(name string) (Feed, error)
}

// Parser downloads a Feed.Url and translates it to a ParsedFeed to
// be checked by a Subscription.
type Parser interface {
	Parse(Feed) ParsedFeed
}

type ParsedFeed struct {
	Title       string
	Description string
	Link        string
	FeedLink    string
	Updated     string
	Published   string
	Items       []*Item
	Custom      map[string]string
}

type Item struct {
	Title       string
	Description string
	Content     string
	GUID        string
	Custom      map[string]string
}
