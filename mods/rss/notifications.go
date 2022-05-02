package rss

type Notification struct {
	Feed    Feed
	Channel string
	Users   []string
}
