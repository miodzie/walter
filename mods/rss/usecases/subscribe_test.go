package usecases

import (
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
