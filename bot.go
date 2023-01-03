// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package walter

import (
	"fmt"
	"strings"
	"walter/log"
)

type Bot interface {
	Connection
	Admin
	Modable
	Messenger
	MessageFormatter
}

type Modable interface {
	Modules() []Module
	AddMods([]Module)
	// ModList [name]config
	ModList() map[string]interface{}
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
		// TODO: Refactor.
		Bots[name].SetName(name)
	}

	return nil
}

func RunBot(bot Bot) error {
	stream, _ := bot.Connect()

	var modList []string
	for _, mod := range bot.Modules() {
		modList = append(modList, mod.Name())
	}
	log.Infof("[%s] Modules: %s\n", bot.Name(), strings.Join(modList, ", "))
	manager, err := NewModManager(bot.Modules(), bot)
	if err != nil {
		return err
	}

	return manager.Run(stream)
}

func RunAll(createModsFor func(Bot) []Module) error {
	errc := make(chan error)
	for name, bot := range Bots {
		log.Infof("Starting connection: %s\n", name)
		bot.AddMods(createModsFor(bot))
		go func(bot Bot) {
			errc <- RunBot(bot)
		}(bot)
	}

	return <-errc
}
