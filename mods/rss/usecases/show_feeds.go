package usecases

import "github.com/miodzie/seras/mods/rss"

type ShowFeeds struct {
}

type ShowFeedsResponse struct {
	Feeds   []rss.Feed
	Message string
	Error   error
}

func (sf *ShowFeeds) Handle(feeds rss.FeedRepository) ShowFeedsResponse {
	var resp ShowFeedsResponse

	resp.Feeds, resp.Error = feeds.All()

	if len(resp.Feeds) == 0 {
		resp.Message = "No feeds available."
	}

    if resp.Error != nil {
        resp.Message = "Failed to fetch feeds."
    }

	return resp
}
