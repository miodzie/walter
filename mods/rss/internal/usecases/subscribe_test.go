// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"errors"
	"github.com/miodzie/walter/mods/rss/internal/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: Clean up tests

func TestSubscribeNoKeywords(t *testing.T) {
	repository := NewInMemRepo()
	subscribe := NewSubscribe(repository)
	feed := &domain.Feed{Id: 1, Name: "news"}
	repository.AddFeed(feed)
	request := SubscribeRequest{
		FeedName: "news",
		Channel:  "#news",
		User:     "adam",
	}

	response, err := subscribe.Subscribe(request)

	assert.Nil(t, err)
	assert.Equal(t, "Subscribed to news", response.Message)
	subs, _ := repository.Subs(SearchParams{User: "adam"})
	assert.Len(t, subs, 1)
	sub := subs[0]
	assert.Equal(t, "news", sub.Feed.Name)
	assert.Equal(t, "#news", sub.Channel)
	assert.Equal(t, "adam", sub.User)
}

func TestSubscribe_Exec_subscribes_a_user_to_a_feed(t *testing.T) {
	repository := NewInMemRepo()
	subscribe := NewSubscribe(repository)
	feed := &domain.Feed{Id: 1, Name: "news"}
	repository.AddFeed(feed)
	request := SubscribeRequest{
		FeedName: "news",
		Channel:  "#news",
		Keywords: "fire",
		User:     "adam",
	}

	response, err := subscribe.Subscribe(request)

	assert.Nil(t, err)
	assert.Equal(t, "Subscribed to news with keywords: fire", response.Message)
	subs, _ := repository.Subs(SearchParams{User: "adam"})
	assert.Len(t, subs, 1)
	sub := subs[0]
	assert.Equal(t, "news", sub.Feed.Name)
	assert.Equal(t, "#news", sub.Channel)
	assert.Equal(t, "fire", sub.Keywords)
	assert.Equal(t, "adam", sub.User)
}

func TestSubscribe_Exec_subscribes_a_user_with_ignore_words(t *testing.T) {
	repository := NewInMemRepo()
	subscribe := NewSubscribe(repository)
	feed := &domain.Feed{Id: 1, Name: "news"}
	repository.AddFeed(feed)
	request := SubscribeRequest{
		FeedName:    "news",
		Channel:     "#news",
		Keywords:    "fire",
		User:        "adam",
		IgnoreWords: "potato,salad",
	}

	response, err := subscribe.Subscribe(request)

	assert.Nil(t, err)
	assert.Equal(t,
		"Subscribed to news with keywords: fire. ignore: potato,salad",
		response.Message)
	subs, _ := repository.Subs(SearchParams{User: "adam"})
	assert.Len(t, subs, 1)
	sub := subs[0]
	assert.Equal(t, sub.Ignore, request.IgnoreWords)
}

func TestSubscribe_Exec_fails_to_find_feed(t *testing.T) {
	repository := NewInMemRepo()
	useCase := NewSubscribe(repository)

	resp, err := useCase.Subscribe(SubscribeRequest{})

	assert.Equal(t, FeedNotFoundError, err)
	assert.Equal(t, "Failed to find feed.", resp.Message)
}

func TestSubscribe_Exec_handles_repository_errors(t *testing.T) {
	repository := NewInMemRepo()
	useCase := NewSubscribe(repository)
	feed := &domain.Feed{Id: 1, Name: "news"}
	repository.AddFeed(feed)
	request := SubscribeRequest{
		FeedName: "news",
		Channel:  "#news",
		Keywords: "fire",
		User:     "adam",
	}
	expectedErr := errors.New("foo")
	repository.ForceError(expectedErr, 1)

	resp, err := useCase.Subscribe(request)

	assert.Equal(t, resp.Message, "Failed to subscribe.")
	assert.ErrorIs(t, expectedErr, err)
}
