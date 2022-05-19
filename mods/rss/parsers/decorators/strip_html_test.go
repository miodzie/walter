package decorators

import (
	"fmt"
	"testing"

	"github.com/miodzie/seras/mods/rss"
)

func TestStripHtml(t *testing.T) {
	feed := &rss.ParsedFeed{
		Title: "<strong>hello</strong> world!",
		Items: []*rss.Item{{Description: "<img src=\"localhost\">cool beans!"}},
	}
	dummy := &rss.NulledParser{Parsed: feed}
	sut := cleanHtml{BaseParser: dummy}

	parsed, _ := sut.ParseURL("")

	if parsed.Title != "hello world!" {
		t.Error("failed to strip html")
	}

	if parsed.Items[0].Description != "cool beans!" {
		fmt.Println(parsed.Title)
		t.Log(parsed.Items[0].Description)
		t.Error("failed to strip html")
	}
}
