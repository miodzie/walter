package domain

import (
	"github.com/miodzie/walter/log"
	"regexp"
	"strings"
)

type ParsedFeed struct {
	Title       string
	Description string
	Link        string
	Updated     string
	Published   string
	Items       []*Item
	Custom      map[string]string
}

func (feed *ParsedFeed) ItemsWithKeywords(keywords []string) []*Item {
	var items []*Item
	for _, i := range feed.Items {
		if i.HasKeywords(keywords) {
			items = append(items, i)
		}
	}
	return items
}

func (feed *ParsedFeed) HasKeywords(keywords []string) bool {
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
