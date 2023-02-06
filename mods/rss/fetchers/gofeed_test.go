// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package fetchers

import (
	"fmt"
	"testing"
)

func SkipTestItCanParseARedditFeed(t *testing.T) {
	p := GoFeed()
	feed, err := p.Fetch("https://www.reddit.com/r/news/.rss")
	if err != nil {
		fmt.Println(err)
		t.Fail()
		t.Log(err)
	}

	fmt.Println(feed)
}
