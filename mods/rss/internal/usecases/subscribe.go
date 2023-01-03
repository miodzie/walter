// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"fmt"
	"walter/mods/rss/internal/internal/domain"
)

type Subscribe struct {
	repository Repository
}

// Subscribe Subscribes a user to a Feed.
func (s *Subscribe) Subscribe(req SubscribeRequest) (SubscribeResponse, error) {
	feed, err := s.repository.FeedByName(req.FeedName)
	if err != nil {
		return SubscribeResponse{
			Message: "Failed to find feed.",
		}, err
	}

	sub := &domain.Subscription{
		FeedId:   feed.Id,
		Channel:  req.Channel,
		Keywords: req.Keywords,
		User:     req.User,
		Ignore:   req.IgnoreWords,
	}
	if err = s.repository.AddSub(sub); err != nil {
		return SubscribeResponse{Message: "Failed to subscribe."}, err
	}

	reply := fmt.Sprintf("Subscribed to %s", feed.Name)
	if sub.Keywords != "" {
		reply += fmt.Sprintf(" with keywords: %s", sub.Keywords)
	}
	if sub.Ignore != "" {
		reply += fmt.Sprintf(". ignore: %s", sub.Ignore)
	}
	return SubscribeResponse{
		Message: reply,
	}, nil
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

func NewSubscribe(repo Repository) *Subscribe { return &Subscribe{repository: repo} }
