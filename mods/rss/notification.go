// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"fmt"
)

type Notification struct {
	Feed    UserFeed
	Item    Item
	Channel string
	User    string
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
