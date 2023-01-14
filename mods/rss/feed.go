// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"github.com/miodzie/walter/log"
	"regexp"
	"strings"
)

// Fetcher downloads a UserFeed.Url and translates it to a Feed to
// be checked by a Subscription.
type Fetcher interface {
	Fetch(rssUrl string) (*Feed, error)
}

type Feed struct {
	Title       string
	Description string
	Link        string
	Updated     string
	Published   string
	Items       []Item
	Custom      map[string]string
}

func (feed *Feed) ItemsWithKeywords(keywords []string) []Item {
	var items []Item
	for _, i := range feed.Items {
		if i.HasKeywords(keywords) {
			items = append(items, i)
		}
	}
	return items
}

func (feed *Feed) HasKeywords(keywords []string) bool {
	for _, item := range feed.Items {
		if item.HasKeywords(keywords) {
			return true
		}
	}

	return false
}

type Item struct {
	Title       string
	Description string
	Content     string
	Link        string
	Links       []string
	GUID        string
	Custom      map[string]string
}

func (i *Item) DescriptionTruncated() string {
	if len(i.Description) < 100 {
		return i.Description
	}
	sp := strings.Split(i.Description, "")

	return strings.Join(sp[:100], "") + "..."
}

func (i *Item) HasKeywords(keywords []string) bool {
	for _, keyword := range keywords {
		reg, err := createWordBoundaryRegex(keyword)
		if err != nil {
			log.Error(err)
			continue
		}
		if reg.MatchString(i.Title) || reg.MatchString(i.Description) || reg.MatchString(i.Content) {
			return true
		}
	}

	return false
}

const WordBoundary = `(?i)\b$WORD$\b`

func createWordBoundaryRegex(word string) (*regexp.Regexp, error) {
	return regexp.Compile(
		strings.Replace(WordBoundary,
			"$WORD$",
			regexp.QuoteMeta(word),
			1))
}

type StubParser struct {
	Parsed *Feed
}

func (p *StubParser) Fetch(url string) (*Feed, error) {
	return p.Parsed, nil
}
