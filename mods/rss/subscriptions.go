package rss

import "strings"


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
