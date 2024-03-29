// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package fetchers

import (
	"testing"

	"github.com/miodzie/walter/mods/rss"
)

func TestStripHtml(t *testing.T) {
	expected := "cool bean's!"
	feed := &rss.Feed{
		Title: "<strong>hello</strong> world!",
		Items: []rss.Item{
			{Description: "<img src=\"localhost\">cool bean&#39;s!"},
			{Description: "egg salad&#39;s!"},
		},
	}
	dummy := &rss.StubParser{Parsed: feed}
	sut := cleanHtml{BaseFetcher: dummy}

	parsed, _ := sut.Fetch("")

	if parsed.Title != "hello world!" {
		t.Error("failed to strip html")
	}

	d := parsed.Items[0].Description
	if d != expected {
		t.Error("failed to strip html")
		t.Errorf("expected: %s, got: %s", expected, d)
	}
	d = parsed.Items[1].Description
	expected = "egg salad's!"
	if d != expected {
		t.Error("failed to strip html")
		t.Errorf("expected: %s, got: %s", expected, d)
	}
}
