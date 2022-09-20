package usecases

import (
	"fmt"

	"github.com/miodzie/seras/mods/rss"
)

type Subscribe struct {
	repository rss.Repository
}

func NewSubscribeUseCase(repo rss.Repository) *Subscribe {
	return &Subscribe{repository: repo}
}

type SubscribeRequest struct {
	FeedName string
	Channel  string
	Keywords string
	User     string
}

type SubscribeResponse struct {
	Message string
}

func (s *Subscribe) Handle(req SubscribeRequest) (SubscribeResponse, error) {
	var resp SubscribeResponse

	feed, err := s.repository.FeedByName(req.FeedName)
	if err != nil {
		resp.Message = "Unknown feed."
		return resp, err
	}

	sub := &rss.Subscription{
		FeedId:   feed.Id,
		Channel:  req.Channel,
		Keywords: req.Keywords,
		User:     req.User,
	}
	resp.Message = fmt.Sprintf("Subscribed to %s with keywords: %s", feed.Name, sub.Keywords)
	if err = s.repository.AddSub(sub); err != nil {
		resp.Message = "Failed to save feed, likely one already exists for this channel and feed."
		return resp, err
	}

	return resp, nil
}
