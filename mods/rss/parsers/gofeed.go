package parsers

import (
	"github.com/miodzie/seras/mods/rss"
	"github.com/mmcdole/gofeed"
)

type GoFeedParser struct {
}

func (receiver *GoFeedParser) Parse(url string) (*rss.ParsedFeed, error) {
	var parsed *rss.ParsedFeed
	fp := gofeed.NewParser()
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
		var pi *rss.Item
		pi.Title = i.Title
		pi.Description = i.Description
		pi.Content = i.Content
		pi.Link = i.Link
		pi.Links = i.Links
		pi.GUID = i.GUID
		pi.Custom = i.Custom
		parsed.Items = append(parsed.Items, pi)
	}

	return parsed, nil
}
