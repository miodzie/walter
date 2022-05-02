package usecases

import "github.com/miodzie/seras/mods/rss"

type AddFeed struct {
	Feeds rss.FeedRepository
}

type AddFeedRequest struct {
	Name string
	Url  string
}

type AddFeedResponse struct {
	Message string
	Error   error
}

func (addFeed *AddFeed) Handle(req AddFeedRequest) AddFeedResponse {
	var resp AddFeedResponse

	var feed *rss.Feed
	feed.Name = req.Name
	feed.Url = req.Url

	err := addFeed.Feeds.Save(feed)

	resp.Message = "Feed saved."
	if err != nil {
		resp.Message = "Failed to save feed."
		resp.Error = err
	}

	return resp
}
