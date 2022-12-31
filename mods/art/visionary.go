// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package art

import (
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/log"
	"time"
)

const MaxLines = 4

var visionary *Visionary
var lastRun time.Time

func newArtistPalette() walter.Stream {
	stream := make(chan walter.Message)

	// We need a visionary to paint his vision!
	// Wait until one is available, asynchronously!?!
	go func() {
		for visionary == nil {
			log.Warn("No visionary, " +
				"waiting and trying again until they're created.")
			time.Sleep(1 * time.Second)
		}
		visionary.artists = append(visionary.artists, stream)
	}()

	return stream
}

type VisionaryFactory struct {
}

func (b *VisionaryFactory) Create(a interface{}) (walter.Module, error) {
	if visionary == nil {
		visionary = &Visionary{
			artists: []chan walter.Message{},
			running: false,
		}
	}

	return visionary, nil
}

type Visionary struct {
	artists []chan walter.Message
	running bool
}

func (mod *Visionary) Name() string {
	return "visionary"
}

func (mod *Visionary) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		// !gm
		msg.Command("gm", mod.gmCommand)
	}
	return nil
}

func (mod *Visionary) gmCommand(msg walter.Message) {
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

func draw(msg walter.Message, art *Picture, artist chan walter.Message) {
	for i := 0; i < MaxLines || art.Completed(); i++ {
		msg.Content = art.NextLine()
		artist <- msg
		time.Sleep(time.Millisecond * 100)
	}
}

func (mod *Visionary) Stop() {
	mod.running = false
}
