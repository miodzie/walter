package rss

import (
	"fmt"
	"strings"

	"github.com/miodzie/seras"
)

// !add_feed {name} {url}
func (mod *RssMod) addFeed(msg seras.Message) {
	// TODO: validate.
	feed := &Feed{
		Name: msg.Arguments[1],
		Url:  msg.Arguments[2],
	}
	fmt.Println(feed.Name, feed.Url)
	err := mod.feeds.Save(feed)
	if err != nil {
		fmt.Println(err)
	}
}

// !subscribe {feed name} {keywords, comma separated}
func (mod *RssMod) subscribe(msg seras.Message) {
	feed, err := mod.feeds.GetByName(msg.Arguments[1])
	if err != nil {
		mod.actions.Send(seras.Message{Content: "Unknown feed."})
		fmt.Println(err)
		return
	}
	// TODO: parse, test
	keywords := strings.Join(msg.Arguments[2:], " ")
	fmt.Println(keywords)
	sub := &Subscription{
		FeedId:   feed.Id,
		Channel:  msg.Channel,
		Keywords: keywords,
		User:     msg.AuthorId,
	}
	err = mod.subs.Save(sub)
	if err != nil {
		fmt.Println(err)
		mod.actions.Send(seras.Message{Content: "Failed to save feed, likely one already exists for this channel and feed."})
	}
}
