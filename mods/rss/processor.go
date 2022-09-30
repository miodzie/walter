package rss

import (
	"fmt"
	"github.com/miodzie/seras/log"
)

type Processor struct {
	repository Repository
	parser     Parser
	// Max notifications sent per channel per process.
	ChannelLimit int
}

func NewProcessor(repo Repository, parser Parser) *Processor {
	return &Processor{
		repository:   repo,
		parser:       parser,
		ChannelLimit: 3,
	}
}

func (p *Processor) Process() ([]*Notification, error) {
	var newNotifications []*Notification
	feeds, _ := p.repository.Feeds()
	for _, feed := range feeds {
		log.Debug("feed: " + feed.Name)
		parsedFeed, err := p.parser.ParseURL(feed.Url)
		if err != nil {
			return newNotifications, err
		}

		subs, err := p.repository.Subs(SubSearchOpt{FeedId: feed.Id})
		if err != nil {
			return newNotifications, err
		}

		cache := newCache()
		for _, subscription := range subs {
			for _, item := range parsedFeed.ItemsWithKeywords(subscription.KeywordsSlice()) {
				if subscription.HasSeen(*item) {
					continue
				}
				key := cache.makeKey(item, subscription)
				if cache.ChannelLimitReached(subscription.Channel, p.ChannelLimit) {
					// Save it for next time, bucko.
					continue
				}
				notification := cache.GetNotification(key)
				if !cache.HasNotification(key) {
					notification = &Notification{
						Channel: subscription.Channel,
						Feed:    *feed,
						Item:    *item,
					}
					newNotifications = append(newNotifications, notification)
					cache.PutNotification(key, notification)
				}

				notification.Users = append(notification.Users, subscription.User)
				subscription.See(*item)
			}
			err := p.repository.UpdateSub(subscription)
			if err != nil {
				log.Error(err)
			}
		}
	}

	return newNotifications, nil
}

func newCache() *cache {
	return &cache{
		channelAmount:     make(map[string]int),
		seenNotifications: map[string]*Notification{},
	}
}

type cache struct {
	channelAmount     map[string]int
	seenNotifications map[string]*Notification
}

func (c *cache) HasNotification(key string) bool {
	_, exists := c.seenNotifications[key]
	return exists
}

func (c *cache) GetNotification(key string) *Notification {
	return c.seenNotifications[key]
}

func (c *cache) PutNotification(key string, notification *Notification) {
	if _, ok := c.channelAmount[notification.Channel]; !ok {
		c.channelAmount[notification.Channel] = 0
	}
	c.channelAmount[notification.Channel] += 1
	c.seenNotifications[key] = notification
}

func (c *cache) makeKey(item *Item, sub *Subscription) string {
	return item.GUID + sub.Channel
}

func (c *cache) ChannelLimitReached(channelId string, limit int) bool {
	if amt, ok := c.channelAmount[channelId]; ok {
		fmt.Println(amt)
		return amt >= limit
	}
	return false
}
