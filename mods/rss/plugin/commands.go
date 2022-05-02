package plugin

import (
	"fmt"
	"strings"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/rss/usecases"
)

// !add_feed {name} {url}
func (mod *RssMod) addFeed(msg seras.Message) {
	var addFeed = &usecases.AddFeed{Feeds: mod.feeds}
	// TODO: validate.
	req := usecases.AddFeedRequest{
		Name: msg.Arguments[1],
		Url:  msg.Arguments[2],
	}

	resp := addFeed.Handle(req)

	if resp.Error != nil {
		fmt.Println(resp.Error)
	}
}

// !feeds
func (mod *RssMod) showFeeds(msg seras.Message) {
	var showFeeds usecases.ShowFeeds

	resp := showFeeds.Handle(mod.feeds)

	if resp.Error != nil {
		mod.actions.Reply(msg, resp.Message)
		return
	}

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
	req := usecases.SubscribeRequest{
		FeedName: msg.Arguments[1],
		Keywords: keywords,
		Channel:  msg.Channel,
		User:     msg.AuthorId,
	}
	var subscribe = &usecases.Subscribe{
		Feeds: mod.feeds,
		Subs:  mod.subs,
	}
	resp := subscribe.Handle(req)

	if resp.Error != nil {
		fmt.Println(resp.Error)
	}

	mod.actions.Reply(msg, resp.Message)
}
