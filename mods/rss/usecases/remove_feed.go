package usecases

import (
	"github.com/miodzie/walter/mods/rss"
)

type RemoveFeed struct {
	repository rss.Repository
}

func (f *RemoveFeed) Remove(name string) error {
	// TODO: I should probably have tests for this repository call setup...
	subs, err := f.repository.Subs(rss.SearchParams{FeedName: name})
	if err != nil {
		return err
	}
	if err = f.repository.RemoveFeed(name); err != nil {
		return err
	}
	for _, s := range subs {
		err = f.repository.RemoveSub(s)
	}
	return err
}

func NewRemoveFeed(repo rss.Repository) *RemoveFeed { return &RemoveFeed{repository: repo} }
