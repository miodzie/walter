package usecases

import (
	"fmt"
	"github.com/miodzie/seras/mods/rss"
	"testing"
)

func TestNewUnsubscribeUseCase_Unsubscribe_unsubs_a_user(t *testing.T) {
	repository := rss.NewInMemRepo()
	feed := rss.Feed{Id: 1, Name: "news", Url: ""}
	if err := repository.AddFeed(&feed); err != nil {
		t.Log(err)
		t.Fail()
	}
	subscription := rss.Subscription{
		FeedId:  feed.Id,
		User:    "john",
		Channel: "#news",
	}
	if err := repository.AddSub(&subscription); err != nil {
		t.Log(err)
		t.Fail()
	}
	useCase := NewUnsubscribe(repository)
	request := UnsubscribeRequest{Channel: "#news", FeedName: "news", User: "john"}

	// Act
	response, err := useCase.Unsubscribe(request)

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
