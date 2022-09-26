package mods

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/bestbot"
	"github.com/miodzie/seras/mods/dong"
	"github.com/miodzie/seras/mods/logger"
	"github.com/miodzie/seras/mods/logger/drivers"
	"github.com/miodzie/seras/mods/moderator"
	rss2 "github.com/miodzie/seras/mods/rss"
	"github.com/miodzie/seras/mods/rss/parsers/decorators"
	"github.com/miodzie/seras/mods/rss/parsers/gofeed"
	rss "github.com/miodzie/seras/mods/rss/plugin"
	sed "github.com/miodzie/seras/mods/sed/plugin"
	"github.com/miodzie/seras/storage/sqlite"
)

func Default(dbPath string) []seras.Module {
	db, err := sqlite.Setup(dbPath)
	if err != nil {
		panic(err)
	}

	return []seras.Module{
		sed.New(),
		dong.New(),
		bestbot.New(),
		moderator.New(),
		logger.New(drivers.NewMultiLogger(drivers.ConsoleLogger{})),
		rss.New(rss.Context{
			Repository: sqlite.NewRssRepository(db),
			Parser:     decorators.StripHtml(gofeed.New()),
			Formatter:  rss2.MinimalFormatter{},
		}),
	}
}
