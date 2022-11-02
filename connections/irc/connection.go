// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package irc

import (
	"crypto/tls"
	"github.com/miodzie/walter/connections/irc/plugin"
	"strings"
	"time"

	"github.com/miodzie/walter"
	irc "github.com/thoj/go-ircevent"
	"sync"
)

type Connection struct {
	irc  *irc.Connection
	mods []walter.Module
	Config
	sync.Mutex
}

func New(conf Config) (*Connection, error) {
	ircCon := irc.IRC(conf.Nick, conf.Username)
	ircCon.UseTLS = true
	ircCon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	ircCon.UseSASL = conf.SASL
	// TODO: Update to use internal logger or NilledLogger
	//ircCon.Log = nil
	ircCon.SASLLogin = conf.SASLUsername
	ircCon.SASLPassword = conf.SASLPassword
	con := &Connection{
		irc:    ircCon,
		Config: conf,
		mods:   []walter.Module{plugin.New(ircCon)},
	}

	return con, nil
}

func (con *Connection) Connect() (walter.Stream, error) {
	con.Lock()
	defer con.Unlock()
	err := con.irc.Connect(con.Server)
	if err != nil {
		return nil, err
	}
	stream := make(chan walter.Message)

	con.irc.AddCallback("*", func(event *irc.Event) {
		var channel string
		var args = event.Arguments
		if event.Code == "PRIVMSG" {
			args = strings.Split(event.Arguments[1], " ")
			channel = event.Arguments[0]
		}
		stream <- walter.Message{
			Content:   event.Message(),
			Arguments: args,
			Target:    channel,
			Author: walter.Author{
				Id:      event.Host,
				Nick:    event.Nick,
				Mention: event.Nick,
			},
			Code:           event.Code,
			ConnectionName: con.Name(),
			Raw:            event.Raw,
			Timestamp:      time.Now(),
		}
	})

	go func() {
		con.irc.Loop()
	}()

	return stream, nil
}

func (con *Connection) Close() error {
	con.Lock()
	defer con.Unlock()
	con.irc.Disconnect()
	con.irc.ClearCallback("*")

	return nil
}
