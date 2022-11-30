// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"errors"
	"github.com/miodzie/walter/mods/rss"
	"testing"
)

func TestSubscribe_Subscribe(t *testing.T) {
	repository := rss.NewInMemRepo()
	subscribe := NewSubscribe(repository)
	feed := &rss.Feed{Id: 1, Name: "news"}
	repository.AddFeed(feed)
	request := SubscribeRequest{
		FeedName: "news",
		Channel:  "#news",
		Keywords: "fire",
		User:     "adam",
	}

	// Act
	response, err := subscribe.Exec(request)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if response.Message != "Subscribed to news with keywords: fire" {
		t.Fail()
	}
	subs, _ := repository.Subs(rss.SearchParams{User: "adam"})
	if len(subs) != 1 {
		t.Fail()
	}
}

func TestSubscribe_Subscribe_with_ignore(t *testing.T) {
	repository := rss.NewInMemRepo()
	subscribe := NewSubscribe(repository)
	feed := &rss.Feed{Id: 1, Name: "news"}
	repository.AddFeed(feed)
	request := SubscribeRequest{
		FeedName:    "news",
		Channel:     "#news",
		Keywords:    "fire",
		User:        "adam",
		IgnoreWords: "potato,salad",
	}

	// Act
	response, err := subscribe.Exec(request)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if response.Message != "Subscribed to news with keywords: fire. "+
		"ignore: potato,salad" {
		t.Fail()
	}
	subs, _ := repository.Subs(rss.SearchParams{User: "adam"})
	if len(subs) != 1 {
		t.Fail()
	}
	sub := subs[0]
	if sub.Ignore != request.IgnoreWords {
		t.Fail()
	}

}

func TestSubscribe_Subscribe_fails_to_find_feed(t *testing.T) {
	repository := rss.NewInMemRepo()
	useCase := NewSubscribe(repository)

	// Act
	resp, err := useCase.Exec(SubscribeRequest{})

	// Assert
	if err.Error() != "feed not found" {
		t.Error(err)
	}
	if resp.Message != "Failed to find feed." {
		t.Fail()
	}
}

func TestSubscribe_Subscribe_fails_to_subscribe(t *testing.T) {
	repository := rss.NewInMemRepo()
	useCase := NewSubscribe(repository)
	feed := &rss.Feed{Id: 1, Name: "news"}
	repository.AddFeed(feed)
	request := SubscribeRequest{
		FeedName: "news",
		Channel:  "#news",
		Keywords: "fire",
		User:     "adam",
	}
	expectedErr := errors.New("foo")
	repository.ForceError(expectedErr, 1)

	// Act
	resp, err := useCase.Exec(request)

	// Assert
	if resp.Message != "Failed to subscribe." {
		t.Fail()
	}
	if err != expectedErr {
		t.Error(err)
	}
}
