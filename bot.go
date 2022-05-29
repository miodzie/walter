package seras

import "fmt"

type Bot interface {
	Connection
	Admin
	Messenger
	MessageFormatter
	Modable
}

type Modable interface {
	Mods() []Module
	AddMods([]Module)
}

var Bots map[string]Bot

func ParseBots(cfg *Config) error {
	Bots = make(map[string]Bot)
	for name, con := range cfg.Bots {
		parser, ok := connectors[name]
		if !ok {
			return fmt.Errorf("unable to parse connector: %s", name)
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
