// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"errors"
	"fmt"
	"github.com/miodzie/walter/mods/rss"
	"testing"
)

func TestNewUnsubscribeUseCase_Unsubscribe_unsubs_a_user(t *testing.T) {
	repository := rss.NewInMemRepo()
	feed := rss.Feed{Id: 1, Name: "news", Url: ""}
	if err := repository.AddFeed(&feed); err != nil {
		t.Fail()
	}
	subscription := rss.Subscription{
		FeedId:  feed.Id,
		User:    "john",
		Channel: "#news",
	}
	if err := repository.AddSub(&subscription); err != nil {
		t.Error(err)
	}
	unsub := NewUnsubscribe(repository)
	request := UnsubscribeRequest{Channel: "#news", FeedName: "news", User: "john"}

	// Act
	response, err := unsub.Exec(request)

	// Assert
	if err != nil {
		t.Error(err)
	}
	subs, err := repository.Subs(rss.SubSearchOpt{FeedId: feed.Id})
	if err != nil {
		t.Error(err)
	}
	if len(subs) != 0 {
		fmt.Println("There should be no subscriptions after a user unsubscribed.")
		t.Fail()
	}
	if response.Message != "Successfully unsubscribed from `news` feed." {
		t.Fail()
	}
}

func TestNewUnsubscribeUseCase_Unsubscribe_failed_to_find_sub(t *testing.T) {
	repository := rss.NewInMemRepo()
	expectedErr := errors.New("expected")
	repository.ForceError(expectedErr, 0)
	unsub := NewUnsubscribe(repository)

	// Act
	resp, err := unsub.Exec(UnsubscribeRequest{})

	// Assert
	if err != expectedErr {
		t.Error(err)
	}
	if resp.Message != "Failed to locate user subscription." {
		t.Fail()
	}
}

func TestNewUnsubscribeUseCase_Unsubscribe_failed_unsub(t *testing.T) {
	repository := rss.NewInMemRepo()
	feed := rss.Feed{Id: 1, Name: "news", Url: ""}
	if err := repository.AddFeed(&feed); err != nil {
		t.Fail()
	}
	subscription := rss.Subscription{
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

	// Act
	resp, err := unsub.Exec(request)

	// Assert
	if err != expectedErr {
		fmt.Println(resp.Message)
		t.Error(err)
	}
	if resp.Message != "Failed to unsubscribe." {
		t.Fail()
	}
}
