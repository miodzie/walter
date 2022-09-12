package seras

import (
	"errors"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

var connectors map[string]BotParser

type Config struct {
	Mods []string
	Bots map[string]map[string]interface{}
}

// BotParser intakes a map of config settings for a particular bot type.
type BotParser interface {
	Parse(map[string]interface{}) (Bot, error)
}

func AddBotParser(name string, parser BotParser) error {
	if connectors == nil {
		connectors = make(map[string]BotParser)
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
