package plugin

import (
	"time"
	"walter"
	"walter/log"
)

type Brush <-chan walter.Message

type ArtistMod struct {
	brush   Brush
	running bool
}

func (mod *ArtistMod) Name() string { return "artist" }
func (mod *ArtistMod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-mod.brush
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
func (mod *ArtistMod) Stop() { mod.running = false }

type ArtistFactory struct{}

func (r *ArtistFactory) Create(any) (walter.Module, error) {
	sheep := &ArtistMod{
		brush:   newArtistPalette(),
		running: false,
	}

	return sheep, nil
}
func newArtistPalette() Brush {
	stream := make(chan walter.Message)
	go func() {
		for visionary == nil {
			log.Warn("no visionary, waiting and trying again until they're created")
			time.Sleep(1 * time.Second)
		}
		visionary.brushes = append(visionary.brushes, stream)
	}()
	return stream
}
