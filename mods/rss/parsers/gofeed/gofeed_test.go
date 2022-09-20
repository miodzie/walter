package gofeed

import (
	"fmt"
	"testing"
)

func TestItCanParseARedditFeed(t *testing.T) {
	p := New()
	feed, err := p.ParseURL("https://www.reddit.com/r/news/.rss")
	if err != nil {
		fmt.Println(err)
		t.Fail()
		t.Log(err)
	}

	fmt.Println(feed)
}
