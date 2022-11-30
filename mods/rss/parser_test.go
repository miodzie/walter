// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestItem_HasKeywords_Ignores_words_within_a_word(t *testing.T) {
	item := &Item{Title: "Financially"}
	assert.False(t, item.HasKeywords([]string{"CIA"}))
	assert.True(t, item.HasKeywords([]string{"financially"}))
}

func TestItem_DescTruncated_returns_the_string_if_less_than_100_chars(t *testing.T) {
	item := &Item{Description: ""}
	for i := 0; i < 99; i++ {
		item.Description += "A"
	}

	assert.Equal(t, item.Description, item.DescriptionTruncated())
}

func TestItem_DescTruncated_shortens_the_description_length_to_100_chars(t *testing.T) {
	item := &Item{Description: ""}
	for i := 0; i < 150; i++ {
		item.Description += "A"
	}
	sp := strings.Split(item.Description, "")
	expected := strings.Join(sp[:100], "") + "..."

	assert.Equal(t, expected, item.DescriptionTruncated())
}

func TestParsedFeed_ItemsWithKeywords(t *testing.T) {
	sub := &Subscription{Keywords: "foo,bar"}
	feed := &ParsedFeed{
		Items: []*Item{{Title: "foo"}, {Title: "bar"}, {Title: "baz"}},
	}

	result := feed.ItemsWithKeywords(sub.KeyWords())

	assert.Len(t, result, 2)
}

func TestParsedFeed_HasKeywords(t *testing.T) {
	sub := &Subscription{Keywords: "fOo,bar,baz"}
	for _, feed := range hasKeywords {
		assert.True(t, feed.HasKeywords(sub.KeyWords()), sub.Keywords)
	}

	for _, feed := range hasNotKeywords {
		assert.False(t, feed.HasKeywords(sub.KeyWords()))
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
