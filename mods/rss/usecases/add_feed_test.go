// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"errors"
	"github.com/miodzie/walter/mods/rss"
	"testing"
)

func TestAddFeed_Exec_adds_a_new_feed_to_the_repository(t *testing.T) {
	repository := rss.NewInMemRepo()
	addFeed := NewAddFeed(repository)

	// Act
	response, err := addFeed.Exec(AddFeedRequest{Name: "foo", Url: "http://localhost.rss"})

	// Assert
	if err != nil {
		t.Error(err)
	}
	if response.Message != "Feed saved." {
		t.Fail()
	}
	feed, err := repository.FeedByName("foo")
	if err != nil {
		t.Error(err)
	}
	if feed.Name != "foo" {
		t.Error("feed was not saved to repository")
	}
}

func TestAddFeed_Exec_handles_repository_errors(t *testing.T) {
	repository := rss.NewInMemRepo()
	expectedErr := errors.New("my error")
	repository.ForceError(expectedErr, 0)
	addFeed := NewAddFeed(repository)

	// Act
	response, err := addFeed.Exec(AddFeedRequest{})

	// Assert
	if err != expectedErr {
		t.Error(err)
	}
	if response.Message != "Failed to save feed." {
		t.Fail()
	}
}
