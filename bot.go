package seras

import (
	"fmt"
	"github.com/miodzie/seras/log"
)

type Bot interface {
	Connection
	Admin
	Modable
	Messenger
	MessageFormatter
}

type Modable interface {
	Mods() []Module
	AddMods([]Module)
}

var Bots map[string]Bot

func ParseBots(cfg *Config) error {
	Bots = make(map[string]Bot)
	for name, con := range cfg.Bots {
		parser, ok := connectors[con["type"].(string)]
		if !ok {
			return fmt.Errorf("unknown parser: %s", name)
		}
		var err error
		Bots[name], err = parser.Parse(con)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunBot(bot Bot) error {
	stream, _ := bot.Connect()
	manager, err := NewModManager(bot.Mods(), bot)
	if err != nil {
		return err
	}

	return manager.Run(stream)
}

func RunAll(addMods func(string) []Module) error {
	errc := make(chan error)
	for name, bot := range Bots {
		log.Info("Starting connection: ", name)
		bot.AddMods(addMods(name))
		go func(bot Bot) {
			errc <- RunBot(bot)
		}(bot)
	}
	return <-errc
}
