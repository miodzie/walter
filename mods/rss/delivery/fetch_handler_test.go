package delivery

import (
	"errors"
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"github.com/miodzie/walter/mods/rss"
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
	golangBlog := rss.Feed{Title: "Go Blog"}
	s.stubFetcher.Add("blog.golang.org", golangBlog)
	randomBlog := rss.Feed{Title: "Random Blog"}
	s.stubFetcher.Add("localhost", randomBlog)

	output := s.fetchHandler.Handle([]string{"blog.golang.org", "localhost"})

	assert.Cmp(<-output, golangBlog)
	assert.Cmp(<-output, randomBlog)
	assert.Cmp(<-output, rss.Feed{})
}

func (s *FetchHandlerSuite) TestItIgnoresOnFetchError(assert *td.T) {
	s.stubFetcher.Add("localhost", rss.Feed{Title: "Blog"})
	s.stubFetcher.AddErr(errors.New("test"))

	output := s.fetchHandler.Handle([]string{"localhost"})

	assert.Cmp(<-output, rss.Feed{})
}

func TestFetcherTestSuite(t *testing.T) {
	tdsuite.Run(t, new(FetchHandlerSuite))
}

///////////////////////////////////////////////////////////////////////////////////////////

type StubFetcher struct {
	feeds map[string]rss.Feed
	err   error
}

func NewStubFetcher() *StubFetcher {
	return &StubFetcher{feeds: make(map[string]rss.Feed)}
}

func (s *StubFetcher) Fetch(rssUrl string) (*rss.Feed, error) {
	if s.err != nil {
		return nil, s.err
	}
	f, ok := s.feeds[rssUrl]
	if !ok {
		return nil, errors.New("url for parsed userFeed wasn't added")
	}
	delete(s.feeds, rssUrl)
	return &f, nil
}

func (s *StubFetcher) Add(url string, blog rss.Feed) {
	s.feeds[url] = blog
}

func (s *StubFetcher) AddErr(err error) {
	s.err = err
}
