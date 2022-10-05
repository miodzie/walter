package mods

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/bestbot"
	"github.com/miodzie/seras/mods/botnet"
	"github.com/miodzie/seras/mods/dong"
	"github.com/miodzie/seras/mods/logger"
	"github.com/miodzie/seras/mods/logger/drivers"
	"github.com/miodzie/seras/mods/rss"
	"github.com/miodzie/seras/mods/rss/parsers/decorators"
	"github.com/miodzie/seras/mods/rss/parsers/gofeed"
	rss_plugin "github.com/miodzie/seras/mods/rss/plugin"
	sed "github.com/miodzie/seras/mods/sed/plugin"
	"github.com/miodzie/seras/storage/sqlite"
)

func init() {
	factories = make(map[string]Factory)
	Register("shepherd", &botnet.ShepherdFactory{})
	Register("sheep", &botnet.SheepFactory{})
	Register("best_bot", &bestbot.ModFactory{})
	Register("dong", &dong.ModFactory{})
	Register("sed", &sed.ModFactory{})
	Register("logger", &logger.ModFactory{DefaultLogger: drivers.NewMultiLogger(drivers.ConsoleLogger{})})
	//Register("sed", &rss_plugin.ModFactory{
	//	Context: rss_plugin.Context{
	//		Repository: sqlite.NewRssRepository(db),
	//		Parser:     decorators.StripHtml(gofeed.New()),
	//		Formatter:  rss.MinimalFormatter{},
	//	},
	//})
}

// List [mod name]config
type List map[string]interface{}
type Config struct {
}

func Create(list List) ([]seras.Module, error) {
	var mods []seras.Module

	for name, config := range list {
		mod, err := Make(name, config)
		if err != nil {
			return nil, err
		}
		mods = append(mods, mod)
	}

	return mods, nil
}

func Default(dbPath string) []seras.Module {
	db, err := sqlite.Setup(dbPath)
	if err != nil {
		panic(err)
	}

	return []seras.Module{
		rss_plugin.New(rss_plugin.Context{
			Repository: sqlite.NewRssRepository(db),
			Parser:     decorators.StripHtml(gofeed.New()),
			Formatter:  rss.MinimalFormatter{},
		}),
	}
}
