// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package art

import (
	"github.com/miodzie/walter"
)

type ArtistFactory struct {
}

func (r *ArtistFactory) Create(a interface{}) (walter.Module, error) {
	sheep := &ArtistMod{
		instructions: newArtistPalette(),
		running:      false,
	}

	return sheep, nil
}

type ArtistMod struct {
	instructions walter.Stream
	running      bool
}

func (mod *ArtistMod) Name() string {
	return "artist"
}

func (mod *ArtistMod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-mod.instructions
		actions.Send(msg)
	}
	go func() {
		// might as well drain the queue for memory leaks.
		for mod.running {
			_ = <-stream
		}
	}()
	return nil
}

func (mod *ArtistMod) Stop() {
	mod.running = false
}
