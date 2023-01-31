// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package decorators

import (
	"html"
	"reflect"

	"github.com/microcosm-cc/bluemonday"
	"github.com/miodzie/walter/mods/rss"
)

type cleanHtml struct {
	BaseFetcher rss.Fetcher
}

func StripHtml(p rss.Fetcher) rss.Fetcher {
	return &cleanHtml{BaseFetcher: p}
}

func (s *cleanHtml) Fetch(url string) (*rss.Feed, error) {
	feed, err := s.BaseFetcher.Fetch(url)
	if err != nil {
		return feed, err
	}

	stripHtml(feed)
	for i, item := range feed.Items {
		stripHtml(&item)
		feed.Items[i] = item
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
