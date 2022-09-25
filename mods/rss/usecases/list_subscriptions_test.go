package usecases

import (
	"github.com/miodzie/seras/mods/rss"
	"testing"
)

func TestListSubscriptions_ListSubscriptions(t *testing.T) {
	repo := rss.NewInMemRepo()
	useCase := NewListSubscriptions(repo)

	feed := &rss.Feed{Id: 1, Name: "news"}
	_ = repo.AddFeed(feed)
	subscription := &rss.Subscription{
		FeedId:  feed.Id,
		User:    "Bob",
		Channel: "#general",
	}
	_ = repo.AddSub(subscription)
	request := ListSubscriptionsRequest{
		User: "Bob",
		Optional: struct{ Channel string }{
			Channel: "#general",
		},
	}

	// Act
	response := useCase.ListSubscriptions(request)

	// Assert
	if response.Error != nil {
		t.Log(response.Error)
		t.Fail()
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
