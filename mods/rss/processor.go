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

func (c *Processor) Handle() ([]*Notification, error) {
	var notifications []*Notification

	n := &Notification{
		Channel: "#chat",
		Users:   []string{"adam"},
		Feed:    Feed{Id: 1},
	}
	notifications = append(notifications, n)

	return notifications, nil
}
