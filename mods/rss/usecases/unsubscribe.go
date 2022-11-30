// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"fmt"
	"github.com/miodzie/walter/mods/rss"
)

type Unsubscribe struct {
	repository rss.Repository
}

func NewUnsubscribe(repo rss.Repository) *Unsubscribe {
	return &Unsubscribe{repository: repo}
}

type UnsubscribeRequest struct {
	User     string
	Channel  string
	FeedName string
}

type UnsubscribeResponse struct {
	Message string
}

func (useCase *Unsubscribe) Exec(
	request UnsubscribeRequest) (UnsubscribeResponse, error) {
	subs, err := useCase.repository.Subs(rss.SearchParams{
		User:     request.User,
		FeedName: request.FeedName,
		Channel:  request.Channel,
	})

	if err != nil || len(subs) != 1 {
		return UnsubscribeResponse{Message: "Failed to locate user subscription."}, err
	}
	err = useCase.repository.RemoveSub(subs[0])
	if err != nil {
		return UnsubscribeResponse{Message: "Failed to unsubscribe."}, err
	}

	return UnsubscribeResponse{
		Message: fmt.Sprintf("Successfully unsubscribed from `%s` feed.", request.FeedName),
	}, nil
}
