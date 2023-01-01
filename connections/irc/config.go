// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package irc

import (
	"errors"
	"github.com/miodzie/walter"
)

var ErrIncorrectType = errors.New("config is not of type: 'irc'")

type Config struct {
	Server       string
	Nick         string
	Username     string
	Channels     []string
	SASL         bool
	SASLUsername string
	SASLPassword string
	walter.BaseConnection
}

func init() {
	if err := walter.AddBotParser("irc", &BotParser{}); err != nil {
		panic(err)
	}
}

func ParseConfig(val map[string]any) (Config, error) {
	var cfg Config

	if val["type"] != "irc" {
		return cfg, ErrIncorrectType
	}

	var ok bool
	cfg.Nick, ok = val["nick"].(string)
	if !ok {
		return cfg, errors.New("unable to parse nick")
	}

	cfg.Username, ok = val["username"].(string)
	if !ok {
		return cfg, errors.New("unable to parse username")
	}

	cfg.SASL, ok = val["sasl"].(bool)
	cfg.SASLUsername, ok = val["sasl_username"].(string)
	cfg.SASLPassword, ok = val["sasl_password"].(string)

	cfg.Channels, _ = val["channels"].([]string)

	cfg.Server, ok = val["server"].(string)
	if !ok {
		return cfg, errors.New("unable to parse server")
	}
	admins, ok := val["admins"].([]any)
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

func (c *BotParser) Parse(val map[string]any) (walter.Bot, error) {
	cfg, err := ParseConfig(val)
	if err != nil {
		return nil, err
	}

	con, err := New(cfg)
	if err != nil {
		return con, err
	}

	return con, nil
}
