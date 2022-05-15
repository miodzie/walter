package interactors

import "github.com/miodzie/seras/mods/rss"

type AddFeed struct {
	Feeds rss.Feeds
}

type AddFeedRequest struct {
	Name string
	Url  string
}

type AddFeedResponse struct {
	Message string
	Error   error
}

func (a *AddFeed) Handle(req AddFeedRequest) AddFeedResponse {
	var resp AddFeedResponse

	var listing *rss.Feed
	listing.Name = req.Name
	listing.Url = req.Url

	err := a.Feeds.Add(listing)

	resp.Message = "Feed saved."
	if err != nil {
		resp.Message = "Failed to save feed."
		resp.Error = err
	}

	return resp
}
