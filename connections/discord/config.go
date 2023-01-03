// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package discord

import (
	"errors"

	"walter"
)

var ErrIncorrectType = errors.New("config is not of type: 'discord'")

type Config struct {
	Token string
	walter.BaseConnection
}

func init() {
	if err := walter.AddBotParser("discord", &BotParser{}); err != nil {
		panic(err)
	}
}

// ParseConfig
// TODO: I seriously need like an internal lib to map this automatically.
func ParseConfig(val map[string]interface{}) (Config, error) {
	var cfg Config

	if val["type"] != "discord" {
		return cfg, ErrIncorrectType
	}

	var ok bool
	cfg.Token, ok = val["token"].(string)
	if !ok {
		return cfg, errors.New("unable to parse token")
	}
	admins, ok := val["admins"].([]interface{})
	if !ok {
		return cfg, errors.New("unable to parse admin")
	}
	for _, a := range admins {
		cfg.Admins = append(cfg.Admins, a.(string))
	}

	cfg.Mods, ok = val["mods"].(map[string]any)

	return cfg, nil
}

type BotParser struct {
}

func (c *BotParser) Parse(val map[string]interface{}) (walter.Bot, error) {
	cfg, err := ParseConfig(val)
	if err != nil {
		return nil, err
	}

	disc, err := New(cfg)
	if err != nil {
		return disc, err
	}

	return disc, nil
}
