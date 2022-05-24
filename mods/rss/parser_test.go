package rss

import (
	"strings"
	"testing"
)

func TestItem_Desc_returns_the_string_if_less_than_100_chars(t *testing.T) {
	item := &Item{Description: ""}
	for i := 0; i < 99; i++ {
		item.Description += "A"
	}

	if item.Description != item.Desc() {
		t.Fail()
	}
}

func TestItem_Desc_shortens_the_description_length_to_100_chars(t *testing.T) {
	item := &Item{Description: ""}
	for i := 0; i < 150; i++ {
		item.Description += "A"
	}
	sp := strings.Split(item.Description, "")
	expected := strings.Join(sp[:100], "") + "..."

	if expected != item.Desc() {
		t.Fail()
	}
}

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
