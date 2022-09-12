package irc

import (
	"errors"
	"github.com/miodzie/seras"
)

var ErrIncorrectType = errors.New("config is not of type: 'irc'")

type Config struct {
	Server       string
	Nick         string
	Username     string
	Channels     []string
	Admins       []string
	Mods         []string
	SASL         bool
	SASLUsername string
	SASLPassword string
}

func ParseConfig(val map[string]interface{}) (Config, error) {
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
	if !ok {
		return cfg, errors.New("unable to parse sasl")
	}
	cfg.SASLUsername, ok = val["sasl_username"].(string)
	if !ok {
		return cfg, errors.New("unable to parse sasl_username")
	}
	cfg.SASLPassword, ok = val["sasl_password"].(string)
	if !ok {
		return cfg, errors.New("unable to parse sasl_password")
	}

	cfg.Server, ok = val["server"].(string)
	if !ok {
		return cfg, errors.New("unable to parse server")
	}
	admins, ok := val["admins"].([]interface{})
	if !ok {
		return cfg, errors.New("unable to parse admin")
	}
	for _, a := range admins {
		cfg.Admins = append(cfg.Admins, a.(string))
	}

	mods, ok := val["mods"].([]interface{})
	if !ok {
		return cfg, errors.New("unable to parse mods")
	}
	for _, a := range mods {
		cfg.Mods = append(cfg.Mods, a.(string))
	}

	return cfg, nil
}

type BotParser struct {
}

func (c *BotParser) Parse(val map[string]interface{}) (seras.Bot, error) {
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
