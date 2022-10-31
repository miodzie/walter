// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"errors"
	"github.com/miodzie/seras/mods/rss"
	"testing"
)

func TestAddFeed_AddFeed(t *testing.T) {
	repository := rss.NewInMemRepo()
	addFeed := NewAddFeed(repository)

	// Act
	response, err := addFeed.Exec(AddFeedRequest{})

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
	addFeed := NewAddFeed(repository)

	// Act
	response, err := addFeed.Exec(AddFeedRequest{})

	// Assert
	if err != expectedErr {
		t.Error(err)
	}
	if response.Message != "Failed to save feed." {
		t.Fail()
	}
}
