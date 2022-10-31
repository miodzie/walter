// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package art

import (
	"github.com/miodzie/seras"
)

type ArtistFactory struct {
}

func (r *ArtistFactory) Create(a interface{}) (seras.Module, error) {
	sheep := &Artist{
		instructions: newArtistPalette(),
		running:      false,
	}

	return sheep, nil
}

type Artist struct {
	instructions seras.Stream
	running      bool
}

func (mod *Artist) Name() string {
	return "artist"
}

func (mod *Artist) Start(stream seras.Stream, actions seras.Actions) error {
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

func (mod *Artist) Stop() {
	mod.running = false
}
