package rss

type FetchHandler struct {
	fetcher Fetcher
}

func (h *FetchHandler) Handle(urls []string) <-chan ParsedFeed {
	output := make(chan ParsedFeed, len(urls))
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
