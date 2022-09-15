package mods

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/bestbot"
	"github.com/miodzie/seras/mods/policing"
	"github.com/miodzie/seras/mods/rss/parsers/decorators"
	"github.com/miodzie/seras/mods/rss/parsers/gofeed"
	rss "github.com/miodzie/seras/mods/rss/plugin"
	sed "github.com/miodzie/seras/mods/sed/plugin"
	"github.com/miodzie/seras/storage/sqlite"
)

func Default(db string) []seras.Module {
	err := sqlite.Setup(db)
	if err != nil {
		panic(err)
	}

	return []seras.Module{
		sed.New(),
		bestbot.New(),
		policing.New(),
		rss.New(
			rss.Context{
				Repository: &sqlite.RssRepository{},
				Parser:     decorators.StripHtml(gofeed.New()),
				Formatter:  nil,
			},
		),
	}
}
