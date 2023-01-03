// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package domain

import (
	"time"
)

// Feed is the allowed and available web feeds that users can subscribe
// to.
type Feed struct {
	// TODO: Id is a database concept, should be removed.
	Id            uint64
	Name          string
	Url           string
	LastPublished time.Time
}
