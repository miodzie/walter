package seras

import (
	"errors"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

var connectors map[string]ConfigParser

type Config struct {
	Mods        []string
	Connections map[string]map[string]interface{}
}

type ConfigParser interface {
	Parse(map[string]interface{}) (Connection, error)
}

func AddConnector(name string, parser ConfigParser) error {
	if connectors == nil {
		connectors = make(map[string]ConfigParser)
	}
	if _, ok := connectors[name]; ok {
		return errors.New("connector already registered")
	}
	connectors[name] = parser
	return nil
}

func ParseToml(file string) (*Config, error) {
	var c Config
	f, err := os.ReadFile(file)
	if err != nil {
		return &c, err
	}

	return &c, toml.Unmarshal(f, &c)
}
