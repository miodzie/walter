package usecases

import (
	"errors"
	"github.com/miodzie/seras/mods/rss"
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
	response, err := getSubs.Get(request)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(response.Subscriptions) == 0 {
		t.Log("no subscriptions returned")
		t.Fail()
	}
	for _, sub := range response.Subscriptions {
		if sub.Channel != "#general" {
			t.Fail()
		}
	}
}

func TestGetSubscriptions_Get_error(t *testing.T) {
	repo := rss.NewInMemRepo()
	expectedErr := errors.New("testing")
	repo.ForceError(expectedErr, 0)
	getSubs := NewGetSubscriptions(repo)

	// Act
	resp, err := getSubs.Get(GetSubscriptionsRequest{User: "Bob"})

	// Assert
	if err != expectedErr {
		t.Error(err)
	}
	if resp.Message != "Failed to retrieve subscriptions." {
		t.Fail()
	}
}
