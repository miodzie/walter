package usecases

import (
	"github.com/miodzie/seras/mods/rss"
	"testing"
)

func TestAddFeed_AddFeed(t *testing.T) {
	repository := rss.NewInMemRepo()
	useCase := NewAddFeed(repository)

	// Act
	response := useCase.AddFeed(AddFeedRequest{})

	// Assert
	if response.Error != nil {
		t.Error(response.Error)
	}
}
