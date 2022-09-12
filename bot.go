package seras

import (
	"fmt"
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

func RunAll() error {
	for name, bot := range Bots {
		//// TODO: Change from defaults to configured.
		// TODO: Change to use configured storage type.
		//bot.AddMods(mods.Default(fmt.Sprintf("%s.sqlite", name)))
		fmt.Printf("Starting %s\n", name)
		// TODO: I'm pretty sure this is blocking lel
		err := RunBot(bot)
		if err != nil {
			return err
		}
	}
	return nil
}
