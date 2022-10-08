package art

import (
	"errors"
	"github.com/miodzie/seras"
)

type ArtistFactory struct {
}

func (r *ArtistFactory) Create(a interface{}) (seras.Module, error) {
	if visionary == nil {
		return nil, errors.New("race condition: restart the bot to try again. sorry")
	}
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
