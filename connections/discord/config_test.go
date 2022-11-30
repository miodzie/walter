// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package discord

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseConfig_parses_into_Config_struct(t *testing.T) {
	val := make(map[string]interface{})
	val["type"] = "discord"
	val["token"] = "my_token"
	val["admins"] = []interface{}{"alice"}
	mods := make(map[string]interface{})
	mods["my_plugin"] = struct{}{}
	val["mods"] = mods

	cfg, err := ParseConfig(val)

	assert.Nil(t, err)
	assert.Equal(t, val["token"], cfg.Token)
	assert.Len(t, cfg.Admins, 1)
	assert.Len(t, cfg.Mods, 1)
}

func TestParseConfig_returns_error_if_type_is_not_discord(t *testing.T) {
	val := make(map[string]interface{})
	val["type"] = "irc"

	_, err := ParseConfig(val)

	assert.Equal(t, ErrIncorrectType, err)
}
