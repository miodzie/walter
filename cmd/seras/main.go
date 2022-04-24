package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/connections/discord"
	"github.com/miodzie/seras/connections/fake"
	"github.com/miodzie/seras/connections/irc"
	"github.com/miodzie/seras/mods"
	"github.com/miodzie/seras/mods/policing"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	interupt(func() {})
	connection := makeDiscord(os.Getenv("DISCORD_TOKEN"))
	messenger := connection
	// messenger := &seras.NullMessenger{}
	stream, _ := connection.Connect()

  // TODO: Abstract to ModuleManager? 
	modules := []seras.Module{mods.NewBestBotMod(), policing.NewPolicingMod()}
	modStreams := []chan seras.Message{}
	for _, mod := range modules {
		modStream := make(chan seras.Message)
		modStreams = append(modStreams, modStream)
		mod.Loop(modStream, messenger)
	}
	cli(messenger)

	for msg := range stream {
		for _, ch := range modStreams {
			ch <- msg
		}
	}

	return nil
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

func makeIrc() (*irc.Connection, error) {
	return irc.New(irc.Config{Host: "irc.libera.chat:6667", Nick: "Guest", Username: "Guest"})
}

func makeFake() *fake.Connection {
	return fake.NewConnection()
}

func makeDiscord(token string) *discord.Connection {
	disc := discord.New(token)
	return disc
}

func makeDiscordMessenger() {

}
