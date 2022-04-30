package mods

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/bestbot"
	"github.com/miodzie/seras/mods/policing"
	"github.com/miodzie/seras/mods/rss"
	"github.com/miodzie/seras/storage/sqlite"
)

func Default() []seras.Module {
	err := sqlite.Setup("database.sqlite")
	if err != nil {
		panic(err)
	}

	return []seras.Module{
		bestbot.New(),
		policing.New(),
		rss.New(&sqlite.FeedRepository{}, &sqlite.SubscriptionRepository{}),
	}
}