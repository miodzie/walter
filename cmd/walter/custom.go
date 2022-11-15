// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

//go:build custom
// +build custom

package main

import (
	"github.com/miodzie/walter/mods"
	"github.com/miodzie/walter/mods/bestbot"
)

func init() {
	mods.Register("best_bot", &bestbot.ModFactory{})
}
