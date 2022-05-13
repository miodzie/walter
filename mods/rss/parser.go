package rss

type DefaultParser struct {
}

func (p *DefaultParser) Parse(feed Feed) (ParsedFeed, error) {
	var parsed ParsedFeed
	return parsed, nil
}
