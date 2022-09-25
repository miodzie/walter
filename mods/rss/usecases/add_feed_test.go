package usecases

import (
	"errors"
	"github.com/miodzie/seras/mods/rss"
	"testing"
)

func TestAddFeed_AddFeed(t *testing.T) {
	repository := rss.NewInMemRepo()
	useCase := NewAddFeed(repository)

	// Act
	response, err := useCase.AddFeed(AddFeedRequest{})

	// Assert
	if err != nil {
		t.Error(err)
	}
	if response.Message != "Feed saved." {
		t.Fail()
	}
}

func TestAddFeed_AddFeed_fails(t *testing.T) {
	repository := rss.NewInMemRepo()
	expectedErr := errors.New("test")
	repository.ForceError(expectedErr, 0)
	useCase := NewAddFeed(repository)

	// Act
	response, err := useCase.AddFeed(AddFeedRequest{})

	// Assert
	if err != expectedErr {
		t.Error(err)
	}
	if response.Message != "Failed to save feed." {
		t.Fail()
	}
}
