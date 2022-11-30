// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"github.com/miodzie/walter/log"
	"sync"
)

type Processor struct {
	// Max notifications sent per channel per process.
	ChannelLimit      int
	repository        Repository
	parser            Parser
	notificationCache notificationCache
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
	p.notificationCache = newNotificationCache()
	var notifications []*Notification
	feeds, _ := p.repository.Feeds()
	for _, feed := range feeds {
		parsedFeed, err := p.parser.ParseURL(feed.Url)
		if err != nil {
			return notifications, err
		}

		subs, err := p.repository.Subs(SearchParams{FeedId: feed.Id})
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
	for _, item := range parsedFeed.ItemsWithKeywords(subscription.KeywordsSlice()) {
		if p.shouldIgnore(subscription, item) {
			continue
		}

		notification, wasNew := p.getOrCreateNotification(subscription, item)
		notification.Users = append(notification.Users, subscription.User)
		if wasNew {
			notifications = append(notifications, notification)
		}
		subscription.Remember(*item)
	}

	err := p.repository.UpdateSub(subscription)
	if err != nil {
		log.Error(err)
	}

	return notifications
}

func (p *Processor) shouldIgnore(subscription *Subscription, item *Item) bool {
	if p.notificationCache.ChannelLimitReached(subscription.Channel, p.ChannelLimit) {
		return true
	}
	if subscription.HasSeen(*item) {
		return true
	}
	if subscription.Ignore != "" &&
		item.HasKeywords(subscription.IgnoreSlice()) {
		return true
	}

	return false
}

func (p *Processor) getOrCreateNotification(subscription *Subscription,
	item *Item) (*Notification, bool) {
	wasNew := false
	key := p.notificationCache.makeKey(item, subscription)
	notification := p.notificationCache.get(key)
	if !p.notificationCache.has(key) {
		notification = &Notification{Channel: subscription.Channel, Feed: *subscription.Feed, Item: *item}
		p.notificationCache.put(key, notification)
		wasNew = true
	}
	return notification, wasNew
}

func newNotificationCache() notificationCache {
	return notificationCache{
		channelAmount:     make(map[string]int),
		seenNotifications: map[string]*Notification{},
	}
}

type notificationCache struct {
	channelAmount     map[string]int
	seenNotifications map[string]*Notification
}

func (c *notificationCache) ChannelLimitReached(channelId string, limit int) bool {
	if amt, ok := c.channelAmount[channelId]; ok {
		return amt >= limit
	}
	return false
}

func (c *notificationCache) has(key string) bool {
	_, exists := c.seenNotifications[key]
	return exists
}

func (c *notificationCache) get(key string) *Notification {
	return c.seenNotifications[key]
}

func (c *notificationCache) put(key string, notification *Notification) {
	if _, ok := c.channelAmount[notification.Channel]; !ok {
		c.channelAmount[notification.Channel] = 0
	}
	c.channelAmount[notification.Channel] += 1
	c.seenNotifications[key] = notification
}

func (c *notificationCache) makeKey(item *Item, sub *Subscription) string {
	return item.GUID + sub.Channel
}
