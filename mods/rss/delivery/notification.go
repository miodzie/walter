// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package delivery

import (
	"fmt"
	"github.com/miodzie/walter/log"
	"github.com/miodzie/walter/mods/rss"
)

type Notification struct {
	Feed         rss.UserFeed
	Item         rss.Item
	Channel      string
	User         string
	DeliveryHook func() error
}

func (n Notification) Address() string {
	return n.Channel
}

func (n Notification) Deliver(deliver func(address string, content string) error) {
	if deliver(n.Channel, n.String()) == nil {
		err := n.DeliveryHook()
		if err != nil {
			log.Error(err)
		}
	}
}

func (n Notification) String() string {
	i := n.Item
	t := "%s\n%s\n%s\n%s\n"
	t = fmt.Sprintf(t,
		i.Title,
		i.DescriptionTruncated(),
		i.Link,
		n.User,
	)

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
		i.Title, i.Link, n.User,
	)
}
