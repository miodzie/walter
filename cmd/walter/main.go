// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/miodzie/walter/log"
	"github.com/miodzie/walter/mods"
	"github.com/miodzie/walter/storage"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/miodzie/walter"
	_ "github.com/miodzie/walter/connections/discord"
	_ "github.com/miodzie/walter/connections/irc"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := initConfig()
	if err != nil {
		return err
	}
	interrupt(func() {})
	err = walter.ParseBots(cfg)
	if err != nil {
		return err
	}

	return walter.RunAll(func(bot walter.Bot) []walter.Module {
		//m := mods.Default(fmt.Sprintf("%s.sqlite", bot.Name()))
		m, err := mods.CreateFromList(bot.ModList())
		if err != nil {
			panic(err)
		}

		return m
	})
}

func initConfig() (*walter.Config, error) {
	file := homeDir() + "/.walter/config.toml"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(homeDir()+"/.walter", 0700)
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(file, []byte(walter.DefaultConfig), 0600)
		if err != nil {
			return nil, err
		}
		log.Warnf(`Config was not found, created.
Please update the config located at: %s And restart.\n`, file)
		os.Exit(0)
	}
	cfg, err := walter.ParseToml(file)
	if err != nil {
		return cfg, err
	}

	err = storage.InitFromConfig(file)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func interrupt(callable func()) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		callable()
		os.Exit(1)
	}()
}

func cli(messenger walter.Messenger) {
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			err := messenger.Send(walter.Message{Content: text})
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

}

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
