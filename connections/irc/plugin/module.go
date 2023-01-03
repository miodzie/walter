// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Package plugin TODO: Consider moving this to the mods/ package?
package plugin

import (
	irc "github.com/thoj/go-ircevent"
	"walter"
)

type Mod struct {
	irc      *irc.Connection
	running  bool
	channels []string
}

func New(irc *irc.Connection, channels []string) *Mod {
	return &Mod{irc: irc, channels: channels}
}

func (mod *Mod) Name() string {
	return "irc"
}

func (mod *Mod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if msg.Code == "INVITE" && actions.IsAdmin(msg.Author.Id) {
			mod.irc.Join(msg.Arguments[1])
		}
		if msg.IsCommand("join") && actions.IsAdmin(msg.Author.Id) {
			mod.irc.Join(msg.Arguments[1])
		}
		if msg.IsCommand("part") && actions.IsAdmin(msg.Author.Id) {
			if len(msg.Arguments) > 1 {
				mod.irc.Part(msg.Arguments[1])
			} else {
				mod.irc.Part(msg.Target)
			}
		}
	}
	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}
