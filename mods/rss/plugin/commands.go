package plugin

import (
	"fmt"
	"strings"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/rss/interactors"
)

// !add_feed {name} {url}
func (mod *RssMod) addFeed(msg seras.Message) {
	if !mod.actions.IsAdmin(msg.Author.Id) {
		mod.actions.Reply(msg, "Only admins can add feeds.")
		return
	}
	if len(msg.Arguments) != 3 {
		mod.actions.Reply(msg, "Incorrect amount of arguments.")
		return
	}

	var addFeed = &interactors.AddFeed{Feeds: mod.Feeds}
	// TODO: validate.
	req := interactors.AddFeedRequest{
		Name: msg.Arguments[1],
		Url:  msg.Arguments[2],
	}

	resp := addFeed.Handle(req)

	if resp.Error != nil {
		fmt.Println(resp.Error)
	}

	mod.actions.Reply(msg, resp.Message)
}

// !feeds
func (mod *RssMod) showFeeds(msg seras.Message) {
	var showFeeds interactors.ShowFeeds

	resp := showFeeds.Handle(mod.Feeds)

	if resp.Error != nil {
		mod.actions.Reply(msg, resp.Message)
		fmt.Println(resp.Error)
		return
	}

	// TODO: Presenter layer.
	var reply = seras.Message{Channel: msg.Channel}
	var parsed []string
	for _, feed := range resp.Feeds {
		parsed = append(parsed, fmt.Sprintf("%s: %s", feed.Name, feed.Url))
	}
	reply.Content = strings.Join(parsed, "\n")
	reply.Content += fmt.Sprintf("\n\nTo subscribe to a feed, use %ssubscribe {name} {keywords}, keywords being comma separated (spaces are ok, e.g. \"spy x family, comedy\")", seras.Token())

	mod.actions.Send(reply)
}

// !subscribe {feed name} {keywords, comma separated}
func (mod *RssMod) subscribe(msg seras.Message) {
	if len(msg.Arguments) < 3 {
		return
	}
	// TODO: validate & parse?
	keywords := strings.Join(msg.Arguments[2:], " ")
	req := interactors.SubscribeRequest{
		FeedName: msg.Arguments[1],
		Keywords: keywords,
		Channel:  msg.Channel,
		User:     msg.Author.Mention,
	}
	var subscribe = &interactors.Subscribe{
		Feeds: mod.Feeds,
		Subs:  mod.Subscriptions,
	}
	resp, err := subscribe.Handle(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	mod.actions.Reply(msg, resp.Message)
}
