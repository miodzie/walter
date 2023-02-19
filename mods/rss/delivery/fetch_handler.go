package delivery

import "github.com/miodzie/walter/mods/rss"

type FetchHandler struct {
	fetcher rss.Fetcher
}

// TODO: this would be so much easier, having an already hydrated struct,
// instead of having to fumble about with the Repository more.
type Envelope struct {
	rss.Feed
	Subs []*rss.Subscription
}

func (h *FetchHandler) Handle(urls []string) <-chan rss.Feed {
	output := make(chan rss.Feed, len(urls))
	go func() {
		for _, u := range urls {
			feed, err := h.fetcher.Fetch(u)
			if err != nil {
				continue
			}
			output <- *feed
		}
		close(output)
	}()
	return output
}

func NewFetchHandler(fetcher rss.Fetcher) *FetchHandler {
	return &FetchHandler{fetcher: fetcher}
}
