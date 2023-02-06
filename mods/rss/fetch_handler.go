package rss

type FetchHandler struct {
	fetcher Fetcher
}

// TODO: this would be so much easier, having an already hydrated struct,
// instead of having to fumble about with the Repository more.
type Envelope struct {
	Feed
	Subs []*Subscription
}

func (h *FetchHandler) Handle(urls []string) <-chan Feed {
	output := make(chan Feed, len(urls))
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

func NewFetchHandler(fetcher Fetcher) *FetchHandler {
	return &FetchHandler{fetcher: fetcher}
}
