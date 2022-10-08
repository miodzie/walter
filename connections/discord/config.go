package discord

import (
	"errors"

	"github.com/miodzie/seras"
)

var ErrIncorrectType = errors.New("config is not of type: 'discord'")

type Config struct {
	Admins []string
	Mods   []string
	Token  string
	seras.ConnectionConfig
}

func init() {
	if err := seras.AddBotParser("discord", &BotParser{}); err != nil {
		panic(err)
	}
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

	disc, err := New(cfg)
	if err != nil {
		return disc, err
	}

	return disc, nil
}
