package seras

import (
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	Mods        []string
	Connections map[string]Con
}

type Con struct {
	Type   string
	Admins []string
	Mods   []string
}

type ConfigParser interface {
	Parse(file string) (Connection, error)
}

func ParseToml(file string) (*Config, error) {
	var c Config
	f, err := os.ReadFile(file)
	if err != nil {
		return &c, err
	}

	return &c, toml.Unmarshal(f, &c)
}
