// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"fmt"
	"walter"
	"walter/mods/sed"
)

func New() *Mod {
	return &Mod{}
}

type Mod struct {
	log     map[string][]walter.Message
	running bool
}

func (mod *Mod) Name() string {
	return "sed"
}

func (mod *Mod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	mod.log = make(map[string][]walter.Message)
	for mod.running {
		msg := <-stream
		// TODO: Parse user input better.
		s := sed.ParseSed(msg.Content)
		if s.Command != ".s" {
			mod.logMsg(msg)
		}
		// NOTE: Technically, tokens can be anything, not just "/". Scan for multiple tokens?
		if s.Command == ".s" || s.Command == "s" {
			for i := len(mod.log[msg.Target]) - 1; i >= 0; i-- {
				m := mod.log[msg.Target][i]
				if s.HasMatch(m.Content) {
					r := fmt.Sprintf("%s meant to say: %s", m.Author.Nick, s.Replace(m.Content))
					actions.Reply(msg, r)
					break
				}
			}
		}
	}

	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}

func (mod *Mod) logMsg(msg walter.Message) {
	mod.log[msg.Target] = append(mod.log[msg.Target], msg)
	if len(mod.log[msg.Target]) > 20 {
		mod.log[msg.Target] = mod.log[msg.Target][1:]
	}
}

type ModFactory struct {
}

func (m ModFactory) Create(config interface{}) (walter.Module, error) {
	return New(), nil
}
