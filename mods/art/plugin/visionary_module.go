package plugin

import (
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/mods/art"
	"time"
)

const MaxLines = 4

var visionary *VisionaryMod
var lastRun time.Time

type VisionaryMod struct {
	artists []chan walter.Message
	running bool
}

func (mod *VisionaryMod) Name() string {
	return "visionary"
}

func (mod *VisionaryMod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		msg.Command("gm", mod.gmCommand)
	}
	return nil
}

func draw(msg walter.Message, picture *art.Picture, artist chan walter.Message) {
	for i := 0; i < MaxLines || picture.Completed(); i++ {
		msg.Content = picture.NextLine()
		artist <- msg
		time.Sleep(time.Millisecond * 100)
	}
}

func (mod *VisionaryMod) Stop() {
	mod.running = false
}

type VisionaryFactory struct{}

func (b *VisionaryFactory) Create(any) (walter.Module, error) {
	if visionary == nil {
		visionary = &VisionaryMod{
			artists: []chan walter.Message{},
			running: false,
		}
	}

	return visionary, nil
}
