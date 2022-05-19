package rss

import "testing"

func TestParsedFeed_ItemsWithKeywords(t *testing.T) {
	sub := &Subscription{Keywords: "foo,bar"}
	feed := &ParsedFeed{
		Items: []*Item{{Title: "foo"}, {Title: "bar"}, {Title: "baz"}},
	}

	result := feed.ItemsWithKeywords(sub.KeywordsSlice())

	if len(result) != 2 {
		t.Error("should contain 2 result")
	}
}

func TestParsedFeed_HasKeywords(t *testing.T) {
	sub := &Subscription{Keywords: "fOo,bar,baz"}
	for _, feed := range hasKeywords {
		if !feed.HasKeywords(sub.KeywordsSlice()) {
			t.Error("feed should have keyword")
		}
	}

	for _, feed := range hasNotKeywords {
		if feed.HasKeywords(sub.KeywordsSlice()) {
			t.Error("feed should not have keyword")
		}
	}
}

var hasKeywords = []*ParsedFeed{
	{Items: []*Item{
		{Title: "FOO"},
	}},
	{Items: []*Item{
		{Title: "baz"},
	}},
	{Items: []*Item{
		{Description: "foo"},
	}},
	{Items: []*Item{
		{Content: "foo"},
	}},
}

var hasNotKeywords = []*ParsedFeed{
	{Items: []*Item{
		{Title: "zab"},
	}},
	{Items: []*Item{
		{Description: "oof"},
	}},
	{Items: []*Item{
		{Content: "oof"},
	}},
}
