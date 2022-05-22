package discord

import (
	"errors"

	"github.com/miodzie/seras"
)

var ErrIncorrectType = errors.New("config is not of type: 'discord'")

type Config struct {
	Admins map[string]bool
	Mods   []string
	Token  string
}

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
	cfg.Admins, ok = val["admins"].(map[string]bool)
	if !ok {
		return cfg, errors.New("unable to parse admin")
	}
	cfg.Mods, ok = val["mods"].([]string)
	if !ok {
		return cfg, errors.New("unable to parse mods")
	}

	return cfg, nil
}

type ConfigParser struct {
}

func (c *ConfigParser) Parse(val map[string]interface{}) (seras.Connection, error) {
	cfg, err := ParseConfig(val)
	if err != nil {
		return nil, err
	}

	disc, err := New(cfg.Token)
	if err != nil {
		return disc, err
	}

	return disc, nil
}
