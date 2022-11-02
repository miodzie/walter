// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"github.com/miodzie/walter/log"
	"time"

	"github.com/miodzie/walter"
	"github.com/miodzie/walter/mods/rss"
)

type RssMod struct {
	actions walter.Actions
	running bool
	Context
}

type Context struct {
	rss.Repository
	rss.Parser
	rss.Formatter
}

func New(ctx Context) *RssMod {
	return &RssMod{Context: ctx}
}
func (mod *RssMod) Name() string {
	return "rss"
}

func (mod *RssMod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	mod.actions = actions
	go mod.checkFeeds()
	for mod.running {
		msg := <-stream
		msg.Command("feeds", mod.showFeeds)
		msg.Command("add_feed", mod.addFeed)
		msg.Command("subscribe", mod.subscribe)
		msg.Command("unsubscribe", mod.unsubscribe)
		msg.Command("subs", mod.subs)
	}

	return nil
}

func (mod *RssMod) checkFeeds() {
	time.Sleep(time.Minute * 1)
	p := rss.NewProcessor(mod.Repository, mod.Parser)
	for mod.running {
		log.Info("Processing feed subscriptions...")
		notifs, err := p.Process()
		if err != nil {
			log.Error(err)
		}
		log.Infof("%d notifications found\n", len(notifs))
		for _, notif := range notifs {
			msg := walter.Message{
				Target:  notif.Channel,
				Content: mod.Format(*notif),
			}
			log.Debug(notif)
			mod.actions.Send(msg)
		}
		time.Sleep(time.Minute * 30)
	}
}

func (mod *RssMod) Stop() {
	mod.running = false
}
