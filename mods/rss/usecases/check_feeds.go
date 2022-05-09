package usecases

import (
	"fmt"

	"github.com/miodzie/seras/mods/rss"
)

type CheckFeeds struct {
	Feeds rss.FeedRepository
	Subs  rss.SubscriptionRepository
}

type CheckFeedsResponse struct {
	Notifications []*rss.Notification
	Error         error
}

func (cl *CheckFeeds) Handle() CheckFeedsResponse {
	var resp CheckFeedsResponse
	feeds, err := cl.Feeds.All()
	if err != nil {
		resp.Error = err
		return resp
	}

	for _, listing := range feeds {
		fmt.Printf("Checking feed: %s: %s\n", listing.Name, listing.Url)
		subs, err := cl.Subs.ByFeedId(listing.Id)
		if err != nil {
			resp.Error = err
		}

		notifications := make(map[string][]*rss.Subscription)
		for _, sub := range subs {
			//  Fetch feed.
			//  Check feed for sub keywords.
			// Add sub to SubsToNotify slice if so.
			fmt.Printf("Channel: %s User: %s, Keywords: %s\n", sub.Channel, sub.User, sub.Keywords)
			shouldNotify := true
			if shouldNotify {
				key := fmt.Sprintf("%d-%s", sub.FeedId, sub.Channel)
				notifications[key] = append(notifications[key], &sub)
			}
		}

		// Sort by Channel and Feed, and group Users to notify into a
		// new struct.
		for _, subs := range notifications {
			notification := &rss.Notification{}
			for _, sub := range subs {
				notification.Channel = sub.Channel // TODO: Set only once
				notification.Users = append(notification.Users, sub.User)
			}
			resp.Notifications = append(resp.Notifications, notification)
		}
	}

	return resp
}
