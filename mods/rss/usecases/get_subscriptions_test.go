// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"errors"
	"github.com/miodzie/walter/mods/rss"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSubscriptions_Get(t *testing.T) {
	repo := rss.NewInMemRepo()
	getSubs := NewGetSubscriptions(repo)

	feed := &rss.Feed{Id: 1, Name: "news"}
	_ = repo.AddFeed(feed)
	subscription := &rss.Subscription{
		FeedId:  feed.Id,
		User:    "Bob",
		Channel: "#general",
	}
	_ = repo.AddSub(subscription)
	request := GetSubscriptionsRequest{
		User: "Bob",
		Optional: struct{ Channel string }{
			Channel: "#general",
		},
	}

	// Act
	response, err := getSubs.Exec(request)

	// Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, response.Subscriptions, 0)
	for _, sub := range response.Subscriptions {
		assert.Equal(t, "#general", sub.Channel)
	}
}

func TestGetSubscriptions_Get_error(t *testing.T) {
	repo := rss.NewInMemRepo()
	expectedErr := errors.New("testing")
	repo.ForceError(expectedErr, 0)
	getSubs := NewGetSubscriptions(repo)

	// Act
	resp, err := getSubs.Exec(GetSubscriptionsRequest{User: "Bob"})

	// Assert
	assert.ErrorIs(t, expectedErr, err)
	assert.Equal(t, "Failed to retrieve subscriptions.", resp.Message)
}
