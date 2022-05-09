// Integration test for FeedRepository.
package sqlite

import (
	"testing"

	"github.com/miodzie/seras/mods/rss"
)

var feedRepo FeedRepository

func TestFeedAll(t *testing.T) {
	feed := &rss.Feed{Name: "another_one", Url: "https://google.com/2"}
	feedRepo.Add(feed)
	feeds, err := feedRepo.All()
	if err != nil {
		t.Error(err)
	}
	if len(feeds) == 0 {
		t.Error(err)
	}
}

func TestFeedAdd(t *testing.T) {
	feed := &rss.Feed{Name: "hackernews", Url: "https://google.com"}
	err := feedRepo.Add(feed)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestFeedGetByName(t *testing.T) {
	err := feedRepo.Add(&rss.Feed{Name: "cool_name", Url: "https://google.com/cool_name"})
	if err != nil {
		t.Error(err)
	}

	feed, err := feedRepo.ByName("cool_name")
	if err != nil {
		t.Error(err)
	}
	if feed.Name != "cool_name" {
		t.Fail()
	}
}
