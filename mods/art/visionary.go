package art

import (
	"github.com/miodzie/seras"
	"time"
)

const MaxLines = 4

var visionary *Visionary
var lastRun time.Time

func newArtistPalette() seras.Stream {
	stream := make(chan seras.Message)
	visionary.artists = append(visionary.artists, stream)

	return stream
}

type VisionaryFactory struct {
}

func (b *VisionaryFactory) Create(a interface{}) (seras.Module, error) {
	if visionary == nil {
		visionary = &Visionary{
			artists: []chan seras.Message{},
			running: false,
		}
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
		// !gm
		msg.Command("gm", mod.gmCommand)
	}
	return nil
}

func (mod *Visionary) gmCommand(msg seras.Message) {
	// Quick throttle impl
	if time.Since(lastRun) < time.Second*2 {
		return
	}
	lastRun = time.Now()
	art := &Picture{Art: gm}
	for !art.Completed() {
		for _, artist := range mod.artists {
			draw(msg, art, artist)
		}
	}
}

func draw(msg seras.Message, art *Picture, artist chan seras.Message) {
	for i := 0; i < MaxLines || art.Completed(); i++ {
		msg.Content = art.NextLine()
		artist <- msg
		time.Sleep(time.Millisecond * 100)
	}
}

func (mod *Visionary) Stop() {
	mod.running = false
}
