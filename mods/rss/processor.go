// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"github.com/miodzie/walter/log"
	"sync"
)

type FeedProcessor struct {
	// Max notifications sent per channel per process.
	ChannelLimit      int
	repository        Repository
	parser            Parser
	notificationCache notificationCache
	sync.Mutex
}

func NewProcessor(repo Repository, parser Parser) *FeedProcessor {
	return &FeedProcessor{
		repository:   repo,
		parser:       parser,
		ChannelLimit: 3,
	}
}

func (p *FeedProcessor) Process() ([]*Notification, error) {
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
			newNotes := p.findNewNotifications(parsedFeed, subscription)
			notifications = append(notifications, newNotes...)
		}
	}

	return notifications, nil
}

func (p *FeedProcessor) findNewNotifications(parsedFeed *ParsedFeed, subscription *Subscription) []*Notification {
	var notifications []*Notification
	for _, item := range parsedFeed.ItemsWithKeywords(subscription.KeyWords()) {
		if p.shouldIgnore(subscription, item) {
			continue
		}
		subscription.Remember(*item)
		notification, wasNew := p.getOrCreateNotification(subscription, item)
		notification.Users = append(notification.Users, subscription.User)
		if wasNew {
			notifications = append(notifications, notification)
		}
	}

	err := p.repository.UpdateSub(subscription)
	if err != nil {
		log.Error(err)
	}

	return notifications
}

func (p *FeedProcessor) shouldIgnore(subscription *Subscription, item *Item) bool {
	return subscription.HasSeen(*item) ||
		(item.HasKeywords(subscription.IgnoreWords()) && subscription.Ignore != "") ||
		p.notificationCache.ChannelLimitReached(subscription.Channel, p.ChannelLimit)
}

func (p *FeedProcessor) getOrCreateNotification(subscription *Subscription,
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
	if _, exists := c.channelAmount[notification.Channel]; !exists {
		c.channelAmount[notification.Channel] = 0
	}
	c.channelAmount[notification.Channel] += 1
	c.seenNotifications[key] = notification
}

func (c *notificationCache) makeKey(item *Item, sub *Subscription) string {
	return item.GUID + sub.Channel
}
