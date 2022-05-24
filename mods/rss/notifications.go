package rss

import (
	"fmt"
	"strings"
)

type Notification struct {
	Feed    Feed
	Item    Item
	Channel string
	Users   []string
}

func (n Notification) String() string {
	i := n.Item
	t := "%s\n%s\n%s\n%s\n"
	t = fmt.Sprintf(t, i.Title, i.Desc(), i.Link, strings.Join(n.Users, ", "))
	return t
}
