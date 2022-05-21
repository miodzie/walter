package discord

import (
	"os"

	"github.com/miodzie/seras"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Admins map[string]bool
	Mods   []string
	Token  string
}

type ConfigParser struct {
}

func (c *ConfigParser) Parse(name, file string) (seras.Connection, error) {
	var cfg Config
	f, err := os.ReadFile(file)
	if err != nil {
	  return nil, err
	}
	toml.Unmarshal(f, cfg)

	disc, err := New(cfg.Token)
	if err != nil {
		return disc, err
	}

	return disc, nil
}
