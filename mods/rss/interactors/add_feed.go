package interactors

import "github.com/miodzie/seras/mods/rss"

type AddFeed struct {
	Repo rss.Repository
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

	var feed rss.Feed
	feed.Name = req.Name
	feed.Url = req.Url

	err := a.Repo.AddFeed(&feed)

	resp.Message = "Feed saved."
	if err != nil {
		resp.Message = "Failed to save feed."
		resp.Error = err
	}

	return resp
}
