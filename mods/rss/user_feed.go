// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"time"
)

// UserFeed is the allowed web feeds that users can subscribe
// to.
type UserFeed struct {
	Id            uint64
	Name          string
	Url           string
	LastPublished time.Time
}
