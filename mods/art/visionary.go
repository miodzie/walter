package art

import (
	"errors"
	"github.com/miodzie/seras"
	"time"
)

const MaxLines = 4

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
			art := &Picture{Art: gm}
			for !art.Completed() {
				for _, artist := range mod.artists {
					if art.Completed() {
						break
					}
					for i := 0; i < MaxLines; i++ {
						msg.Content = art.NextLine()
						artist <- msg
						time.Sleep(time.Millisecond * 50)
						if art.Completed() {
							break
						}
					}
				}
			}
		}
	}
	return nil
}

func (mod *Visionary) Stop() {
	mod.running = false
}
