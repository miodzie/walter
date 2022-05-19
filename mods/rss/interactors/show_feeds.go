package interactors

import (
	"text/template"

	"github.com/miodzie/seras/mods/rss"
)

type ShowFeeds struct{}

type ShowFeedsResponse struct {
	Feeds   []*rss.Feed // NOTE: Using the domain model is crossing a boundry, violates DIP.
	Message string
	Error   error
}

func (sl *ShowFeeds) Handle(feeds rss.Feeds) ShowFeedsResponse {
	var resp ShowFeedsResponse
	template.New("test")

	resp.Feeds, resp.Error = feeds.All()

	if len(resp.Feeds) == 0 {
		resp.Message = "No feeds available."
	}

	if resp.Error != nil {
		resp.Message = "Failed to fetch feeds."
	}

	return resp
}
