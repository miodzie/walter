package mods

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/bestbot"
	"github.com/miodzie/seras/mods/policing"
	"github.com/miodzie/seras/mods/rss"
	"github.com/miodzie/seras/storage/memory"
)

func Default() []seras.Module {
	return []seras.Module{
	    bestbot.New(),
	    policing.New(),
	    rss.New(&memory.InMemListenerRepository{}),
	}
}
