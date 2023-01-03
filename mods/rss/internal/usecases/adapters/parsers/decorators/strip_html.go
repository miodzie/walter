// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package decorators

import (
	"github.com/miodzie/walter/mods/rss/internal/usecases"
	"html"
	"reflect"

	"github.com/microcosm-cc/bluemonday"
	"github.com/miodzie/walter/mods/rss/internal/internal/domain"
)

type cleanHtml struct {
	BaseParser usecases.Parser
}

func StripHtml(p usecases.Parser) usecases.Parser {
	return &cleanHtml{BaseParser: p}
}

func (s *cleanHtml) ParseURL(url string) (*domain.ParsedFeed, error) {
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
			f.SetString(html.UnescapeString(p.Sanitize(f.String())))
		}
	}
}
