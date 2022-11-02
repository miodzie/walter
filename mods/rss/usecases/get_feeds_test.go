// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"github.com/miodzie/walter/mods/rss"
	"testing"
)

func TestGetFeeds_Get(t *testing.T) {
	repository := rss.NewInMemRepo()
	repository.AddFeed(&rss.Feed{Name: "news"})
	getFeeds := NewGetFeeds(repository)

	// Act
	response, err := getFeeds.Exec()

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(response.Feeds) != 1 {
		t.Fail()
	}
	feed := response.Feeds[0]
	if feed.Name != "news" {
		t.Fail()
	}
}

func TestGetFeeds_Get_returns_empty_when_no_feeds(t *testing.T) {
	repository := rss.NewInMemRepo()
	getFeeds := NewGetFeeds(repository)

	// Act
	response, err := getFeeds.Exec()

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(response.Feeds) != 0 {
		t.Fail()
	}
	if response.Message != "No feeds available." {
		t.Fail()
	}
}
