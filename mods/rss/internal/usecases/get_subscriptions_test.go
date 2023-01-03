// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"walter/mods/rss/internal/internal/domain"
)

func TestGetSubscriptions_Exec_gets_subscriptions_for_a_user(t *testing.T) {
	repository := NewInMemRepo()
	getSubs := NewGetSubscriptions(repository)

	feed := &domain.Feed{Id: 1, Name: "news"}
	_ = repository.AddFeed(feed)
	subscription := &domain.Subscription{
		FeedId:  feed.Id,
		User:    "Bob",
		Channel: "#general",
	}
	_ = repository.AddSub(subscription)
	request := GetSubscriptionsRequest{
		User: "Bob",
		Optional: struct{ Channel string }{
			Channel: "#general",
		},
	}

	response, err := getSubs.Get(request)

	assert.Nil(t, err)
	assert.NotEmpty(t, response.Subscriptions, 0)
	for _, sub := range response.Subscriptions {
		assert.Equal(t, "#general", sub.Channel)
	}
}

func TestGetSubscriptions_Exec_handles_repository_errors(t *testing.T) {
	repo := NewInMemRepo()
	expectedErr := errors.New("testing")
	repo.ForceError(expectedErr, 0)
	getSubs := NewGetSubscriptions(repo)

	resp, err := getSubs.Get(GetSubscriptionsRequest{User: "Bob"})

	assert.ErrorIs(t, expectedErr, err)
	assert.Equal(t, "Failed to retrieve subscriptions.", resp.Message)
}
