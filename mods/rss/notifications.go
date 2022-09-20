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
	t = fmt.Sprintf(t, i.Title, i.DescTruncated(), i.Link, strings.Join(n.Users, ", "))

	return t
}

type Formatter interface {
	Format(Notification) string
}

type DefaultFormatter struct {
}

func (d DefaultFormatter) Format(notification Notification) string {
	return notification.String()
}

type MinimalFormatter struct {
}

func (m MinimalFormatter) Format(n Notification) string {
	i := n.Item
	return fmt.Sprintf(
		"%s - %s : %s",
		i.Title, i.Link, strings.Join(n.Users, ","),
	)
}
