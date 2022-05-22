package seras

import (
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	Mods        []string
	Connections map[string]map[string]interface{}
}

type ConfigParser interface {
	Parse(map[string]interface{}) (Connection, error)
}

func ParseToml(file string) (*Config, error) {
	var c Config
	f, err := os.ReadFile(file)
	if err != nil {
		return &c, err
	}

	return &c, toml.Unmarshal(f, &c)
}
