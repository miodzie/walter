// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"fmt"
)

type Unsubscribe struct {
	repository Repository
}

func (useCase *Unsubscribe) Unsubscribe(request UnsubscribeRequest) (UnsubscribeResponse, error) {
	subs, err := useCase.repository.Subs(SearchParams{
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

type UnsubscribeRequest struct {
	User     string
	Channel  string
	FeedName string
}

type UnsubscribeResponse struct {
	Message string
}

func NewUnsubscribe(repo Repository) *Unsubscribe { return &Unsubscribe{repository: repo} }
