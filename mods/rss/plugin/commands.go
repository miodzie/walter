// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"flag"
	"fmt"
	"github.com/miodzie/walter/log"
	"strings"

	"github.com/miodzie/walter"
	"github.com/miodzie/walter/mods/rss/usecases"
)

// !add_feed {name} {url}
func (mod *RssMod) addFeed(msg walter.Message) {
	if !mod.actions.IsAdmin(msg.Author.Id) {
		mod.actions.Reply(msg, "Only admins can add feeds.")
		return
	}
	if len(msg.Arguments) != 3 {
		mod.actions.Reply(msg, "Incorrect amount of arguments.")
		return
	}

	addFeed := usecases.NewAddFeed(mod.Repository)
	// TODO: validate.
	request := usecases.AddFeedRequest{
		Name: msg.Arguments[1],
		Url:  msg.Arguments[2],
	}

	resp, err := addFeed.Add(request)

	if err != nil {
		log.Error(err)
	}

	mod.actions.Reply(msg, resp.Message)
}

// !feeds
func (mod *RssMod) showFeeds(msg walter.Message) {
	getFeeds := usecases.NewGetFeeds(mod.Repository)
	resp, err := getFeeds.Get()

	if err != nil {
		mod.actions.Reply(msg, resp.Message)
		log.Error(err)
		return
	}

	var reply = walter.Message{Target: msg.Target}
	var parsed []string
	for _, feed := range resp.Feeds {
		parsed = append(parsed, fmt.Sprintf("%s: %s", feed.Name, feed.Url))
	}
	reply.Content = strings.Join(parsed, "\n")
	if len(parsed) == 0 {
		reply.Content = "No feeds available. Ask an admin to add some."
	}
	mod.actions.Send(reply)
	reply.Content = fmt.Sprintf("To subscribe to a feed, use %ssubscribe {name} {keywords}, keywords being comma separated (spaces are ok, e.g. \"spy x family, comedy\")", walter.Token())
	mod.actions.Send(reply)
}

// !subscribe {feed name} {keywords, comma separated}
func (mod *RssMod) subscribe(msg walter.Message) {
	if len(msg.Arguments) < 3 {
		mod.actions.Reply(msg, fmt.Sprintf("To subscribe to a feed, "+
			"use %ssubscribe my_feed -keywords=foo,bar -ignore=baz,buzz "+
			"keywords being comma separated (spaces are ok, e.g. \"spy x family, comedy\")", walter.Token()))
		return
	}
	subCmd := flag.NewFlagSet("subscribe", flag.ContinueOnError)
	keywords := subCmd.String("keywords", "", "keywords")
	ignore := subCmd.String("ignore", "", "ignore")
	if err := subCmd.Parse(msg.Arguments[2:]); err != nil {
		log.Error(err)
		mod.actions.Reply(msg, "Failed to parse !subscribe commands.")
		return
	}

	req := usecases.SubscribeRequest{
		FeedName:    msg.Arguments[1],
		Keywords:    *keywords,
		Channel:     msg.Target,
		User:        msg.Author.Mention,
		IgnoreWords: *ignore,
	}
	var subscribe = usecases.NewSubscribe(mod.Repository)
	resp, err := subscribe.Subscribe(req)
	if err != nil {
		log.Error(err)
	}

	mod.actions.Reply(msg, resp.Message)
}

// !unsubscribe {feed name}
func (mod *RssMod) unsubscribe(msg walter.Message) {
	if len(msg.Arguments) != 2 {
		mod.actions.Reply(msg, "Invalid amount of arguments. !unsubscribe $feedName")
		return
	}
	feedName := msg.Arguments[1]
	request := usecases.UnsubscribeRequest{
		User:     msg.Author.Mention,
		Channel:  msg.Target,
		FeedName: feedName,
	}
	unsubscribe := usecases.NewUnsubscribe(mod.Repository)
	response, err := unsubscribe.Unsubscribe(request)
	if err != nil {
		log.Error(err)
	}
	mod.actions.Reply(msg, response.Message)
}

func (mod *RssMod) subs(msg walter.Message) {
	request := usecases.GetSubscriptionsRequest{
		User:     msg.Author.Mention,
		Optional: struct{ Channel string }{msg.Target},
	}

	getSubs := usecases.NewGetSubscriptions(mod.Repository)
	response, err := getSubs.Get(request)
	if err != nil {
		log.Error(err)
		mod.actions.Reply(msg, "oh noes i brokededz")
		return
	}
	if len(response.Subscriptions) == 0 {
		mod.actions.Reply(msg, "No subscriptions in this channel.")
		return
	}
	var feeds []string
	for _, sub := range response.Subscriptions {
		feeds = append(feeds, sub.Feed)
	}
	reply := fmt.Sprintf("Subscribed to: %s", strings.Join(feeds, ", "))
	mod.actions.Reply(msg, reply)
}
