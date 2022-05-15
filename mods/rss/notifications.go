package rss

type Notification struct {
	Feed    Feed
	Item    Item
	Channel string
	Users   []string
}
