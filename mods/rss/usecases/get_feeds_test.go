package usecases

import (
	"github.com/miodzie/seras/mods/rss"
	"testing"
)

func TestGetFeeds_Get(t *testing.T) {
	repository := rss.NewInMemRepo()
	repository.AddFeed(&rss.Feed{Name: "news"})
	useCase := NewGetFeeds(repository)

	// Act
	response, err := useCase.Get()

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(response.Feeds) != 1 {
		t.Fail()
	}
	feed := response.Feeds[0]
	if feed.Name != "news" {
		t.Fail()
	}
}

func TestGetFeeds_Get_returns_empty_when_no_feeds(t *testing.T) {
	repository := rss.NewInMemRepo()
	useCase := NewGetFeeds(repository)

	// Act
	response, err := useCase.Get()

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(response.Feeds) != 0 {
		t.Fail()
	}
	if response.Message != "No feeds available." {
		t.Fail()
	}
}
