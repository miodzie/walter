package rss

import (
	"strings"
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
	Parse(Feed) (ParsedFeed, error)
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
	Raw         string
}

func (feed *ParsedFeed) ItemsWithKeywords(keywords []string) []*Item {
	var items []*Item
	for _, i := range feed.Items {
		if i.HasKeywords(keywords) {
			items = append(items, i)
		}
	}
	return items
}

func (feed *ParsedFeed) HasKeywords(keywords []string) bool {
	// for _, k := range keywords {
	// 	if strings.Contains(feed.Raw, k) {
	// 		return true
	// 	}
	// }
	// return false
	for _, keyword := range keywords {
		checks := []bool{
			strings.Contains(feed.Title, keyword),
			strings.Contains(feed.Description, keyword),
		}
		if anyTrue(checks) {
			return true
		}
	}

	for _, item := range feed.Items {
		if item.HasKeywords(keywords) {
			return true
		}
	}

	return false
}

type Item struct {
	Title       string
	Description string
	Content     string
	GUID        string
	Custom      map[string]string
}

func (i *Item) HasKeywords(keywords []string) bool {
	for _, keyword := range keywords {
		checks := []bool{
			strings.Contains(i.Title, keyword),
			strings.Contains(i.Description, keyword),
			strings.Contains(i.Content, keyword),
		}
		if anyTrue(checks) {
			return true
		}
	}

	return false
}

func anyTrue(checks []bool) bool {
	for _, found := range checks {
		if found {
			return true
		}
	}

	return false
}
