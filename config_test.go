// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package seras

import (
	"testing"
)

func TestParseToml(t *testing.T) {
	cfg, err := ParseToml("config.toml")
	if err != nil {
		t.Error(err)
	}

	if len(cfg.Bots) != 2 {
		t.Fail()
	}
	d, ok := cfg.Bots["discord"]
	if !ok {
		t.Fail()
	}
	if d["type"] != "discord" {
		t.Fail()
	}
}
