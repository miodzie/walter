// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package moderator

import (
	"time"
	"walter/log"

	"walter"
)

type Mod struct {
	running bool
}

func New() *Mod {
	return &Mod{}
}

func (mod *Mod) Name() string {
	return "moderator"
}

func (mod *Mod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if IsSpam(msg) {
			err := actions.TimeoutUser(msg.Target, msg.Author.Id, time.Now().Add(time.Minute*1))
			if err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}
