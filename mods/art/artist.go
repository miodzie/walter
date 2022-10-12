package art

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/log"
	"time"
)

type ArtistFactory struct {
}

func (r *ArtistFactory) Create(a interface{}) (seras.Module, error) {
	// We need a visionary to paint his vision!
	// Wait until one is available.
	for visionary == nil {
		log.Warn("No visionary, " +
			"waiting and trying again until they're created.")
		time.Sleep(1 * time.Second)
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
