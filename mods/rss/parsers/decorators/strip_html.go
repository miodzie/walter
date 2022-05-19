package decorators

import (
	"reflect"

	"github.com/microcosm-cc/bluemonday"
	"github.com/miodzie/seras/mods/rss"
)

type cleanHtml struct {
	BaseParser rss.Parser
}

func StripHtml(p rss.Parser) rss.Parser {
	return &cleanHtml{BaseParser: p}
}

func (s *cleanHtml) ParseURL(url string) (*rss.ParsedFeed, error) {
	feed, err := s.BaseParser.ParseURL(url)
	if err != nil {
		return feed, err
	}

	stripHtml(feed)
	for _, i := range feed.Items {
		stripHtml(i)
	}

	return feed, nil
}

func stripHtml(any interface{}) {
	p := bluemonday.StripTagsPolicy()
	r := reflect.ValueOf(any).Elem()
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		if f.Kind() == reflect.String {
			f.SetString(p.Sanitize(f.String()))
		}
	}
}
