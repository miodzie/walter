package interactors

import (
	"fmt"

	"github.com/miodzie/seras/mods/rss"
)

type Subscribe struct {
	Feeds rss.Feeds
	Subs  rss.Subscriptions
}

type SubscribeRequest struct {
	FeedName string
	Channel  string
	Keywords string
	User     string
}

type SubscribeResponse struct {
	Message string
	Error   error
}

func (s *Subscribe) Handle(req SubscribeRequest) SubscribeResponse {
	var resp SubscribeResponse

	feed, err := s.Feeds.ByName(req.FeedName)
	if err != nil {
		resp.Message = "Unknown feed. Use !feeds to see available."
		resp.Error = err
		return resp
	}

	sub := &rss.Subscription{
		FeedId:   feed.Id,
		Channel:  req.Channel,
		Keywords: req.Keywords,
		User:     req.User,
	}
	resp.Message = fmt.Sprintf("Subscribed to %s with keywords: %s", feed.Name, sub.Keywords)
	if resp.Error = s.Subs.Add(sub); resp.Error != nil {
		resp.Message = "Failed to save feed, likely one already exists for this channel and feed."
	}

	return resp
}
