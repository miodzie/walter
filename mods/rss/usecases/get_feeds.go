// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"github.com/miodzie/seras/mods/rss"
)

type GetFeeds struct{ repository rss.Repository }

func NewGetFeeds(repository rss.Repository) *GetFeeds {
	return &GetFeeds{repository: repository}
}

type GetFeedsResponse struct {
	Feeds   []*GetFeedsFeed
	Message string
}

type GetFeedsFeed struct {
	Name string
	Url  string
}

func (g *GetFeeds) Exec() (GetFeedsResponse, error) {
	var resp GetFeedsResponse
	feeds, err := g.repository.Feeds()
	for _, f := range feeds {
		resp.Feeds = append(resp.Feeds, &GetFeedsFeed{Name: f.Name, Url: f.Url})
	}

	if len(resp.Feeds) == 0 {
		resp.Message = "No feeds available."
	}

	if err != nil {
		resp.Message = "Failed to fetch feeds."
	}

	return resp, err
}
