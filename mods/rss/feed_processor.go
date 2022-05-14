package rss

import (
	"fmt"
)

type CheckFeeds struct {
	feeds  Feeds
	subs   Subscriptions
	parser Parser
}

func NewCheckFeeds(f Feeds, s Subscriptions, parser Parser) *CheckFeeds {
	return &CheckFeeds{
		feeds:  f,
		subs:   s,
		parser: parser,
	}
}

func (c *CheckFeeds) Handle() ([]*Notification, error) {
	var notifications []*Notification
	feeds, err := c.feeds.All()
	if err != nil {
		return notifications, err
	}

	for _, feed := range feeds {
		// fetch and parse the feed:
		_, err := c.parser.Parse(feed.Url)

		// for each feed, pull its subscriptions.
		subs, err := c.subs.ByFeedId(feed.Id)
		if err != nil {
			return notifications, err
		}

		subsToNotify := make(map[string][]*Subscription)
		for _, sub := range subs {
			//  Fetch feed.
			//  Check feed for sub keywords.
			// Add sub to SubsToNotify slice if so.
			shouldNotify := true
			if shouldNotify {
				key := fmt.Sprintf("%d-%s", sub.FeedId, sub.Channel)
				subsToNotify[key] = append(subsToNotify[key], sub)
			}
		}

		// Sort by Channel and Feed, and group Users to notify into a
		// new struct.
		for _, subs := range subsToNotify {
			notification := &Notification{}
			for _, sub := range subs {
				notification.Channel = sub.Channel // TODO: Set only once
				notification.Users = append(notification.Users, sub.User)
			}
			notifications = append(notifications, notification)
		}
	}

	return notifications, nil
}
