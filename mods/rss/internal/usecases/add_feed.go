// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import "github.com/miodzie/walter/mods/rss/internal/internal/domain"

type AddFeed struct {
	repository Repository
}

func NewAddFeed(repository Repository) *AddFeed {
	return &AddFeed{repository: repository}
}

type AddFeedRequest struct {
	Name string
	Url  string
}

type AddFeedResponse struct {
	Message string
}

func (a *AddFeed) Add(req AddFeedRequest) (AddFeedResponse, error) {
	resp := AddFeedResponse{Message: "Feed saved."}

	var feed domain.Feed
	feed.Name = req.Name
	feed.Url = req.Url

	err := a.repository.AddFeed(&feed)

	if err != nil {
		resp.Message = "Failed to save feed."
	}

	return resp, err
}
