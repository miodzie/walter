// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package moderator

import (
	"regexp"

	"github.com/miodzie/walter"
)

// If the string has more than 3 all caps words.
const allCapsRegex = `(\b[A-Z]+\s?\b){3,}`

func IsSpam(msg walter.Message) bool {
	r, _ := regexp.Compile(allCapsRegex)

	return r.MatchString(msg.Content)
}
