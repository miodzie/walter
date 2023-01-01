// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

// NoteSorter sorts new notifications for the subscriptions and feed,
// ignoring already seen items for a subscription.
type NoteSorter struct {
	cache        noteCache
	channelLimit int
}

func NewNoteSorter(channelLimit int) *NoteSorter {
	return &NoteSorter{
		channelLimit: channelLimit,
		cache:        newNoteCache()}
}

func (s *NoteSorter) sort(subs []*Subscription, feed *ParsedFeed) (notes []*Notification) {
	for _, sub := range subs {
		newNotes := s.findNewNotificationsFor(sub, feed)
		notes = append(notes, newNotes...)
	}
	return notes
}

func (s *NoteSorter) findNewNotificationsFor(sub *Subscription, feed *ParsedFeed) []*Notification {
	var notes []*Notification
	for _, item := range feed.ItemsWithKeywords(sub.KeyWords()) {
		// TODO(refactor): Don't like how this and sub.Remember()
		//   are called down here. Maybe create an interface? "IgnoreStrategy"
		if s.shouldIgnore(sub, item) {
			continue
		}
		sub.Remember(*item)
		notification, wasNew := s.getOrCreateNotification(sub, item)
		notification.Users = append(notification.Users, sub.User)
		if wasNew {
			notes = append(notes, notification)
		}
	}
	return notes
}

func (s *NoteSorter) shouldIgnore(subscription *Subscription, item *Item) bool {
	return subscription.HasSeen(*item) ||
		(item.HasKeywords(subscription.IgnoreWords()) && subscription.Ignore != "") ||
		s.cache.ChannelLimitReached(subscription.Channel, s.channelLimit)
}

func (s *NoteSorter) getOrCreateNotification(subscription *Subscription, item *Item) (*Notification, bool) {
	wasNew := false
	key := s.cache.makeKey(item, subscription)
	notification := s.cache.get(key)
	if !s.cache.has(key) {
		notification = &Notification{Channel: subscription.Channel, Feed: *subscription.Feed, Item: *item}
		s.cache.put(key, notification)
		wasNew = true
	}
	return notification, wasNew
}

func newNoteCache() noteCache {
	return noteCache{
		channelAmount:     make(map[string]int),
		seenNotifications: map[string]*Notification{},
	}
}

type noteCache struct {
	channelAmount     map[string]int
	seenNotifications map[string]*Notification
}

func (c *noteCache) ChannelLimitReached(channelId string, limit int) bool {
	if amt, ok := c.channelAmount[channelId]; ok {
		return amt >= limit
	}
	return false
}

func (c *noteCache) has(key string) bool {
	_, exists := c.seenNotifications[key]
	return exists
}

func (c *noteCache) get(key string) *Notification {
	return c.seenNotifications[key]
}

func (c *noteCache) put(key string, notification *Notification) {
	if _, exists := c.channelAmount[notification.Channel]; !exists {
		c.channelAmount[notification.Channel] = 0
	}
	c.channelAmount[notification.Channel] += 1
	c.seenNotifications[key] = notification
}

func (c *noteCache) makeKey(item *Item, sub *Subscription) string {
	return item.GUID + sub.Channel
}
