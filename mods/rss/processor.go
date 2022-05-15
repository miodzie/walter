package rss

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
		parsed, _ := p.parser.ParseURL(feed.Url)
		subs, _ := p.subs.ByFeedId(feed.Id)

		seen := make(map[string]*Notification)
		for _, sub := range subs {
			for _, item := range parsed.ItemsWithKeywords(sub.KeywordsSlice()) {
				key := item.GUID + sub.Channel

				if noti, ok := seen[key]; ok {
					noti.Users = append(noti.Users, sub.User)
				} else {
					notifications = append(notifications, &Notification{
						Channel: sub.Channel,
						Users:   []string{sub.User},
						Feed:    *feed,
						Item:    *item,
					})
					seen[key] = notifications[len(notifications)-1]
				}
			}
		}
	}

	return notifications, nil
}
