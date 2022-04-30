package rss

import (
	"strings"
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

type Subscription struct {
	Id       uint64
	FeedId   uint64
	User     string
	Keywords string
	Channel  string
	Feed     *Feed
	Seen     map[string]interface{} // [guid]item
}

func (sub *Subscription) KeywordsSlice() []string {
	return strings.Split(sub.Keywords, ",")
}

type SubscriptionRepository interface {
	Save(*Subscription) error
	GetByFeedId(id uint64) ([]Subscription, error)
}
