package rss

import (
	"errors"
	"fmt"
	"strings"

	"github.com/miodzie/seras"
	"github.com/mmcdole/gofeed"
)

type Listener struct {
	FeedUrl string
	// TODO: Abstract gofeed to internal interface/struct.
	Seen map[string]*gofeed.Item
	// Keywords to check in the RSS Item.
	Keywords []string
	// Channel to notify users in.
	Channel string
	// Users to be notified.
	Notifees []string
}

func (listener *Listener) Process() ([]seras.Message, error) {
	msgs := []seras.Message{}
	fd := gofeed.NewParser()
	feed, err := fd.ParseURL(listener.FeedUrl)
	if err != nil {
		return msgs, err
	}
	for _, item := range feed.Items {
		if listener.ShouldNotify(item) {
			msgs = append(msgs, listener.CreateNotifcation(item))
		}
	}

	return msgs, nil
}

func (listener *Listener) CreateNotifcation(item *gofeed.Item) seras.Message {
	msg := seras.Message{Channel: listener.Channel}
    var images string
    for _, enclosure := range item.Enclosures {
        images += enclosure.URL
    }
	template := "Hot off the press!\n%s\n%s\nLink: %s\n\n%s"
	msg.Content = fmt.Sprintf(
		template,
		item.Title,
		images,
		item.Link,
		strings.Join(listener.Notifees, ", "),
	)

	return msg
}

func (listener *Listener) ShouldNotify(item *gofeed.Item) bool {
	if listener.HasSeen(item) {
		return false
	}
	listener.AddItem(item)

	// check Keywords
	title := strings.ToLower(item.Title)
	for _, keyword := range listener.Keywords {
		keyword = strings.ToLower(keyword)
		if strings.Contains(title, keyword) {
			return true
		}
	}

	return false
}

func (listener *Listener) HasSeen(item *gofeed.Item) bool {
	_, seen := listener.Seen[item.GUID]

	return seen
}

func (listener *Listener) AddItem(item *gofeed.Item) error {
	if listener.HasSeen(item) {
		return errors.New("item already exists in listener")
	}
	listener.Seen[item.GUID] = item

	return nil
}
