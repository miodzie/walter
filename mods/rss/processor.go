package rss

import "fmt"

type Processor struct {
	repo Repository
	parser Parser
}

func NewProcessor(repo Repository, parser Parser) *Processor {
	return &Processor{
        repo: repo,
		parser: parser,
	}
}

func (p *Processor) Handle() ([]*Notification, error) {
	var notifications []*Notification

	feeds, _ := p.repo.AllFeeds()
	for _, feed := range feeds {
		parsed, err := p.parser.ParseURL(feed.Url)
		if err != nil {
			return notifications, err
		}
		subs, err := p.repo.SubByFeedId(feed.Id)
		if err != nil {
			return notifications, err
		}
		seen := make(map[string]*Notification)
		for _, sub := range subs {
			for _, item := range parsed.ItemsWithKeywords(sub.KeywordsSlice()) {
				if sub.HasSeen(*item) {
					continue
				}
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
				sub.See(*item)
			}
			err := p.repo.UpdateSub(sub)
			if err != nil {
				// TODO: remove?
				fmt.Println(err)
			}
		}
	}

	return notifications, nil
}
