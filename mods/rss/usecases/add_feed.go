package usecases

import "github.com/miodzie/seras/mods/rss"

type AddFeed struct {
	repository rss.Repository
}

func NewAddFeed(repository rss.Repository) *AddFeed {
	return &AddFeed{repository: repository}
}

type AddFeedRequest struct {
	Name string
	Url  string
}

type AddFeedResponse struct {
	Message string
}

func (a *AddFeed) AddFeed(req AddFeedRequest) (AddFeedResponse, error) {
	var resp AddFeedResponse

	var feed rss.Feed
	feed.Name = req.Name
	feed.Url = req.Url

	err := a.repository.AddFeed(&feed)

	resp.Message = "Feed saved."
	if err != nil {
		resp.Message = "Failed to save feed."
	}

	return resp, err
}
