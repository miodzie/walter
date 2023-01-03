// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package decorators

import (
	"testing"

	"github.com/miodzie/walter/mods/rss/internal/domain"
)

func TestStripHtml(t *testing.T) {
	expected := "cool bean's!"
	feed := &domain.ParsedFeed{
		Title: "<strong>hello</strong> world!",
		Items: []*domain.Item{{Description: "<img src=\"localhost\">cool bean&#39;s!"}},
	}
	dummy := &domain.StubParser{Parsed: feed}
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
