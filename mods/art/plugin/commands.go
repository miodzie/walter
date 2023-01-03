package plugin

import (
	"time"
	"walter"
	"walter/mods/art"
)

func (mod *VisionaryMod) gmCommand(msg walter.Message) {
	if time.Since(lastRun) < time.Second*2 {
		return
	}
	lastRun = time.Now()
	picture := &art.Picture{Art: art.GM}
	for !picture.Completed() {
		for _, brush := range mod.brushes {
			draw(msg.Target, picture, brush)
		}
	}
}
