package rss

import (
	"fmt"
)

type Processor struct {
	feeds  Feeds
	subs   Subscriptions
	parser Parser
}

func NewProcessor(f Feeds, s Subscriptions, parser Parser) *Processor {
	return &Processor{
		feeds:  f,
		subs:   s,
		parser: parser,
	}
}

func (p *Processor) Handle() ([]*Notification, error) {
	var notifications []*Notification

	feeds, _ := p.feeds.All()
	for _, feed := range feeds {
		parsed, _ := p.parser.Parse(feed.Url)
		subs, _ := p.subs.ByFeedId(feed.Id)
		seenNotifs := make(map[string]*Notification)
		for _, sub := range subs {
			if parsed.HasKeywords(sub.KeywordsSlice()) {
				key := fmt.Sprintf("%d-%s", feed.Id, sub.Channel)
				if seen, ok := seenNotifs[key]; ok {
					seen.Users = append(seen.Users, sub.User)
				} else {
					notifications = append(notifications, &Notification{
						Channel: sub.Channel,
						Users:   []string{sub.User},
						Feed:    *feed,
					})
					seenNotifs[key] = notifications[len(notifications)-1]
				}
			}
		}
	}

	return notifications, nil
}
