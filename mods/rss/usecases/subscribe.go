package usecases

import "github.com/miodzie/seras/mods/rss"

type Subscribe struct {
	Feeds rss.FeedRepository
	Subs  rss.SubscriptionRepository
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

	feed, err := s.Feeds.GetByName(req.FeedName)
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
	resp.Error = s.Subs.Save(sub)
	if resp.Error != nil {
		resp.Message = "Failed to save feed, likely one already exists for this channel and feed."
	}

	return resp
}
