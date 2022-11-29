// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"github.com/miodzie/walter/log"
	"sync"
)

type Processor struct {
	repository Repository
	parser     Parser
	// Max notifications sent per channel per process.
	ChannelLimit int
	cache        *notificationCache
	sync.Mutex
}

func NewProcessor(repo Repository, parser Parser) *Processor {
	return &Processor{
		repository:   repo,
		parser:       parser,
		ChannelLimit: 3,
	}
}

func (p *Processor) Process() ([]*Notification, error) {
	p.Lock()
	defer p.Unlock()
	p.cache = newCache()
	var notifications []*Notification
	feeds, _ := p.repository.Feeds()
	for _, feed := range feeds {
		log.Debug("feed: " + feed.Name)
		parsedFeed, err := p.parser.ParseURL(feed.Url)
		if err != nil {
			return notifications, err
		}

		subs, err := p.repository.Subs(SubSearchOpt{FeedId: feed.Id})
		if err != nil {
			return notifications, err
		}

		for _, subscription := range subs {
			newNotes := p.processSubscription(parsedFeed, subscription)
			notifications = append(notifications, newNotes...)
		}
	}

	return notifications, nil
}

func (p *Processor) processSubscription(parsedFeed *ParsedFeed, subscription *Subscription) []*Notification {
	var notifications []*Notification
	// Fetch all items with a subscriptions keywords.
	for _, item := range parsedFeed.ItemsWithKeywords(subscription.KeywordsSlice()) {
		// Ignore
		if p.cache.ChannelLimitReached(subscription.Channel, p.ChannelLimit) {
			continue
		}
		if subscription.HasSeen(*item) {
			continue
		}
		if subscription.Ignore != "" &&
			item.HasKeywords(subscription.IgnoreSlice()) {
			continue
		}

		key := p.cache.makeKey(item, subscription)
		notification := p.cache.GetNotification(key)
		if !p.cache.HasNotification(key) {
			notification = &Notification{
				Channel: subscription.Channel,
				Feed:    *subscription.Feed,
				Item:    *item,
			}
			notifications = append(notifications, notification)
			p.cache.PutNotification(key, notification)
		}

		notification.Users = append(notification.Users, subscription.User)
		subscription.See(*item)
	}
	err := p.repository.UpdateSub(subscription)
	if err != nil {
		log.Error(err)
	}

	return notifications
}

func newCache() *notificationCache {
	return &notificationCache{
		channelAmount:     make(map[string]int),
		seenNotifications: map[string]*Notification{},
	}
}

type notificationCache struct {
	channelAmount     map[string]int
	seenNotifications map[string]*Notification
}

func (c *notificationCache) HasNotification(key string) bool {
	_, exists := c.seenNotifications[key]
	return exists
}

func (c *notificationCache) GetNotification(key string) *Notification {
	return c.seenNotifications[key]
}

func (c *notificationCache) PutNotification(key string, notification *Notification) {
	if _, ok := c.channelAmount[notification.Channel]; !ok {
		c.channelAmount[notification.Channel] = 0
	}
	c.channelAmount[notification.Channel] += 1
	c.seenNotifications[key] = notification
}

func (c *notificationCache) makeKey(item *Item, sub *Subscription) string {
	return item.GUID + sub.Channel
}

func (c *notificationCache) ChannelLimitReached(channelId string, limit int) bool {
	if amt, ok := c.channelAmount[channelId]; ok {
		return amt >= limit
	}
	return false
}
