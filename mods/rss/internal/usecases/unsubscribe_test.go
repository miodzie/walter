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

func TestNewUnsubscribeUseCase_Exec_unsubscribes_a_user(t *testing.T) {
	repository := NewInMemRepo()
	feed := domain.Feed{Id: 1, Name: "news", Url: ""}
	if err := repository.AddFeed(&feed); err != nil {
		t.Fail()
	}
	subscription := domain.Subscription{
		FeedId:  feed.Id,
		User:    "john",
		Channel: "#news",
	}
	if err := repository.AddSub(&subscription); err != nil {
		t.Error(err)
	}
	unsub := NewUnsubscribe(repository)
	request := UnsubscribeRequest{Channel: "#news", FeedName: "news", User: "john"}

	response, err := unsub.Unsubscribe(request)

	assert.Nil(t, err)
	subs, err := repository.Subs(SearchParams{FeedId: feed.Id})
	assert.Nil(t, err)
	assert.Empty(t, subs, "There should be no subscriptions after a user unsubscribed.")
	assert.Equal(t, "Successfully unsubscribed from `news` feed.", response.Message)
}

func TestNewUnsubscribeUseCase_Exec_failed_to_find_sub(t *testing.T) {
	repository := NewInMemRepo()
	expectedErr := errors.New("expected")
	repository.ForceError(expectedErr, 0)
	unsub := NewUnsubscribe(repository)

	resp, err := unsub.Unsubscribe(UnsubscribeRequest{})

	assert.ErrorIs(t, expectedErr, err)
	assert.Equal(t, "Failed to locate user subscription.", resp.Message)
}

func TestNewUnsubscribeUseCase_Exec_handles_repository_errors(t *testing.T) {
	repository := NewInMemRepo()
	feed := domain.Feed{Id: 1, Name: "news", Url: ""}
	if err := repository.AddFeed(&feed); err != nil {
		t.Fail()
	}
	subscription := domain.Subscription{
		FeedId:  feed.Id,
		User:    "john",
		Channel: "#news",
	}
	if err := repository.AddSub(&subscription); err != nil {
		t.Error(err)
	}

	expectedErr := errors.New("expected")
	repository.ForceError(expectedErr, 1)
	unsub := NewUnsubscribe(repository)
	request := UnsubscribeRequest{Channel: "#news", FeedName: "news", User: "john"}

	resp, err := unsub.Unsubscribe(request)

	assert.ErrorIs(t, expectedErr, err)
	assert.Equal(t, "Failed to unsubscribe.", resp.Message)
}
