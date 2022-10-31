// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Package sed is a toy implementation of sed for making _real_ corrections
// only, no silly seds.
package sed

import (
	"strings"
)

func ParseSed(str string) *Sed {
	sed := Sed{}
	split := strings.Split(str, "/")
	sed.Command, split = split[0], split[1:]
	if len(split)%2 != 0 {
		sed.Options, split = split[len(split)-1], split[:len(split)-1]
	}
	sed.Replacements = split

	return &sed
}

type Sed struct {
	Command      string
	Replacements []string
	Options      string
}

func (s Sed) Replace(str string) string {
	r := s.Replacements
	for len(r) != 0 {
		subject := r[0]
		r = r[1:]
		replacement := r[0]
		r = r[1:]

		if s.HasOption("g") {
			str = strings.ReplaceAll(str, subject, replacement)
		} else {
			str = strings.Replace(str, subject, replacement, 1)
		}
	}

	return str
}

func (s Sed) HasMatch(str string) bool {
	// 😎
	return str != s.Replace(str)
}

func (s Sed) HasOption(opt string) bool {
	return strings.Contains(s.Options, opt)
}
