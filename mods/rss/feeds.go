package rss

import (
	"time"
)

// Feeds are the allowed and available RSS feeds that users can subscribe
// to.
type Feed struct {
	Id            uint
    Name          string
	Url           string
	LastPublished time.Time
}

type FeedRepository interface {
	All() ([]Feed, error)
	Get(id string) (*Feed, error)
	GetByName(name string) (*Feed, error)
	Save(*Feed) error
	Delete(id string) error
}

type Subscription struct {
	Id       uint
	Feed     *Feed
	User     string
	Keywords []string
	Channel  string
	Seen     map[string]interface{} // [guid]item
}

type SubscriptionRepository interface {
	All() []Subscription
	Save(*Subscription) error
	Delete(subId string) error
	GetByFeedId(id uint) ([]Subscription, error)
}
