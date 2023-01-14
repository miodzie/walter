package rss

import (
	"errors"
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

type FetchHandlerSuite struct {
	fetchHandler *FetchHandler
	stubFetcher  *StubFetcher
}

func (s *FetchHandlerSuite) PreTest(*td.T, string) error {
	s.stubFetcher = NewStubFetcher()
	s.fetchHandler = NewFetchHandler(s.stubFetcher)
	return nil
}

func (s *FetchHandlerSuite) TestFetcherPipelineReturnsParsedFeeds(assert *td.T) {
	golangBlog := ParsedFeed{Title: "Go Blog"}
	s.stubFetcher.Add("blog.golang.org", golangBlog)
	randomBlog := ParsedFeed{Title: "Random Blog"}
	s.stubFetcher.Add("localhost", randomBlog)

	output := s.fetchHandler.Handle([]string{"blog.golang.org", "localhost"})

	assert.Cmp(<-output, golangBlog)
	assert.Cmp(<-output, randomBlog)
	assert.Cmp(<-output, ParsedFeed{})
}

func (s *FetchHandlerSuite) TestItIgnoresOnFetchError(assert *td.T) {
	s.stubFetcher.Add("localhost", ParsedFeed{Title: "Blog"})
	s.stubFetcher.AddErr(errors.New("test"))

	output := s.fetchHandler.Handle([]string{"localhost"})

	assert.Cmp(<-output, ParsedFeed{})
}

func TestFetcherTestSuite(t *testing.T) {
	tdsuite.Run(t, new(FetchHandlerSuite))
}

///////////////////////////////////////////////////////////////////////////////////////////

type StubFetcher struct {
	feeds map[string]ParsedFeed
	err   error
}

func NewStubFetcher() *StubFetcher {
	return &StubFetcher{feeds: make(map[string]ParsedFeed)}
}

func (s *StubFetcher) Fetch(rssUrl string) (*ParsedFeed, error) {
	if s.err != nil {
		return nil, s.err
	}
	f, ok := s.feeds[rssUrl]
	if !ok {
		return nil, errors.New("url for parsed feed wasn't added")
	}
	delete(s.feeds, rssUrl)
	return &f, nil
}

func (s *StubFetcher) Add(url string, blog ParsedFeed) {
	s.feeds[url] = blog
}

func (s *StubFetcher) AddErr(err error) {
	s.err = err
}
