// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package fetchers

import (
	"crypto/tls"
	"github.com/miodzie/walter/mods/rss"
	"github.com/mmcdole/gofeed"
	"net/http"
)

type GoFeedFetcher struct{}

func GoFeed() *GoFeedFetcher { return &GoFeedFetcher{} }

func (receiver *GoFeedFetcher) Fetch(url string) (*rss.Feed, error) {
	var parsed rss.Feed
	fp := gofeed.NewParser()
	// If TLSClientConfig is non-nil, HTTP/2 support may not be enabled by default.
	fp.Client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{}}}
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	parsed.Title = feed.Title
	parsed.Description = feed.Description
	parsed.Link = feed.Link
	parsed.Updated = feed.Updated
	parsed.Published = feed.Published
	parsed.Custom = feed.Custom

	for _, i := range feed.Items {
		var pi rss.Item
		pi.Title = i.Title
		pi.Description = i.Description
		pi.Content = i.Content
		pi.Link = i.Link
		pi.Links = i.Links
		pi.GUID = i.GUID
		pi.Custom = i.Custom
		parsed.Items = append(parsed.Items, pi)
	}

	return &parsed, nil
}
