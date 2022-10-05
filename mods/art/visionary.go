package art

import (
	"errors"
	"github.com/miodzie/seras"
)

var visionary *Visionary

func newArtistPalette() seras.Stream {
	stream := make(chan seras.Message)
	visionary.artists = append(visionary.artists, stream)

	return stream
}

type VisionaryFactory struct {
}

func (b *VisionaryFactory) Create(a interface{}) (seras.Module, error) {
	if visionary != nil {
		return nil, errors.New("there can only be one")
	}
	visionary = &Visionary{
		artists: []chan seras.Message{},
		running: false,
	}

	return visionary, nil
}

type Visionary struct {
	artists []chan seras.Message
	running bool
}

func (mod *Visionary) Name() string {
	return "visionary"
}

func (mod *Visionary) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if msg.IsCommand("gm") {
			msg.Content = "g'mrn frens"
			for _, a := range mod.artists {
				// Blocking, ofc.
				a <- msg
			}
		}
	}
	return nil
}

func (mod *Visionary) Stop() {
	mod.running = false
}
