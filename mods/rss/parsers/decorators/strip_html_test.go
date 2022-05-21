package decorators

import (
	"testing"

	"github.com/miodzie/seras/mods/rss"
)

func TestStripHtml(t *testing.T) {
	expected := "cool bean's!"
	feed := &rss.ParsedFeed{
		Title: "<strong>hello</strong> world!",
		Items: []*rss.Item{{Description: "<img src=\"localhost\">cool bean&#39;s!"}},
	}
	dummy := &rss.NulledParser{Parsed: feed}
	sut := cleanHtml{BaseParser: dummy}

	parsed, _ := sut.ParseURL("")

	if parsed.Title != "hello world!" {
		t.Error("failed to strip html")
	}

	d := parsed.Items[0].Description
	if d != expected {
		t.Error("failed to strip html")
		t.Errorf("expected: %s, got: %s", expected, d)
	}
}
