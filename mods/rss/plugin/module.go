// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"github.com/miodzie/walter/log"
	"github.com/miodzie/walter/mods/rss/delivery"
	"time"

	"github.com/miodzie/walter"
	"github.com/miodzie/walter/mods/rss"
)

type RssMod struct {
	actions walter.Actions
	running bool
	Services
}

type Services struct {
	rss.Repository
	rss.Fetcher
	delivery.Formatter
}

func (mod *RssMod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	mod.actions = actions
	go mod.checkFeeds()
	for mod.running {
		msg := <-stream
		msg.Command("feeds", mod.showFeeds)
		msg.Command("add_feed", mod.addFeed)
		msg.Command("remove_feed", mod.removeFeed)
		msg.Command("subscribe", mod.subscribe)
		msg.Command("unsubscribe", mod.unsubscribe)
		msg.Command("subs", mod.subs)
	}

	return nil
}

func (mod *RssMod) checkFeeds() {
	processor := delivery.NewProcessor(mod.Fetcher, mod.Repository)
	for mod.running {
		log.Info("Processing feed subscriptions...")
		deliveries, err := processor.Process()
		if err != nil {
			log.Error(err)
			return
		}
		deliveries = delivery.ThrottleByChannel(deliveries, 3)
		total := 0
		for mail := range deliveries {
			mail.Deliver(func(address string, content string) error {
				msg := walter.Message{Target: address, Content: content}
				if err := mod.actions.Send(msg); err != nil {
					log.Error(err)
					return err
				}
				total++
				return nil
			})
		}
		log.Infof("%d notifications delivered\n", total)
		time.Sleep(time.Minute * 30)
	}
}
func New(services Services) *RssMod { return &RssMod{Services: services} }
func (mod *RssMod) Name() string    { return "rss" }
func (mod *RssMod) Stop()           { mod.running = false }
