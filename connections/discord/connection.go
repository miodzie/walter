// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/miodzie/walter"
	"sync"
)

type Connection struct {
	session *discordgo.Session
	stream  chan walter.Message
	mods    []walter.Module
	Config
	sync.Mutex
}

func New(config Config) (*Connection, error) {
	disc := &Connection{Config: config}
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return disc, err
	}
	disc.session = session

	return disc, nil
}

func (con *Connection) Connect() (walter.Stream, error) {
	con.Lock()
	defer con.Unlock()

	con.session.AddHandler(con.onMessageCreate)
	con.stream = make(chan walter.Message)
	con.session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageTyping | discordgo.IntentsDirectMessages

	return con.stream, con.session.Open()
}

func (con *Connection) Close() error {
	con.Lock()
	defer con.Unlock()

	close(con.stream)

	return con.session.Close()
}
