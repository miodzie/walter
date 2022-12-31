package plugin

import (
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/mods/art"
	"time"
)

func (mod *VisionaryMod) gmCommand(msg walter.Message) {
	if time.Since(lastRun) < time.Second*2 {
		return
	}
	lastRun = time.Now()
	picture := &art.Picture{Art: art.GM}
	for !picture.Completed() {
		for _, artist := range mod.artists {
			draw(msg, picture, artist)
		}
	}
}
