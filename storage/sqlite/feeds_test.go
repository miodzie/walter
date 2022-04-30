// Integration test for FeedRepository.
package sqlite

import (
	"fmt"
	"testing"

	"github.com/miodzie/seras/mods/rss"
)

var feedRepo FeedRepository


func TestAll(t *testing.T) {
	feed := &rss.Feed{Name: "another_one", Url: "https://google.com/2"}
	feedRepo.Save(feed)
	feeds, err := feedRepo.All()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if len(feeds) == 0 {
		t.Error(err)
		t.Fail()
	}
	for _, f := range feeds {
		fmt.Println(f)
	}
}

func TestSave(t *testing.T) {
	feed := &rss.Feed{Name: "hackernews", Url: "https://google.com"}
	err := feedRepo.Save(feed)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestGetByName(t *testing.T) {
	err := feedRepo.Save(&rss.Feed{Name: "cool_name", Url: "https://google.com/cool_name"})
	if err != nil {
		t.Error(err)
		t.Fail()
	}

    feed, err := feedRepo.GetByName("cool_name")
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    if feed.Name != "cool_name" {
        t.Fail()
    }
}
