package mods

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/bestbot"
	"github.com/miodzie/seras/mods/policing"
)

func Default() []seras.Module {
	return []seras.Module{bestbot.New(), policing.New()}
}
