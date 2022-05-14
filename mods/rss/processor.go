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
		parsed, _ := p.parser.Parse(feed.Url)
		subs, _ := p.subs.ByFeedId(feed.Id)
		for _, sub := range subs {
			if parsed.HasKeywords(sub.KeywordsSlice()) {
				notifications = append(notifications, &Notification{
					Channel: sub.Channel,
					Users:   []string{sub.User},
					Feed:    *feed,
				})
			}
		}
	}

	return notifications, nil
}
