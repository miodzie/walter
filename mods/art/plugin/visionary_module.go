package plugin

import (
	"time"
	"walter"
	"walter/mods/art"
)

const MaxLines = 4

var visionary *VisionaryMod
var lastRun time.Time

type VisionaryMod struct {
	brushes []chan walter.Message
	running bool
}

func (mod *VisionaryMod) Name() string { return "visionary" }

func (mod *VisionaryMod) Start(stream walter.Stream, _ walter.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		msg.Command("gm", mod.gmCommand)
	}
	return nil
}

func draw(target string, picture *art.Picture, artist chan walter.Message) {
	for i := 0; i < MaxLines || picture.Completed(); i++ {
		artist <- walter.Message{
			Content: picture.NextLine(),
			Target:  target}
		time.Sleep(time.Millisecond * 100)
	}
}

func (mod *VisionaryMod) Stop() { mod.running = false }

type VisionaryFactory struct{}

func (b *VisionaryFactory) Create(any) (walter.Module, error) {
	if visionary == nil {
		visionary = &VisionaryMod{
			brushes: []chan walter.Message{},
			running: false,
		}
	}

	return visionary, nil
}
