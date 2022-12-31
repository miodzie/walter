package plugin

import (
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/log"
	"time"
)

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
		// Drain channel for memory leaks.
		for mod.running {
			_ = <-stream
		}
	}()
	return nil
}

func (mod *ArtistMod) Stop() {
	mod.running = false
}

type ArtistFactory struct{}

func (r *ArtistFactory) Create(any) (walter.Module, error) {
	sheep := &ArtistMod{
		instructions: newArtistPalette(),
		running:      false,
	}

	return sheep, nil
}
func newArtistPalette() walter.Stream {
	stream := make(chan walter.Message)
	go func() {
		for visionary == nil {
			log.Warn("no visionary, waiting and trying again until they're created")
			time.Sleep(1 * time.Second)
		}
		visionary.artists = append(visionary.artists, stream)
	}()
	return stream
}
