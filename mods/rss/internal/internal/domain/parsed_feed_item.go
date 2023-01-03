package domain

import (
	"github.com/miodzie/walter/log"
	"regexp"
	"strings"
)

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