package usecases

import (
	"errors"
	"github.com/miodzie/seras/mods/rss"
	"testing"
)

func TestSubscribe_Subscribe(t *testing.T) {
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

	// Act
	response, err := useCase.Subscribe(request)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if response.Message != "Subscribed to news with keywords: fire" {
		t.Fail()
	}
	subs, _ := repository.Subs(rss.SubSearchOpt{User: "adam"})
	if len(subs) != 1 {
		t.Fail()
	}
}

func TestSubscribe_Subscribe_fails_to_find_feed(t *testing.T) {
	repository := rss.NewInMemRepo()
	useCase := NewSubscribe(repository)

	// Act
	resp, err := useCase.Subscribe(SubscribeRequest{})

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
	resp, err := useCase.Subscribe(request)

	// Assert
	if resp.Message != "Failed to subscribe." {
		t.Fail()
	}
	if err != expectedErr {
		t.Error(err)
	}
}
