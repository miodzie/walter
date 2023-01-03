// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package domain

import (
	"fmt"
	"strings"
)

type Notification struct {
	// TODO: Feed doesn't need to be on here.
	Feed    Feed
	Item    ParsedItem
	Channel string
	Users   []string
}

func (n Notification) String() string {
	i := n.Item
	t := "%s\n%s\n%s\n%s\n"
	t = fmt.Sprintf(t,
		i.Title,
		i.DescriptionTruncated(),
		i.Link,
		strings.Join(n.Users, ", "))

	return t
}
