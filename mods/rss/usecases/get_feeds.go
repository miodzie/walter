package usecases

import (
	"github.com/miodzie/seras/mods/rss"
)

type GetFeeds struct{ repository rss.Repository }

func NewGetFeeds(repository rss.Repository) *GetFeeds {
	return &GetFeeds{repository: repository}
}

type GetFeedsResponse struct {
	Feeds   []*rss.Feed // NOTE: Using the domain model is crossing a boundary, violates DIP.
	Message string
}

func (g *GetFeeds) Get() (GetFeedsResponse, error) {
	var err error
	var resp GetFeedsResponse

	resp.Feeds, err = g.repository.Feeds()

	if len(resp.Feeds) == 0 {
		resp.Message = "No feeds available."
	}

	if err != nil {
		resp.Message = "Failed to fetch feeds."
	}

	return resp, err
}
