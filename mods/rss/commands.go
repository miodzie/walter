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
	err := mod.feeds.Save(feed)
	if err != nil {
		fmt.Println(err)
	}
}

// !subscribe {feed name} {keywords, comma separated}
func (mod *RssMod) subscribe(msg seras.Message) {
	feed, err := mod.feeds.GetByName(msg.Arguments[1])
	if err != nil {
		panic(err)
	}
	// TODO: write a test damn it
	keywords := strings.Join(msg.Arguments[2:], "")
	sub := &Subscription{
		Feed:     feed,
		Channel:  msg.Channel,
		Keywords: strings.Split(keywords, ","),
		User:     msg.Author,
	}
	err = mod.subs.Save(sub)
	if err != nil {
		panic(err)
	}
}
