package usecases

import (
	"fmt"
	"github.com/miodzie/seras/mods/rss"
)

type Subscribe struct {
	repository rss.Repository
}

func NewSubscribe(repo rss.Repository) *Subscribe {
	return &Subscribe{repository: repo}
}

type SubscribeRequest struct {
	FeedName    string
	Channel     string
	Keywords    string
	User        string
	IgnoreWords string
}

type SubscribeResponse struct {
	Message string
}

// Subscribe subscribes a user to a Feed.
func (s *Subscribe) Subscribe(req SubscribeRequest) (SubscribeResponse, error) {
	feed, err := s.repository.FeedByName(req.FeedName)
	if err != nil {
		return SubscribeResponse{
			Message: "Failed to find feed.",
		}, err
	}

	sub := &rss.Subscription{
		FeedId:   feed.Id,
		Channel:  req.Channel,
		Keywords: req.Keywords,
		User:     req.User,
		Ignore:   req.IgnoreWords,
	}
	if err = s.repository.AddSub(sub); err != nil {
		return SubscribeResponse{Message: "Failed to subscribe."}, err
	}

	reply := fmt.Sprintf("Subscribed to %s with keywords: %s", feed.Name, sub.Keywords)
	if sub.Ignore != "" {
		reply += fmt.Sprintf(". ignore: %s", sub.Ignore)
	}
	return SubscribeResponse{
		Message: reply,
	}, nil
}
