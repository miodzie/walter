// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"github.com/miodzie/walter/mods/rss/internal/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFeeds_Exec_gets_all_feeds(t *testing.T) {
	repository := NewInMemRepo()
	repository.AddFeed(&domain.Feed{Name: "news"})
	getFeeds := NewGetFeeds(repository)

	response, err := getFeeds.Get()

	assert.Nil(t, err)
	assert.Len(t, response.Feeds, 1)
	feed := response.Feeds[0]
	assert.Equal(t, "news", feed.Name)
}

func TestGetFeeds_Exec_returns_empty_when_no_feeds(t *testing.T) {
	repository := NewInMemRepo()
	getFeeds := NewGetFeeds(repository)

	response, err := getFeeds.Get()

	assert.Nil(t, err)
	assert.Empty(t, response.Feeds)
	assert.Equal(t, "No feeds available.", response.Message)
}
