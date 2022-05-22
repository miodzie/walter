package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/connections/discord"
	"github.com/miodzie/seras/connections/fake"
	"github.com/miodzie/seras/connections/irc"
	"github.com/miodzie/seras/mods"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := initConfig()
	if err != nil {
		return err
	}
	interupt(func() {})
	seras.AddBotParser("discord", &discord.BotParser{})
	err = seras.ParseBots(cfg)
	if err != nil {
		return err
	}

	// Hard code for now.
	bot := seras.Bots["discord"]
	startCli(bot)
	bot.AddMods(mods.Default("database.sqlite"))

	return seras.RunBot(bot)
}

func interupt(callable func()) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		callable()
		os.Exit(1)
	}()
}

func startCli(messenger seras.Messenger) {
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

func initConfig() (*seras.Config, error) {
	cfg, err := seras.ParseToml(UserHomeDir() + "/.seras/config.toml")
	if err != nil {
		return cfg, err
	}
	err = godotenv.Load()
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func makeIrc() (*irc.Connection, error) {
	return irc.New(irc.Config{Server: "irc.libera.chat:6667", Nick: "Guest", Username: "Guest"})
}

func makeFake() *fake.Connection {
	return fake.NewConnection()
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
