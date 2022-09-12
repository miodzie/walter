package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/connections/discord"
	"github.com/miodzie/seras/mods"
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
	_ = seras.AddBotParser("discord", &discord.BotParser{})
	err = seras.ParseBots(cfg)
	if err != nil {
		return err
	}

	// Hard code for now.
	bot := seras.Bots["discord"]
	cli(bot)
	bot.AddMods(mods.Default("database.sqlite"))

	return seras.RunBot(bot)
}

func initConfig() (*seras.Config, error) {
	cfg, err := seras.ParseToml(homeDir() + "/.seras/config.toml")
	if err != nil {
		return cfg, err
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

func cli(messenger seras.Messenger) {
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			err := messenger.Send(seras.Message{Content: text})
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
