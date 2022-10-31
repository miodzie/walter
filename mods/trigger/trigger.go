// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package trigger

import (
	"regexp"
	"strings"
)

const WORD = `(?i)\b{WORD}\b`

type Trigger struct {
	Id    uint64
	Word  string
	Reply string
	Regex string
}

func (t *Trigger) Check(str string) bool {
	reg := strings.ReplaceAll(WORD, "{WORD}", t.Word)
	r, _ := regexp.Compile(reg)

	return r.Match([]byte(str))
}

type Repository interface {
	All() ([]Trigger, error)
	Add(*Trigger) error
	Remove(id uint64)
}
