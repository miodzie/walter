package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/miodzie/seras/log"
	"github.com/miodzie/seras/mods"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/miodzie/seras"
	_ "github.com/miodzie/seras/connections/discord"
	_ "github.com/miodzie/seras/connections/irc"
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
	err = seras.ParseBots(cfg)
	if err != nil {
		return err
	}

	return seras.RunAll(func(bot seras.Bot) []seras.Module {
		//m := mods.Default(fmt.Sprintf("%s.sqlite", bot.Name()))
		m, err := mods.CreateFromList(bot.ModList())
		if err != nil {
			panic(err)
		}

		return m
	})
}

func initConfig() (*seras.Config, error) {
	file := homeDir() + "/.seras/config.toml"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(homeDir()+"/.seras", 0700)
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(file, []byte(seras.DefaultConfig), 0600)
		if err != nil {
			return nil, err
		}
		log.Warnf("Config was not found, created.\nPlease update the config located at: %s\nAnd restart.", file)
		os.Exit(0)
	}
	cfg, err := seras.ParseToml(file)
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
