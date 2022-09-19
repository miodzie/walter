package interactors

import (
	"github.com/miodzie/seras/mods/rss"
)

type ShowFeeds struct{}

type ShowFeedsResponse struct {
	Feeds   []*rss.Feed // NOTE: Using the domain model is crossing a boundary, violates DIP.
	Message string
	Error   error
}

func (sl *ShowFeeds) Handle(repo rss.Repository) ShowFeedsResponse {
	var resp ShowFeedsResponse

	resp.Feeds, resp.Error = repo.AllFeeds()

	if len(resp.Feeds) == 0 {
		resp.Message = "No feeds available."
	}

	if resp.Error != nil {
		resp.Message = "Failed to fetch feeds."
	}

	return resp
}
