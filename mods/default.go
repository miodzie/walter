package mods

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/bestbot"
	"github.com/miodzie/seras/mods/dong"
	"github.com/miodzie/seras/mods/policing"
	rss2 "github.com/miodzie/seras/mods/rss"
	"github.com/miodzie/seras/mods/rss/parsers/decorators"
	"github.com/miodzie/seras/mods/rss/parsers/gofeed"
	rss "github.com/miodzie/seras/mods/rss/plugin"
	sed "github.com/miodzie/seras/mods/sed/plugin"
	"github.com/miodzie/seras/storage/sqlite"
)

func Default(dbPath string) []seras.Module {
	err := sqlite.Setup(dbPath)
	if err != nil {
		panic(err)
	}

	return []seras.Module{
		sed.New(),
		dong.New(),
		bestbot.New(),
		policing.New(),
		rss.New(
			rss.Context{
				Repository: &sqlite.RssRepository{},
				Parser:     decorators.StripHtml(gofeed.New()),
				Formatter:  rss2.MinimalFormatter{},
			},
		),
	}
}
