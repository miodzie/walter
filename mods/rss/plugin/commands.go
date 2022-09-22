package plugin

import (
	"fmt"
	"strings"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/rss/usecases"
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

	var addFeed = &usecases.AddFeed{Repo: mod.Repository}
	// TODO: validate.
	req := usecases.AddFeedRequest{
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
	var showFeeds usecases.ShowFeeds

	resp := showFeeds.Handle(mod.Repository)

	if resp.Error != nil {
		mod.actions.Reply(msg, resp.Message)
		fmt.Println("WHAT" + resp.Message)
		fmt.Println(resp.Error)
		return
	}

	// TODO: Presenter layer.
	var reply = seras.Message{Target: msg.Target}
	var parsed []string
	for _, feed := range resp.Feeds {
		parsed = append(parsed, fmt.Sprintf("%s: %s", feed.Name, feed.Url))
	}
	reply.Content = strings.Join(parsed, "\n")
	if len(parsed) == 0 {
		reply.Content = "No feeds available. Ask an admin to add some."
	}
	mod.actions.Send(reply)
	reply.Content = fmt.Sprintf("To subscribe to a feed, use %ssubscribe {name} {keywords}, keywords being comma separated (spaces are ok, e.g. \"spy x family, comedy\")", seras.Token())
	mod.actions.Send(reply)
}

// !subscribe {feed name} {keywords, comma separated}
func (mod *RssMod) subscribe(msg seras.Message) {
	if len(msg.Arguments) < 3 {
		mod.actions.Reply(msg, fmt.Sprintf("To subscribe to a feed, use %ssubscribe {name} {keywords}, keywords being comma separated (spaces are ok, e.g. \"spy x family, comedy\")", seras.Token()))
		return
	}
	// TODO: validate & parse?
	keywords := strings.Join(msg.Arguments[2:], " ")
	req := usecases.SubscribeRequest{
		FeedName: msg.Arguments[1],
		Keywords: keywords,
		Channel:  msg.Target,
		User:     msg.Author.Mention,
	}
	var subscribe = usecases.NewSubscribeUseCase(mod.Repository)
	resp, err := subscribe.Handle(req)
	// TODO: Probably remove err return argument.
	fmt.Println(err)

	mod.actions.Reply(msg, resp.Message)
}

// !unsubscribe {feed name}
func (mod *RssMod) unsubscribe(msg seras.Message) {
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
	fmt.Printf("%+v\n", request)
	uc := usecases.NewUnsubscribeUseCase(mod.Repository)
	response := uc.Handle(request)
	mod.actions.Reply(msg, response.Message)
}

func (mod *RssMod) subs(msg seras.Message) {
	request := usecases.ListSubscriptionsRequest{
		User:     msg.Author.Mention,
		Optional: struct{ Channel string }{msg.Target},
	}

	useCase := usecases.NewListSubscriptionsUseCase(mod.Repository)
	response := useCase.Handle(request)
	if response.Error != nil {
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
